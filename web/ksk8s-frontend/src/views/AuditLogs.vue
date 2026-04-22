<template>
  <div>
    <div class="page-header">
      <h2>{{ $t('audit.title') }}</h2>
    </div>

    <el-card shadow="never" style="margin-bottom: 16px">
      <el-form :model="filters" inline>
        <el-form-item :label="$t('audit.timeRange')">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            :range-separator="$t('audit.to')"
            :start-placeholder="$t('audit.start')"
            :end-placeholder="$t('audit.end')"
            value-format="YYYY-MM-DDTHH:mm:ss"
            @change="onDateChange"
          />
        </el-form-item>
        <el-form-item :label="$t('audit.action')">
          <el-input v-model="filters.action" placeholder="e.g. run_task" clearable />
        </el-form-item>
        <el-form-item :label="$t('audit.user')">
          <el-input v-model="filters.username" :placeholder="$t('audit.user')" clearable />
        </el-form-item>
        <el-form-item :label="$t('audit.highRisk')">
          <el-select v-model="filters.is_high_risk" :placeholder="$t('common.pleaseSelect')" clearable style="width: 100px">
            <el-option :label="$t('common.yes')" :value="true" />
            <el-option :label="$t('common.no')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchLogs">{{ $t('common.search') }}</el-button>
          <el-button @click="resetFilters">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="logs" v-loading="loading" stripe>
      <el-table-column prop="id" :label="$t('common.id')" width="60" />
      <el-table-column prop="created_at" :label="$t('audit.time')" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="username" :label="$t('audit.user')" width="120" />
      <el-table-column prop="action" :label="$t('audit.action')" width="160" />
      <el-table-column prop="resource_type" :label="$t('audit.resource')" width="120" />
      <el-table-column prop="resource_id" :label="$t('audit.resourceId')" width="100" />
      <el-table-column prop="status_code" :label="$t('common.status')" width="80">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status_code)" size="small">{{ row.status_code }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_high_risk" :label="$t('audit.risk')" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.is_high_risk" type="danger" size="small">{{ $t('audit.high') }}</el-tag>
          <el-tag v-else type="info" size="small">{{ $t('audit.normal') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip_address" :label="$t('audit.ip')" width="140" />
      <el-table-column prop="details" :label="$t('audit.details')">
        <template #default="{ row }">
          <pre style="margin: 0; font-size: 12px; max-height: 80px; overflow: auto">{{ formatDetails(row.details) }}</pre>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[20, 50, 100]"
      layout="total, sizes, prev, pager, next"
      style="margin-top: 20px"
      @change="fetchLogs"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { auditAPI } from '../api/client'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const logs = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(50)
const total = ref(0)
const dateRange = ref(null)

const filters = ref({
  start_time: '',
  end_time: '',
  action: '',
  username: '',
  is_high_risk: null
})

function onDateChange(val) {
  if (val && val.length === 2) {
    filters.value.start_time = val[0] ? val[0] + 'Z' : ''
    filters.value.end_time = val[1] ? val[1] + 'Z' : ''
  } else {
    filters.value.start_time = ''
    filters.value.end_time = ''
  }
}

function resetFilters() {
  dateRange.value = null
  filters.value = { start_time: '', end_time: '', action: '', username: '', is_high_risk: null }
  page.value = 1
  fetchLogs()
}

function statusType(code) {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 400 && code < 500) return 'warning'
  if (code >= 500) return 'danger'
  return 'info'
}

function formatTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleString()
}

function formatDetails(details) {
  if (!details) return ''
  try {
    const obj = JSON.parse(details)
    return JSON.stringify(obj, null, 2)
  } catch {
    return details
  }
}

async function fetchLogs() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filters.value.start_time) params.start_time = filters.value.start_time
    if (filters.value.end_time) params.end_time = filters.value.end_time
    if (filters.value.action) params.action = filters.value.action
    if (filters.value.username) params.username = filters.value.username
    if (filters.value.is_high_risk !== null && filters.value.is_high_risk !== undefined) {
      params.is_high_risk = filters.value.is_high_risk
    }
    const res = await auditAPI.list(params)
    logs.value = res.data.logs
    total.value = res.data.total
  } catch (err) {
    ElMessage.error(t('audit.failedToLoad'))
  } finally {
    loading.value = false
  }
}

onMounted(fetchLogs)
</script>
