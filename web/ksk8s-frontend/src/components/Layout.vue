<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <span class="logo-text">ksk8s</span>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        class="sidebar-menu"
        background-color="#1a1a2e"
        text-color="#b0b0c3"
        active-text-color="#409eff"
      >
        <el-menu-item index="/clusters">
          <el-icon><Grid /></el-icon>
          <span>{{ $t('nav.clusters') }}</span>
        </el-menu-item>
        <el-menu-item index="/templates">
          <el-icon><Document /></el-icon>
          <span>{{ $t('nav.templates') }}</span>
        </el-menu-item>
        <el-menu-item index="/audit-logs">
          <el-icon><List /></el-icon>
          <span>{{ $t('nav.auditLogs') }}</span>
        </el-menu-item>
        <el-menu-item v-if="auth.isPlatformAdmin" index="/users">
          <el-icon><UserFilled /></el-icon>
          <span>{{ $t('nav.userManagement') }}</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>{{ $t('nav.settings') }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-right">
          <el-select v-model="locale" size="small" style="width: 110px">
            <el-option label="English" value="en" />
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="繁體中文" value="zh-TW" />
            <el-option label="Français" value="fr" />
          </el-select>
          <el-tag size="small" :type="roleTagType">{{ roleLabel }}</el-tag>
          <span class="username">{{ auth.user?.username }}</span>
          <el-button size="small" @click="logout">{{ $t('auth.logout') }}</el-button>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { Grid, Document, List, UserFilled, Setting } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()
const { t, locale } = useI18n()

function normalizeRole(role) {
  if (role === 'admin') return 'platform_admin'
  if (role === 'viewer') return 'security_auditor'
  return role
}

const roleLabel = computed(() => {
  const map = {
    platform_admin: t('role.platform_admin'),
    cluster_admin: t('role.cluster_admin'),
    security_auditor: t('role.security_auditor')
  }
  return map[normalizeRole(auth.user?.role)] || auth.user?.role
})

const roleTagType = computed(() => {
  const map = {
    platform_admin: 'danger',
    cluster_admin: 'warning',
    security_auditor: 'info'
  }
  return map[normalizeRole(auth.user?.role)] || ''
})

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.sidebar {
  background: #1a1a2e;
  color: #fff;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #2a2a3e;
}
.logo-text {
  font-size: 22px;
  font-weight: bold;
  color: #409eff;
}
.sidebar-menu {
  border-right: none;
}
.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.username {
  color: #606266;
}
.main {
  background: #f5f7fa;
  padding: 20px;
}
</style>
