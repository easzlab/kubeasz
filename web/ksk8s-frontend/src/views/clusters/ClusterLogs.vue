<template>
  <div>
    <h3>{{ $t('logs.title') }}</h3>
    <el-table :data="tasks" v-loading="loading" stripe style="margin-top: 16px"
    >
      <el-table-column prop="id" :label="$t('common.id')" width="60" />
      <el-table-column prop="task_type" :label="$t('common.type')" />
      <el-table-column prop="step_number" :label="$t('common.step')" />
      <el-table-column prop="status" :label="$t('common.status')"
      >
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('common.created')" />
      <el-table-column :label="$t('common.actions')" width="120"
      >
        <template #default="{ row }">
          <el-button size="small" @click="viewLogs(row)">{{ $t('logs.taskLogs') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="logDialogVisible" :title="$t('logs.taskLogs')" width="900px"
    >
      <LogViewer v-if="selectedTask" :task-id="selectedTask.id" :cluster-id="clusterId" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { taskAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import LogViewer from '../../components/LogViewer.vue'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()
const clusterId = route.params.id
const tasks = ref([])
const loading = ref(false)
const logDialogVisible = ref(false)
const selectedTask = ref(null)

function statusType(status) {
  switch (status) {
    case 'running': return 'primary'
    case 'success': return 'success'
    case 'failed': return 'danger'
    case 'aborted': return 'warning'
    case 'awaiting_approval': return 'info'
    default: return undefined
  }
}

async function fetchTasks() {
  loading.value = true
  try {
    const res = await taskAPI.list(clusterId)
    tasks.value = res.data
  } catch (err) {
    ElMessage.error(t('logs.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function viewLogs(task) {
  selectedTask.value = task
  logDialogVisible.value = true
}

onMounted(fetchTasks)
</script>
