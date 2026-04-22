import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('ksk8s_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('ksk8s_user') || 'null'))
  const otpCode = ref('')

  function normalizeRole(role) {
    if (role === 'admin') return 'platform_admin'
    if (role === 'viewer') return 'security_auditor'
    return role
  }

  const effectiveRole = computed(() => normalizeRole(user.value?.role))
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => effectiveRole.value === 'platform_admin')
  const isClusterAdmin = computed(() => effectiveRole.value === 'cluster_admin' || effectiveRole.value === 'platform_admin')
  const isPlatformAdmin = computed(() => effectiveRole.value === 'platform_admin')
  const isAuditor = computed(() => effectiveRole.value === 'security_auditor' || effectiveRole.value === 'platform_admin')
  const isReadOnly = computed(() => effectiveRole.value === 'security_auditor')
  const otpEnabled = computed(() => user.value?.otp_enabled || false)

  async function login(credentials) {
    const res = await authAPI.login(credentials)
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('ksk8s_token', token.value)
    localStorage.setItem('ksk8s_user', JSON.stringify(user.value))
    otpCode.value = ''
    return res.data
  }

  async function register(data) {
    return authAPI.register(data)
  }

  async function otpSetup() {
    const res = await authAPI.otpSetup()
    return res.data
  }

  async function otpVerify(code) {
    const res = await authAPI.otpVerify({ code })
    if (user.value) {
      user.value.otp_enabled = true
      localStorage.setItem('ksk8s_user', JSON.stringify(user.value))
    }
    return res.data
  }

  async function otpDisable(password, code) {
    const res = await authAPI.otpDisable({ password, otp_code: code })
    if (user.value) {
      user.value.otp_enabled = false
      localStorage.setItem('ksk8s_user', JSON.stringify(user.value))
    }
    return res.data
  }

  function logout() {
    token.value = ''
    user.value = null
    otpCode.value = ''
    localStorage.removeItem('ksk8s_token')
    localStorage.removeItem('ksk8s_user')
  }

  return {
    token, user, otpCode,
    isLoggedIn, isAdmin, isClusterAdmin, isPlatformAdmin, isAuditor, isReadOnly, otpEnabled,
    login, register, logout,
    otpSetup, otpVerify, otpDisable
  }
})
