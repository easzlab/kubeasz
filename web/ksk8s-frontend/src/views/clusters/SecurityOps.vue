<template>
  <div>
    <h3>{{ $t('security.title') }}</h3>

    <el-card shadow="hover" style="max-width: 600px; margin-top: 20px">
      <div class="op-title">{{ $t('security.caRenewal') }}</div>
      <div class="op-desc">{{ $t('security.caRenewalDesc') }}</div>
      <el-button type="primary" size="small" @click="renewCA" style="margin-top: 12px">{{ $t('security.renewCA') }}</el-button>
    </el-card>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { taskAPI } from '../../api/client'
import { requestWithOTP } from '../../composables/useOTPDialog'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()

async function renewCA() {
  try {
    await ElMessageBox.confirm(t('security.renewConfirm'), t('common.confirm'), { type: 'warning' })
    const res = await requestWithOTP((code) => taskAPI.run(route.params.id, 'kca-renew', code))
    ElMessage.success(t('security.renewStarted', { id: res.data.id }))
  } catch (err) {
    if (err !== 'cancel' && err.response?.data?.error !== 'otp_required') {
      ElMessage.error(t(`error.${err.response?.data?.error}`) || t('security.failedToRenew'))
    }
  }
}
</script>

<style scoped>
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
