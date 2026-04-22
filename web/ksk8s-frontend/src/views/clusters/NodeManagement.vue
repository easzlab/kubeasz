<template>
  <div>
    <h3>{{ $t('node.title') }}</h3>

    <el-alert
      v-if="nodes.length === 0"
      :title="$t('node.noNodesHint')"
      type="info"
      :closable="false"
      style="margin: 16px 0"
    />

    <el-table v-if="nodes.length > 0" :data="nodes" style="width: 100%; margin-top: 16px">
      <el-table-column prop="group_name" :label="$t('node.group')" width="120" />
      <el-table-column prop="ip_address" :label="$t('node.ipAddress')" width="160" />
      <el-table-column prop="k8s_nodename" :label="$t('node.nodeName')" />
      <el-table-column :label="$t('common.actions')" width="280">
        <template #default="scope">
          <el-button
            v-if="!auth.isReadOnly"
            size="small"
            type="primary"
            @click="openSSH(scope.row.ip_address)"
          >
            SSH
          </el-button>
          <el-button
            v-if="!auth.isReadOnly"
            size="small"
            type="danger"
            :loading="deleteLoading[scope.row.ip_address]"
            @click="confirmDelete(scope.row)"
          >
            {{ $t('node.deleteNode') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <template v-if="!auth.isReadOnly">
    <h4 style="margin-top: 24px">{{ $t('node.addNode') }}</h4>
    <el-form :model="newNode" label-width="120px" style="max-width: 500px">
      <el-form-item :label="$t('node.group')">
        <el-select v-model="newNode.group_name">
          <el-option label="etcd" value="etcd" />
          <el-option label="kube_master" value="kube_master" />
          <el-option label="kube_node" value="kube_node" />
          <el-option label="harbor" value="harbor" />
          <el-option label="ex_lb" value="ex_lb" />
          <el-option label="chrony" value="chrony" />
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('node.ipAddress')">
        <el-input v-model="newNode.ip_address" :placeholder="$t('node.placeholderIP')" />
      </el-form-item>
      <el-form-item :label="$t('node.nodeName')">
        <el-input v-model="newNode.k8s_nodename" :placeholder="$t('node.placeholderNodeName')" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="addNode" :loading="loading">{{ $t('node.addNode') }}</el-button>
      </el-form-item>
    </el-form>
    </template>

    <el-card v-if="activeTask" style="margin-top: 24px" shadow="hover">
      <template #header>
        <span>{{ $t('node.activeTask') }}: {{ activeTask.task_type }} — {{ activeTask.status }}</span>
      </template>
      <p>Task ID: {{ activeTask.id }}</p>
      <p>Status: <el-tag :type="taskStatusType(activeTask.status)">{{ activeTask.status }}</el-tag></p>
      <p v-if="activeTask.error_message" style="color: #f56c6c">{{ activeTask.error_message }}</p>
      <el-button size="small" @click="viewTaskLogs">{{ $t('node.viewLogs') }}</el-button>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { clusterAPI, nodeOpsAPI } from '../../api/client'
import { requestWithOTP } from '../../composables/useOTPDialog'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()
const nodes = ref([])
const loading = ref(false)
const deleteLoading = ref({})
const activeTask = ref(null)

const newNode = ref({
  group_name: 'kube_node',
  ip_address: '',
  k8s_nodename: ''
})

async function loadNodes() {
  try {
    const res = await clusterAPI.getConfig(route.params.id)
    const data = res.data
    const all = []
    for (const g of ['etcd', 'kube_master', 'kube_node', 'harbor', 'ex_lb', 'chrony']) {
      for (const n of (data.nodes[g] || [])) {
        all.push({ ...n, group_name: g })
      }
    }
    nodes.value = all
  } catch (err) {
    console.error(t('node.failedToLoad'), err)
  }
}

async function openSSH(ip) {
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
    query: { addr: `${ip}:22`, user: 'root' }
  })
  window.open(url.href, '_blank')
}

async function confirmDelete(node) {
  try {
    await ElMessageBox.confirm(
      t('node.deleteConfirm', { ip: node.ip_address, group: node.group_name }),
      t('common.confirm'),
      { confirmButtonText: t('node.deleteNode'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    deleteLoading.value[node.ip_address] = true
    try {
      const res = await requestWithOTP((code) =>
        nodeOpsAPI.remove(route.params.id, {
          group_name: node.group_name,
          ip_address: node.ip_address
        }, code)
      )
      ElMessage.success(t(`message.${res.data.message}`) || res.data.message)
      if (res.data.task) {
        activeTask.value = res.data.task
        pollTaskStatus(res.data.task.id)
      }
      await loadNodes()
    } catch (err) {
      if (err.response?.data?.error !== 'otp_required') {
        ElMessage.error(t(`error.${err.response?.data?.error}`) || t('node.failedToDelete'))
      }
    } finally {
      deleteLoading.value[node.ip_address] = false
    }
  } catch {
  }
}

async function addNode() {
  if (!newNode.value.ip_address) {
    ElMessage.error(t('node.ipRequired'))
    return
  }
  loading.value = true
  try {
    const res = await requestWithOTP((code) =>
      nodeOpsAPI.add(route.params.id, {
        group_name: newNode.value.group_name,
        ip_address: newNode.value.ip_address,
        k8s_nodename: newNode.value.k8s_nodename
      }, code)
    )
    ElMessage.success(t(`message.${res.data.message}`) || res.data.message)
    if (res.data.task) {
      activeTask.value = res.data.task
      pollTaskStatus(res.data.task.id)
    }
    newNode.value = { group_name: 'kube_node', ip_address: '', k8s_nodename: '' }
    await loadNodes()
  } catch (err) {
    if (err.response?.data?.error !== 'otp_required') {
      ElMessage.error(t(`error.${err.response?.data?.error}`) || t('node.failedToAdd'))
    }
  } finally {
    loading.value = false
  }
}

async function pollTaskStatus(taskId) {
  const interval = setInterval(async () => {
    try {
      const res = await clusterAPI.get(route.params.id)
      if (activeTask.value && activeTask.value.status !== 'running') {
        clearInterval(interval)
      }
    } catch {
      clearInterval(interval)
    }
  }, 3000)
}

function viewTaskLogs() {
  if (!activeTask.value) return
  router.push(`/clusters/${route.params.id}/logs`)
}

function taskStatusType(status) {
  switch (status) {
    case 'success': return 'success'
    case 'running': return 'primary'
    case 'failed': case 'aborted': return 'danger'
    case 'awaiting_approval': return 'warning'
    default: return 'info'
  }
}

onMounted(loadNodes)
</script>
