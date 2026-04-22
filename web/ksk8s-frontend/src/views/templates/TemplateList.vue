<template>
  <div>
    <div class="page-header">
      <h2>{{ $t('template.templates') }}</h2>
    </div>
    <el-table :data="templates" v-loading="loading" stripe
    >
      <el-table-column prop="id" :label="$t('common.id')" width="60" />
      <el-table-column prop="name" :label="$t('common.name')" />
      <el-table-column prop="description" :label="$t('common.description')" />
      <el-table-column prop="is_default" :label="$t('cluster.default')" width="100"
      >
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success">{{ $t('cluster.default') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="280"
      >
        <template #default="{ row }">
          <el-button size="small" @click="$router.push(`/templates/${row.id}`)">{{ $t('common.view') }}</el-button>
          <el-button size="small" @click="setDefault(row.id)" :disabled="row.is_default"
          >{{ $t('template.setDefault') }}</el-button>
          <el-button size="small" type="danger" @click="deleteTemplate(row.id)"
          >{{ $t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { templateAPI } from '../../api/client'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const templates = ref([])
const loading = ref(false)

async function fetchTemplates() {
  loading.value = true
  try {
    const res = await templateAPI.list()
    templates.value = res.data
  } catch (err) {
    ElMessage.error(t('template.failedToLoad'))
  } finally {
    loading.value = false
  }
}

async function setDefault(id) {
  try {
    await templateAPI.setDefault(id)
    ElMessage.success(t('template.defaultTemplateUpdated'))
    fetchTemplates()
  } catch (err) {
    ElMessage.error(t('template.failedToSetDefault'))
  }
}

async function deleteTemplate(id) {
  try {
    await ElMessageBox.confirm(t('template.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await templateAPI.delete(id)
    ElMessage.success(t('template.templateDeleted'))
    fetchTemplates()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(t('template.failedToDelete'))
    }
  }
}

onMounted(fetchTemplates)
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
