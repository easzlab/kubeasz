import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('../components/Layout.vue'),
    redirect: '/clusters',
    children: [
      {
        path: 'clusters',
        name: 'ClusterList',
        component: () => import('../views/clusters/ClusterList.vue')
      },
      {
        path: 'clusters/create',
        name: 'ClusterCreate',
        component: () => import('../views/clusters/ClusterCreate.vue')
      },
      {
        path: 'clusters/:id',
        name: 'ClusterDetail',
        component: () => import('../views/clusters/ClusterDetail.vue'),
        redirect: to => `/clusters/${to.params.id}/install`,
        children: [
          {
            path: 'install',
            name: 'InstallPipeline',
            component: () => import('../views/clusters/InstallPipeline.vue')
          },
          {
            path: 'lifecycle',
            name: 'LifecycleOps',
            component: () => import('../views/clusters/LifecycleOps.vue')
          },
          {
            path: 'nodes',
            name: 'NodeManagement',
            component: () => import('../views/clusters/NodeManagement.vue')
          },
          {
            path: 'config',
            name: 'ClusterConfig',
            component: () => import('../views/clusters/ClusterConfig.vue')
          },
          {
            path: 'security',
            name: 'SecurityOps',
            component: () => import('../views/clusters/SecurityOps.vue')
          },
          {
            path: 'logs',
            name: 'ClusterLogs',
            component: () => import('../views/clusters/ClusterLogs.vue')
          }
        ]
      },
      {
        path: 'templates',
        name: 'TemplateList',
        component: () => import('../views/templates/TemplateList.vue')
      },
      {
        path: 'templates/:id',
        name: 'TemplateDetail',
        component: () => import('../views/templates/TemplateDetail.vue')
      },
      {
        path: 'audit-logs',
        name: 'AuditLogs',
        component: () => import('../views/AuditLogs.vue')
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('../views/users/UserManagement.vue'),
        meta: { requirePlatformAdmin: true }
      },
      {
        path: 'settings',
        name: 'UserSettings',
        component: () => import('../views/users/UserSettings.vue')
      }
    ]
  },
  {
    path: '/webssh',
    name: 'WebSSH',
    component: () => import('../views/WebSSHPage.vue'),
    meta: { public: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (!to.meta?.public && !auth.isLoggedIn) {
    next('/login')
    return
  }
  if (to.path === '/login' && auth.isLoggedIn) {
    next('/')
    return
  }
  if (to.meta?.requirePlatformAdmin && !auth.isPlatformAdmin) {
    next('/')
    return
  }
  next()
})

export default router
