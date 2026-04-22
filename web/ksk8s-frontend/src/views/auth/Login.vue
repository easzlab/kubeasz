<template>
  <div class="login-container">
    <el-card class="login-card" shadow="always">
      <template #header>
        <div class="login-header">
          <h2>ksk8s</h2>
          <p>{{ $t('auth.platformTitle') }}</p>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('auth.login')" name="login">
          <el-form :model="loginForm" @submit.prevent="handleLogin">
            <el-form-item>
              <el-input v-model="loginForm.username" :placeholder="$t('auth.username')" prefix-icon="User" />
            </el-form-item>
            <el-form-item>
              <el-input v-model="loginForm.password" type="password" :placeholder="$t('auth.password')" prefix-icon="Lock" show-password />
            </el-form-item>
            <el-form-item v-if="otpRequired">
              <el-input v-model="loginForm.otp_code" :placeholder="$t('auth.otpCode')" prefix-icon="Key" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleLogin" :loading="loading" style="width: 100%">
                {{ otpRequired ? $t('auth.verifyAndLogin') : $t('auth.login') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane v-if="regEnabled" :label="$t('auth.register')" name="register">
          <el-form :model="registerForm" @submit.prevent="handleRegister">
            <el-form-item>
              <el-input v-model="registerForm.username" :placeholder="$t('auth.username')" prefix-icon="User" />
            </el-form-item>
            <el-form-item>
              <el-input v-model="registerForm.password" type="password" :placeholder="$t('auth.password')" prefix-icon="Lock" show-password />
            </el-form-item>
            <el-form-item>
              <el-input v-model="registerForm.email" :placeholder="$t('auth.email')" prefix-icon="Message" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleRegister" :loading="loading" style="width: 100%">
                {{ $t('auth.register') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { authAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()
const loading = ref(false)
const activeTab = ref('login')
const otpRequired = ref(false)
const regEnabled = ref(false)

const loginForm = ref({ username: '', password: '', otp_code: '' })
const registerForm = ref({ username: '', password: '', email: '' })

onMounted(async () => {
  try {
    const res = await authAPI.settings()
    regEnabled.value = res.data.registration_enabled
  } catch {
    regEnabled.value = false
  }
})

async function handleLogin() {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning(t('auth.pleaseEnterUsernamePassword'))
    return
  }
  if (otpRequired.value && !loginForm.value.otp_code) {
    ElMessage.warning(t('auth.pleaseEnterOTP'))
    return
  }
  loading.value = true
  try {
    await auth.login(loginForm.value)
    ElMessage.success(t('auth.loginSuccess'))
    otpRequired.value = false
    router.push('/')
  } catch (err) {
    const error = err.response?.data?.error
    if (error === 'otp_required') {
      otpRequired.value = true
      ElMessage.warning(t('auth.otpRequired'))
    } else {
      ElMessage.error(t(`error.${err.response?.data?.message}`) || t(`error.${err.response?.data?.error}`) || t('auth.loginFailed'))
    }
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!registerForm.value.username || !registerForm.value.password) {
    ElMessage.warning(t('auth.pleaseEnterUsernamePassword'))
    return
  }
  loading.value = true
  try {
    await auth.register(registerForm.value)
    ElMessage.success(t('auth.registerSuccess'))
    activeTab.value = 'login'
  } catch (err) {
    ElMessage.error(t(`error.${err.response?.data?.error}`) || t('auth.registerFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1a1a2e;
}
.login-card {
  width: 400px;
}
.login-header {
  text-align: center;
}
.login-header h2 {
  margin: 0;
  color: #409eff;
}
.login-header p {
  margin: 8px 0 0;
  color: #909399;
  font-size: 14px;
}
</style>
