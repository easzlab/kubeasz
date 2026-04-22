<template>
  <div>
    <el-steps :active="currentStepIndex" finish-status="success" align-center>
      <el-step v-for="step in steps" :key="step.num" :title="step.label" :description="step.desc" />
    </el-steps>

    <div class="current-step-hint" v-if="currentStep">
      <el-alert
        :title="t('install.currentStep', { current: currentStepIndex + 1, total: steps.length, label: currentStep.label, desc: currentStep.desc })"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="step-actions" v-if="currentStep">
      <el-button
        v-if="currentStepIndex > 0"
        @click="goBackStep"
        :disabled="isInstalling"
      >
        {{ $t('install.goBack') }}
      </el-button>
      <el-button
        v-if="!activeTask"
        type="primary"
        @click="runStep(currentStep.num)"
        :loading="starting"
      >
        {{ $t('install.runStep', { num: currentStep.num }) }}
      </el-button>
      <el-button
        v-if="activeTask && activeTask.status === 'running'"
        type="danger"
        @click="abortTask"
        :loading="aborting"
      >
        {{ $t('install.abort') }}
      </el-button>
      <el-button
        v-if="activeTask && activeTask.status === 'awaiting_approval'"
        type="success"
        @click="approveTask"
        :loading="approving"
      >
        {{ $t('install.approveContinue') }}
      </el-button>
      <el-button
        v-if="activeTask && activeTask.status !== 'running' && activeTask.status !== 'awaiting_approval'"
        type="warning"
        @click="retryStep"
        :loading="starting"
      >
        {{ $t('install.retry') }}
      </el-button>
      <el-button
        v-if="activeTask && activeTask.status !== 'running' && activeTask.status !== 'awaiting_approval' && currentStepIndex < steps.length - 1"
        type="success"
        @click="continueStep"
        :loading="starting"
      >
        {{ $t('install.continue') }}
      </el-button>
      <el-button
        v-if="igniteIP"
        type="info"
        @click="openIgniteSSH"
      >
        {{ $t('install.websshMaster') }}
      </el-button>
      <el-button
        type="info"
        @click="openKubeaszSSH"
      >
        {{ $t('install.websshIgnite') }}
      </el-button>
    </div>

    <div v-if="activeTask" class="task-info">
      <el-tag :type="taskStatusType(activeTask.status)">{{ activeTask.status }}</el-tag>
      <span class="task-meta">Task #{{ activeTask.id }}</span>
    </div>

    <LogViewer v-if="activeTask" :task-id="activeTask.id" :cluster-id="clusterId" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { taskAPI, clusterAPI } from '../../api/client'
import { requestWithOTP } from '../../composables/useOTPDialog'
import { ElMessage } from 'element-plus'
import LogViewer from '../../components/LogViewer.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const clusterId = computed(() => route.params.id)

const steps = [
  { num: '01', label: t('install.prepare'), desc: t('install.prepareDesc') },
  { num: '02', label: t('install.etcd'), desc: t('install.etcdDesc') },
  { num: '03', label: t('install.runtime'), desc: t('install.runtimeDesc') },
  { num: '04', label: t('install.master'), desc: t('install.masterDesc') },
  { num: '05', label: t('install.worker'), desc: t('install.workerDesc') },
  { num: '06', label: t('install.network'), desc: t('install.networkDesc') },
  { num: '07', label: t('install.addons'), desc: t('install.addonsDesc') }
]

const currentStepIndex = ref(0)
const activeTask = ref(null)
const starting = ref(false)
const aborting = ref(false)
const approving = ref(false)
const igniteIP = ref('')
let pollInterval = null

const currentStep = computed(() => steps[currentStepIndex.value])

const isInstalling = computed(() => {
  return activeTask.value && activeTask.value.status === 'running'
})

async function loadProgress() {
  try {
    const res = await clusterAPI.get(clusterId.value)
    if (res.data.install_step_index != null && res.data.install_step_index >= 0) {
      currentStepIndex.value = res.data.install_step_index
    }
  } catch (err) {
    console.error('Failed to load install progress', err)
  }
}

async function saveProgress(index) {
  try {
    await clusterAPI.update(clusterId.value, { install_step_index: index })
  } catch (err) {
    console.error('Failed to save install progress', err)
  }
}

async function loadIgniteIP() {
  try {
    const res = await clusterAPI.getConfig(clusterId.value)
    const masters = res.data.nodes?.kube_master || []
    if (masters.length > 0) {
      igniteIP.value = masters[0].ip_address
    }
  } catch (err) {
    console.error('Failed to load ignite IP', err)
  }
}

async function openIgniteSSH() {
  if (!igniteIP.value) {
    ElMessage.warning(t('install.noMasterIP'))
    return
  }
  try {
    await requestWithOTP((code) => {
      if (code) return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: { ok: true } })
    })
  } catch {
    return
  }
  const url = router.resolve({
    path: '/webssh',
    query: { addr: `${igniteIP.value}:22`, user: 'root' }
  })
  window.open(url.href, '_blank')
}

function openKubeaszSSH() {
  const host = window.location.hostname
  const url = router.resolve({
    path: '/webssh',
    query: { addr: `${host}:22`, user: 'root' }
  })
  window.open(url.href, '_blank')
}

function taskStatusType(status) {
  switch (status) {
    case 'running': return 'primary'
    case 'success': return 'success'
    case 'failed': return 'danger'
    case 'aborted': return 'warning'
    case 'awaiting_approval': return 'info'
    default: return undefined
  }
}

async function runStep(stepNum) {
  starting.value = true
  try {
    const res = await requestWithOTP((code) => taskAPI.run(clusterId.value, stepNum, code))
    activeTask.value = res.data
    ElMessage.success(t('install.stepStarted', { num: stepNum }))
    startPolling()
  } catch (err) {
    if (err.response?.data?.error !== 'otp_required') {
      ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('install.failedToStart'))
    }
  } finally {
    starting.value = false
  }
}

async function abortTask() {
  if (!activeTask.value) return
  aborting.value = true
  try {
    await requestWithOTP((code) => taskAPI.abort(clusterId.value, activeTask.value.id, code))
    ElMessage.info(t('install.taskAborted'))
    await fetchTaskStatus()
  } catch (err) {
    if (err.response?.data?.error !== 'otp_required') {
      ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('install.failedToAbort'))
    }
  } finally {
    aborting.value = false
  }
}

async function approveTask() {
  if (!activeTask.value) return
  approving.value = true
  try {
    await requestWithOTP((code) => taskAPI.approve(clusterId.value, activeTask.value.id, code))
    ElMessage.success(t('install.taskApproved'))
    activeTask.value = null
  } catch (err) {
    if (err.response?.data?.error !== 'otp_required') {
      ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('install.failedToApprove'))
    }
  } finally {
    approving.value = false
  }
}

async function retryStep() {
  if (!currentStep.value) return
  activeTask.value = null
  await runStep(currentStep.value.num)
}

async function goBackStep() {
  if (currentStepIndex.value > 0) {
    activeTask.value = null
    currentStepIndex.value--
    await saveProgress(currentStepIndex.value)
    ElMessage.success(t('install.backToStep', { num: currentStepIndex.value + 1, label: steps[currentStepIndex.value].label }))
  }
}

async function continueStep() {
  activeTask.value = null
  if (currentStepIndex.value < steps.length - 1) {
    currentStepIndex.value++
    await saveProgress(currentStepIndex.value)
    const nextStep = steps[currentStepIndex.value].num
    await runStep(nextStep)
  }
}

async function fetchTaskStatus() {
  if (!activeTask.value) return
  try {
    const res = await taskAPI.status(clusterId.value, activeTask.value.id)
    activeTask.value = res.data.task
  } catch (err) {
    console.error('Failed to fetch task status', err)
  }
}

function startPolling() {
  if (pollInterval) clearInterval(pollInterval)
  pollInterval = setInterval(fetchTaskStatus, 2000)
}

onMounted(() => {
  loadProgress()
  loadIgniteIP()
  startPolling()
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
.current-step-hint {
  margin: 20px 0;
}
.step-actions {
  margin: 30px 0;
  text-align: center;
}
.task-info {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}
.task-meta {
  color: #909399;
  font-size: 13px;
}
</style>
