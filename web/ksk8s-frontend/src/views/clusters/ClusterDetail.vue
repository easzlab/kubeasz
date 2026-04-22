<template>
  <div>
    <div class="cluster-header">
      <h2>{{ cluster.name }}</h2>
      <el-tag :type="statusType(cluster.status)">{{ cluster.status }}</el-tag>
    </div>
    <p class="description">{{ cluster.description }}</p>
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane :label="$t('cluster.config')" name="config">
        <router-view />
      </el-tab-pane>
      <el-tab-pane :label="$t('cluster.install')" name="install">
        <router-view />
      </el-tab-pane>
      <el-tab-pane :label="$t('cluster.lifecycle')" name="lifecycle">
        <router-view />
      </el-tab-pane>
      <el-tab-pane :label="$t('cluster.nodes')" name="nodes">
        <router-view />
      </el-tab-pane>
      <el-tab-pane :label="$t('cluster.security')" name="security">
        <router-view />
      </el-tab-pane>
      <el-tab-pane :label="$t('cluster.logs')" name="logs">
        <router-view />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { clusterAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const cluster = ref({})
const activeTab = ref(route.path.split('/').pop() || 'config')

function statusType(status) {
  switch (status) {
    case 'active': return 'success'
    case 'draft': return 'info'
    case 'error': return 'danger'
    default: return undefined
  }
}

async function fetchCluster() {
  try {
    const res = await clusterAPI.get(route.params.id)
    cluster.value = res.data
  } catch (err) {
    ElMessage.error(t('cluster.failedToLoad'))
  }
}

watch(activeTab, (tab) => {
  router.push(`/clusters/${route.params.id}/${tab}`)
})

watch(() => route.path, (path) => {
  const tab = path.split('/').pop()
  if (tab && tab !== activeTab.value) {
    activeTab.value = tab
  }
})

onMounted(fetchCluster)
</script>

<style scoped>
.cluster-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.description {
  color: #606266;
  margin-bottom: 20px;
}
</style>
