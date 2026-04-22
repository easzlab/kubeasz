import axios from 'axios'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('ksk8s_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('ksk8s_token')
      localStorage.removeItem('ksk8s_user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default client

function withOTPHeader(config, otpCode) {
  if (otpCode) {
    config = config || {}
    config.headers = config.headers || {}
    config.headers['X-OTP-Code'] = otpCode
  }
  return config
}

export const authAPI = {
  login: (data) => client.post('/auth/login', data),
  register: (data) => client.post('/auth/register', data),
  settings: () => client.get('/auth/settings'),
  otpSetup: () => client.post('/auth/otp/setup'),
  otpVerify: (data) => client.post('/auth/otp/verify', data),
  otpDisable: (data) => client.post('/auth/otp/disable', data)
}

export const clusterAPI = {
  list: () => client.get('/clusters'),
  get: (id) => client.get(`/clusters/${id}`),
  create: (data) => client.post('/clusters', data),
  update: (id, data) => client.put(`/clusters/${id}`, data),
  delete: (id, otpCode) => client.delete(`/clusters/${id}`, withOTPHeader(null, otpCode)),
  getConfig: (id) => client.get(`/clusters/${id}/config`),
  saveConfig: (id, data) => client.put(`/clusters/${id}/config`, data),
  generateConfig: (id) => client.post(`/clusters/${id}/generate-config`)
}

export const nodeOpsAPI = {
  add: (clusterId, data, otpCode) => client.post(`/clusters/${clusterId}/nodes`, data, withOTPHeader(null, otpCode)),
  remove: (clusterId, data, otpCode) => client.delete(`/clusters/${clusterId}/nodes`, { ...withOTPHeader(null, otpCode), data })
}

export const taskAPI = {
  list: (clusterId) => client.get(`/clusters/${clusterId}/tasks`),
  get: (clusterId, taskId) => client.get(`/clusters/${clusterId}/tasks/${taskId}`),
  run: (clusterId, step, otpCode) => client.post(`/clusters/${clusterId}/steps/${step}/run`, null, withOTPHeader(null, otpCode)),
  abort: (clusterId, taskId, otpCode) => client.post(`/clusters/${clusterId}/tasks/${taskId}/abort`, null, withOTPHeader(null, otpCode)),
  approve: (clusterId, taskId, otpCode) => client.post(`/clusters/${clusterId}/tasks/${taskId}/approve`, null, withOTPHeader(null, otpCode)),
  retry: (clusterId, step, otpCode) => client.post(`/clusters/${clusterId}/steps/${step}/retry`, null, withOTPHeader(null, otpCode)),
  logs: (clusterId, taskId, offset = 0, limit = 1000) =>
    client.get(`/clusters/${clusterId}/tasks/${taskId}/logs`, { params: { offset, limit } }),
  status: (clusterId, taskId) => client.get(`/clusters/${clusterId}/tasks/${taskId}/status`)
}

export const auditAPI = {
  list: (params) => client.get('/audit-logs', { params })
}

export const userAPI = {
  list: () => client.get('/users'),
  get: (id) => client.get(`/users/${id}`),
  create: (data) => client.post('/users', data),
  updateRole: (id, data) => client.put(`/users/${id}/role`, data),
  resetPassword: (id, data) => client.put(`/users/${id}/password`, data),
  toggleOTP: (id, data) => client.put(`/users/${id}/otp`, data),
  updateLanguage: (id, data) => client.put(`/users/${id}/language`, data),
  getRegistrationSetting: () => client.get('/settings/registration'),
  setRegistrationSetting: (data) => client.put('/settings/registration', data),
  bindCluster: (data) => client.post('/users/bind-cluster', data),
  unbindCluster: (data) => client.post('/users/unbind-cluster', data),
  listBindings: (id) => client.get(`/users/${id}/bindings`)
}

export const templateAPI = {
  list: () => client.get('/templates'),
  get: (id) => client.get(`/templates/${id}`),
  getParsed: (id) => client.get(`/templates/${id}/parsed`),
  create: (data) => client.post('/templates', data),
  update: (id, data) => client.put(`/templates/${id}`, data),
  delete: (id) => client.delete(`/templates/${id}`),
  setDefault: (id) => client.post(`/templates/${id}/set-default`)
}
