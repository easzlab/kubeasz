<template>
  <div>
    <h2>{{ $t('ssh.title') }}</h2>
    <p v-if="addr" style="color: #606266; margin-bottom: 12px"
      >{{ $t('ssh.connecting', { addr, user: user || 'root' }) }}</p
    >
    <p v-else style="color: #f56c6c; margin-bottom: 12px"
      >{{ $t('ssh.noTarget') }}</p
    >
    <WebSSH :addr="addr || '127.0.0.1:22'" :user="user || 'root'" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import WebSSH from '../components/WebSSH.vue'

const route = useRoute()
const addr = computed(() => route.query.addr || '')
const user = computed(() => route.query.user || 'root')
</script>
