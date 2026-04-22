<template>
  <div>
    <div class="page-header">
      <h2>{{ $t('user.title') }}</h2>
      <el-button type="primary" @click="showCreate = true">
        <el-icon><Plus /></el-icon> {{ $t('user.createUser') }}
      </el-button>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="id" :label="$t('common.id')" width="60" />
      <el-table-column prop="username" :label="$t('auth.username')" />
      <el-table-column prop="email" :label="$t('auth.email')" />
      <el-table-column prop="role" :label="$t('user.role')" width="160">
        <template #default="{ row }">
          <el-tag :type="roleTagType(row.role)">{{ $t('role.' + row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="380">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">{{ $t('common.view') }}</el-button>
          <el-button size="small" @click="editRole(row)">{{ $t('user.changeRole') }}</el-button>
          <el-button size="small" @click="openResetPassword(row)">{{ $t('user.resetPassword') }}</el-button>
          <el-button
            v-if="row.role === 'cluster_admin'"
            size="small"
            type="primary"
            @click="openBind(row)"
          >
            {{ $t('user.bindCluster') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Create User Dialog -->
    <el-dialog v-model="showCreate" :title="$t('user.createUser')" width="400px">
      <el-form :model="newUser" label-width="80px">
        <el-form-item :label="$t('auth.username')">
          <el-input v-model="newUser.username" />
        </el-form-item>
        <el-form-item :label="$t('auth.password')">
          <el-input v-model="newUser.password" type="password" show-password />
        </el-form-item>
        <el-form-item :label="$t('auth.email')">
          <el-input v-model="newUser.email" />
        </el-form-item>
        <el-form-item :label="$t('user.role')">
          <el-select v-model="newUser.role">
            <el-option :label="$t('role.platform_admin')" value="platform_admin" />
            <el-option :label="$t('role.cluster_admin')" value="cluster_admin" />
            <el-option :label="$t('role.security_auditor')" value="security_auditor" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="createUser" :loading="saving">{{ $t('common.create') }}</el-button>
      </template>
    </el-dialog>

    <!-- User Detail Dialog -->
    <el-dialog v-model="showDetail" :title="$t('user.userDetail')" width="400px">
      <el-form v-if="detailUser" label-width="120px">
        <el-form-item :label="$t('common.id')">
          <span>{{ detailUser.id }}</span>
        </el-form-item>
        <el-form-item :label="$t('auth.username')">
          <span>{{ detailUser.username }}</span>
        </el-form-item>
        <el-form-item :label="$t('auth.email')">
          <span>{{ detailUser.email || '-' }}</span>
        </el-form-item>
        <el-form-item :label="$t('user.role')">
          <el-tag :type="roleTagType(detailUser.role)">{{ $t('role.' + detailUser.role) }}</el-tag>
        </el-form-item>
        <el-form-item :label="$t('user.otp')">
          <el-switch
            v-model="detailUser.otp_enabled"
            :active-text="$t('common.enabled')"
            :inactive-text="$t('common.disabled')"
            @change="toggleUserOTP"
          />
        </el-form-item>
        <el-form-item :label="$t('user.language')">
          <el-select v-model="detailUser.language" style="width: 150px" @change="updateUserLanguage">
            <el-option label="English" value="en" />
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="繁體中文" value="zh-TW" />
            <el-option label="Français" value="fr" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.created')">
          <span>{{ formatDate(detailUser.created_at) }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDetail = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- Edit Role Dialog -->
    <el-dialog v-model="showEdit" :title="$t('user.changeRole')" width="300px">
      <el-select v-model="editRoleValue" style="width: 100%">
        <el-option :label="$t('role.platform_admin')" value="platform_admin" />
        <el-option :label="$t('role.cluster_admin')" value="cluster_admin" />
        <el-option :label="$t('role.security_auditor')" value="security_auditor" />
      </el-select>
      <template #footer>
        <el-button @click="showEdit = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveRole" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Reset Password Dialog -->
    <el-dialog v-model="showResetPassword" :title="$t('user.resetPassword')" width="400px">
      <el-form :model="resetPasswordForm" label-width="100px">
        <el-form-item :label="$t('auth.username')">
          <span>{{ resetPasswordForm.username }}</span>
        </el-form-item>
        <el-form-item :label="$t('user.newPassword')">
          <el-input v-model="resetPasswordForm.password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResetPassword = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="doResetPassword" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Bind Cluster Dialog -->
    <el-dialog v-model="showBind" :title="$t('user.bindCluster')" width="500px">
      <p style="margin-bottom: 12px; color: #606266;">
        {{ $t('auth.username') }}: <strong>{{ bindUser?.username }}</strong>
      </p>
      <el-form label-width="120px">
        <el-form-item :label="$t('user.selectCluster')">
          <el-select v-model="selectedClusterId" style="width: 100%">
            <el-option
              v-for="c in availableClusters"
              :key="c.id"
              :label="c.name"
              :value="c.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBind = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="doBind" :loading="saving">{{ $t('user.bindCluster') }}</el-button>
      </template>

      <div v-if="userBindings.length > 0" style="margin-top: 16px">
        <el-divider />
        <p style="font-weight: bold; margin-bottom: 8px">{{ $t('user.alreadyBound') }}</p>
        <el-table :data="userBindings" size="small">
          <el-table-column prop="cluster_id" :label="$t('audit.resourceId')" width="100" />
          <el-table-column :label="$t('common.actions')" width="120">
            <template #default="{ row }">
              <el-button size="small" type="danger" @click="doUnbind(row.cluster_id)">{{ $t('user.unbind') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <!-- Registration Setting -->
    <el-card style="margin-top: 24px">
      <template #header>
        <span>{{ $t('user.systemSettings') }}</span>
      </template>
      <el-form label-width="180px">
        <el-form-item :label="$t('user.selfRegistration')">
          <el-switch
            v-model="regEnabled"
            :active-text="$t('common.enabled')"
            :inactive-text="$t('common.disabled')"
            @change="toggleRegistration"
          />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { userAPI, clusterAPI } from '../../api/client'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const users = ref([])
const clusters = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreate = ref(false)
const showEdit = ref(false)
const showBind = ref(false)
const showDetail = ref(false)
const showResetPassword = ref(false)
const editUser = ref(null)
const editRoleValue = ref('')
const regEnabled = ref(false)
const bindUser = ref(null)
const selectedClusterId = ref(null)
const userBindings = ref([])
const detailUser = ref(null)
const resetPasswordForm = ref({ userId: null, username: '', password: '' })

const newUser = ref({
  username: '',
  password: '',
  email: '',
  role: 'cluster_admin'
})

const availableClusters = computed(() => {
  const boundIds = new Set(userBindings.value.map(b => b.cluster_id))
  return clusters.value.filter(c => !boundIds.has(c.id))
})

const roleTagType = (role) => ({
  platform_admin: 'danger',
  cluster_admin: 'warning',
  security_auditor: 'info'
}[role] || '')

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString()
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await userAPI.list()
    users.value = res.data
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToLoad'))
  } finally {
    loading.value = false
  }
}

async function fetchClusters() {
  try {
    const res = await clusterAPI.list()
    clusters.value = res.data
  } catch {
    clusters.value = []
  }
}

async function fetchSettings() {
  try {
    const res = await userAPI.getRegistrationSetting()
    regEnabled.value = res.data.registration_enabled
  } catch {
    regEnabled.value = false
  }
}

async function createUser() {
  saving.value = true
  try {
    await userAPI.create(newUser.value)
    ElMessage.success(t('user.userCreated'))
    showCreate.value = false
    newUser.value = { username: '', password: '', email: '', role: 'cluster_admin' }
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToCreate'))
  } finally {
    saving.value = false
  }
}

function openDetail(row) {
  detailUser.value = { ...row }
  showDetail.value = true
}

async function toggleUserOTP(val) {
  saving.value = true
  try {
    const res = await userAPI.toggleOTP(detailUser.value.id, { enabled: val })
    ElMessage.success(t('user.otpUpdated'))
    detailUser.value.otp_enabled = res.data.otp_enabled
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToUpdateOTP'))
    detailUser.value.otp_enabled = !val
  } finally {
    saving.value = false
  }
}

async function updateUserLanguage(val) {
  saving.value = true
  try {
    const res = await userAPI.updateLanguage(detailUser.value.id, { language: val })
    ElMessage.success(t('user.languageUpdated'))
    detailUser.value.language = res.data.language
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToUpdateLanguage'))
  } finally {
    saving.value = false
  }
}

function editRole(row) {
  editUser.value = row
  editRoleValue.value = row.role
  showEdit.value = true
}

async function saveRole() {
  saving.value = true
  try {
    await userAPI.updateRole(editUser.value.id, { role: editRoleValue.value })
    ElMessage.success(t('user.roleUpdated'))
    showEdit.value = false
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToUpdateRole'))
  } finally {
    saving.value = false
  }
}

function openResetPassword(row) {
  resetPasswordForm.value = { userId: row.id, username: row.username, password: '' }
  showResetPassword.value = true
}

async function doResetPassword() {
  if (!resetPasswordForm.value.password || resetPasswordForm.value.password.length < 6) {
    ElMessage.warning(t('auth.password') + ' ' + t('common.required') + ' (min 6)')
    return
  }
  saving.value = true
  try {
    await userAPI.resetPassword(resetPasswordForm.value.userId, { password: resetPasswordForm.value.password })
    ElMessage.success(t('user.passwordUpdated'))
    showResetPassword.value = false
    resetPasswordForm.value = { userId: null, username: '', password: '' }
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToResetPassword'))
  } finally {
    saving.value = false
  }
}

async function openBind(row) {
  bindUser.value = row
  selectedClusterId.value = null
  showBind.value = true
  await loadBindings(row.id)
}

async function loadBindings(userId) {
  try {
    const res = await userAPI.listBindings(userId)
    userBindings.value = res.data
  } catch {
    userBindings.value = []
  }
}

async function doBind() {
  if (!selectedClusterId.value) {
    ElMessage.warning(t('user.pleaseSelectCluster'))
    return
  }
  saving.value = true
  try {
    await userAPI.bindCluster({
      user_id: bindUser.value.id,
      cluster_id: selectedClusterId.value
    })
    ElMessage.success(t('user.clusterBound'))
    selectedClusterId.value = null
    await loadBindings(bindUser.value.id)
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToBind'))
  } finally {
    saving.value = false
  }
}

async function doUnbind(clusterId) {
  saving.value = true
  try {
    await userAPI.unbindCluster({
      user_id: bindUser.value.id,
      cluster_id: clusterId
    })
    ElMessage.success(t('user.clusterUnbound'))
    await loadBindings(bindUser.value.id)
    fetchUsers()
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToUnbind'))
  } finally {
    saving.value = false
  }
}

async function toggleRegistration(val) {
  saving.value = true
  try {
    await userAPI.setRegistrationSetting({ enabled: val })
    ElMessage.success(t('user.settingUpdated'))
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('user.failedToUpdateSetting'))
    regEnabled.value = !val
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchUsers()
  fetchClusters()
  fetchSettings()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
