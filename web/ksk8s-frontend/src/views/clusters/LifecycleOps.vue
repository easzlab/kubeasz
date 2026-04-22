<template>
  <div>
    <h3>{{ $t('lifecycle.title') }}</h3>
    <el-row :gutter="16" style="margin-top: 20px">
      <el-col :span="8" v-for="op in visibleOps" :key="op.key">
        <el-card class="op-card" shadow="hover">
          <div class="op-title">{{ op.label }}</div>
          <div class="op-desc">{{ op.description }}</div>
          <el-button
            :type="op.danger ? 'danger' : 'primary'"
            size="small"
            @click="runOp(op)"
            style="margin-top: 12px"
          >
            {{ op.buttonText }}
          </el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="logDialogVisible" :title="$t('logs.taskLogs')" width="900px">
      <LogViewer v-if="activeTask" :task-id="activeTask.id" :cluster-id="route.params.id" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { taskAPI } from '../../api/client'
import { requestWithOTP } from '../../composables/useOTPDialog'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import LogViewer from '../../components/LogViewer.vue'

const route = useRoute()
const auth = useAuthStore()
const { t } = useI18n()

const allOps = [
  { key: 'start', label: t('lifecycle.start'), description: t('lifecycle.startDesc'), buttonText: t('lifecycle.startBtn'), danger: false },
  { key: 'stop', label: t('lifecycle.stop'), description: t('lifecycle.stopDesc'), buttonText: t('lifecycle.stopBtn'), danger: true },
  { key: 'upgrade', label: t('lifecycle.upgrade'), description: t('lifecycle.upgradeDesc'), buttonText: t('lifecycle.upgradeBtn'), danger: false },
  { key: 'backup', label: t('lifecycle.backup'), description: t('lifecycle.backupDesc'), buttonText: t('lifecycle.backupBtn'), danger: false },
  { key: 'restore', label: t('lifecycle.restore'), description: t('lifecycle.restoreDesc'), buttonText: t('lifecycle.restoreBtn'), danger: false },
  { key: 'destroy', label: t('lifecycle.destroy'), description: t('lifecycle.destroyDesc'), buttonText: t('lifecycle.destroyBtn'), danger: true }
]

const visibleOps = computed(() => {
  if (auth.isReadOnly) return []
  if (auth.isPlatformAdmin) return allOps
  return allOps.filter(op => op.key !== 'destroy' && op.key !== 'restore')
})

const logDialogVisible = ref(false)
const activeTask = ref(null)

async function runOp(op) {
  if (op.danger) {
    try {
      await ElMessageBox.confirm(
        t('lifecycle.dangerConfirm', { op: op.key }),
        t('common.dangerousOp'),
        { type: 'warning', confirmButtonClass: 'el-button--danger' }
      )
    } catch {
      return
    }
  }
  try {
    const res = await requestWithOTP((code) => taskAPI.run(route.params.id, op.key, code))
    activeTask.value = res.data
    ElMessage.success(t('lifecycle.taskStarted', { op: op.label, id: res.data.id }))
    logDialogVisible.value = true
  } catch (err) {
    if (err.response?.data?.error !== 'otp_required') {
      ElMessage.error(t(`error.${err.response?.data?.error}`) || t('lifecycle.failedToStart', { op: op.label }))
    }
  }
}
</script>

<style scoped>
.op-card {
  margin-bottom: 16px;
}
.op-title {
  font-weight: bold;
  font-size: 16px;
}
.op-desc {
  color: #909399;
  font-size: 13px;
  margin-top: 4px;
}
</style>
