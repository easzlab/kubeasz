<template>
  <div>
    <h2>{{ $t('settings.title') }}</h2>

    <el-card shadow="hover" style="max-width: 600px; margin-top: 20px">
      <template #header>
        <span>{{ $t('settings.accountInfo') }}</span>
      </template>
      <el-form label-width="120px">
        <el-form-item :label="$t('auth.username')">
          <el-input :model-value="auth.user?.username" disabled />
        </el-form-item>
        <el-form-item :label="$t('user.role')">
          <el-tag :type="roleTagType">{{ roleLabel }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="hover" style="max-width: 600px; margin-top: 20px">
      <template #header>
        <span>{{ $t('settings.otp') }}</span>
      </template>
      <div class="op-desc">{{ $t('settings.otpDesc') }}</div>

      <div v-if="auth.otpEnabled" style="margin-top: 16px">
        <el-alert type="success" :closable="false">{{ $t('settings.otpEnabled') }}</el-alert>
        <el-form :model="disableForm" label-width="120px" style="margin-top: 12px">
          <el-form-item :label="$t('auth.password')">
            <el-input v-model="disableForm.password" type="password" :placeholder="$t('auth.password')" />
          </el-form-item>
          <el-form-item :label="$t('auth.otpCode')">
            <el-input v-model="disableForm.code" :placeholder="$t('auth.otpCode')" />
          </el-form-item>
          <el-form-item>
            <el-button type="danger" size="small" @click="disableOTP" :loading="otpLoading">{{ $t('settings.disableOTP') }}</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div v-else style="margin-top: 16px">
        <div v-if="setupData.secret">
          <p>{{ $t('settings.scanQR') }}</p>
          <div style="margin: 12px 0">
            <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="QR Code" style="width: 200px; height: 200px" />
            <el-skeleton v-else style="width: 200px; height: 200px" animated />
          </div>
          <p style="font-size: 12px; color: #909399">{{ $t('settings.secret') }}: {{ setupData.secret }}</p>
          <el-form :model="verifyForm" label-width="120px" style="margin-top: 12px">
            <el-form-item :label="$t('auth.otpCode')">
              <el-input v-model="verifyForm.code" :placeholder="$t('settings.enterOTP')" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="small" @click="verifyOTP" :loading="otpLoading">{{ $t('settings.verifyEnable') }}</el-button>
              <el-button size="small" @click="cancelSetup">{{ $t('common.cancel') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div v-else>
          <el-button type="primary" size="small" @click="setupOTP" :loading="otpLoading">{{ $t('settings.setupOTP') }}</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const auth = useAuthStore()
const otpLoading = ref(false)
const setupData = ref({ secret: '', url: '' })
const verifyForm = ref({ code: '' })
const disableForm = ref({ password: '', code: '' })
const qrCodeDataUrl = ref('')

function normalizeRole(role) {
  if (role === 'admin') return 'platform_admin'
  if (role === 'viewer') return 'security_auditor'
  return role
}

const roleLabel = computed(() => {
  const role = normalizeRole(auth.user?.role)
  return t(`role.${role}`) || auth.user?.role
})

const roleTagType = computed(() => {
  const map = {
    platform_admin: 'danger',
    cluster_admin: 'warning',
    security_auditor: 'info'
  }
  return map[normalizeRole(auth.user?.role)] || ''
})

watch(() => setupData.value.url, async (url) => {
  if (!url) {
    qrCodeDataUrl.value = ''
    return
  }
  try {
    qrCodeDataUrl.value = await QRCode.toDataURL(url, { width: 200, margin: 2 })
  } catch {
    qrCodeDataUrl.value = ''
  }
})

async function setupOTP() {
  otpLoading.value = true
  try {
    const data = await auth.otpSetup()
    setupData.value = { secret: data.secret, url: data.url }
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('settings.failedToSetup'))
  } finally {
    otpLoading.value = false
  }
}

async function verifyOTP() {
  if (!verifyForm.value.code) {
    ElMessage.warning(t('settings.enterOTPCode'))
    return
  }
  otpLoading.value = true
  try {
    await auth.otpVerify(verifyForm.value.code)
    ElMessage.success(t('settings.otpEnabledSuccess'))
    setupData.value = { secret: '', url: '' }
    verifyForm.value.code = ''
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('error.invalid_otp'))
  } finally {
    otpLoading.value = false
  }
}

async function disableOTP() {
  if (!disableForm.value.password || !disableForm.value.code) {
    ElMessage.warning(t('settings.enterPasswordAndOTP'))
    return
  }
  otpLoading.value = true
  try {
    await auth.otpDisable(disableForm.value.password, disableForm.value.code)
    ElMessage.success(t('settings.otpDisabled'))
    disableForm.value = { password: '', code: '' }
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('settings.failedToDisable'))
  } finally {
    otpLoading.value = false
  }
}

function cancelSetup() {
  setupData.value = { secret: '', url: '' }
  verifyForm.value.code = ''
}
</script>

<style scoped>
.op-desc {
  color: #909399;
  font-size: 13px;
  margin-top: 4px;
}
</style>
