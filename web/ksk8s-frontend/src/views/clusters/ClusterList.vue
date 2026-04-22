<template>
  <div>
    <div class="page-header">
      <h2>{{ $t('cluster.clusters') }}</h2>
      <el-button v-if="!auth.isReadOnly" type="primary" @click="$router.push('/clusters/create')">
        <el-icon><Plus /></el-icon> {{ $t('cluster.createCluster') }}
      </el-button>
    </div>
    <el-table :data="clusters" v-loading="loading" stripe>
      <el-table-column prop="id" :label="$t('common.id')" width="60" />
      <el-table-column prop="name" :label="$t('common.name')" />
      <el-table-column prop="description" :label="$t('common.description')" />
      <el-table-column prop="status" :label="$t('common.status')" width="120">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push(`/clusters/${row.id}`)">{{ $t('common.manage') }}</el-button>
          <el-button v-if="!auth.isReadOnly" size="small" type="danger" @click="deleteCluster(row.id)">{{ $t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { clusterAPI } from '../../api/client'
import { requestWithOTP } from '../../composables/useOTPDialog'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const auth = useAuthStore()
const { t } = useI18n()
const clusters = ref([])
const loading = ref(false)

function statusType(status) {
  switch (status) {
    case 'active': return 'success'
    case 'draft': return 'info'
    case 'error': return 'danger'
    default: return undefined
  }
}

async function fetchClusters() {
  loading.value = true
  try {
    const res = await clusterAPI.list()
    clusters.value = res.data
  } catch (err) {
    ElMessage.error(t('cluster.failedToLoadList'))
  } finally {
    loading.value = false
  }
}

async function deleteCluster(id) {
  try {
    await ElMessageBox.confirm(t('cluster.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await requestWithOTP((code) => clusterAPI.delete(id, code))
    ElMessage.success(t('cluster.clusterDeleted'))
    fetchClusters()
  } catch (err) {
    if (err !== 'cancel' && err.response?.data?.error !== 'otp_required') {
      ElMessage.error(t(`error.${err.response?.data?.error}`) || t('cluster.failedToDelete'))
    }
  }
}

onMounted(fetchClusters)
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
