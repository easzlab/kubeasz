<template>
  <div>
    <h2>{{ $t('cluster.createCluster') }}</h2>
    <el-form :model="form" label-width="120px" style="max-width: 600px; margin-top: 20px">
      <el-form-item :label="$t('cluster.clusterName')" required>
        <el-input v-model="form.name" :placeholder="$t('cluster.enterClusterName')" />
      </el-form-item>
      <el-form-item :label="$t('common.description')">
        <el-input v-model="form.description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item :label="$t('cluster.template')">
        <el-select v-model="form.template_id" :placeholder="$t('cluster.selectTemplate')" clearable style="width: 100%">
          <el-option
            v-for="t in templates"
            :key="t.id"
            :label="t.name + (t.is_default ? ` (${$t('cluster.default')})` : '')"
            :value="t.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submit" :loading="loading">{{ $t('common.create') }}</el-button>
        <el-button @click="$router.push('/clusters')">{{ $t('common.cancel') }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { clusterAPI, templateAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()
const loading = ref(false)
const templates = ref([])
const form = ref({ name: '', description: '', template_id: null })

async function fetchTemplates() {
  try {
    const res = await templateAPI.list()
    templates.value = res.data
  } catch (err) {
    console.error(t('template.failedToLoad'), err)
  }
}

async function submit() {
  if (!form.value.name) {
    ElMessage.warning(t('cluster.clusterName') + ' ' + t('common.required'))
    return
  }
  loading.value = true
  try {
    const payload = { ...form.value }
    if (payload.template_id) {
      payload.template_id = Number(payload.template_id)
    }
    const res = await clusterAPI.create(payload)
    ElMessage.success(t('cluster.clusterCreated'))
    router.push(`/clusters/${res.data.id}`)
  } catch (err) {
    ElMessage.error(t(`error.${err.response?.data?.error}`) || t('cluster.failedToCreate'))
  } finally {
    loading.value = false
  }
}

onMounted(fetchTemplates)
</script>
