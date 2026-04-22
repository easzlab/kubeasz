<template>
  <div class="log-viewer">
    <div class="log-header">
      <span class="log-title">{{ $t('install.liveLogs') }}</span>
      <el-tag size="small" :type="wsConnected ? 'success' : 'danger'">
        {{ wsConnected ? $t('ssh.connected') : $t('ssh.disconnected') }}
      </el-tag>
      <el-button size="small" :icon="CopyDocument" @click="copyLogs">{{ $t('common.copy') }}</el-button>
      <el-button size="small" @click="clearLogs">{{ $t('common.clear') }}</el-button>
    </div>
    <div ref="logContainer" class="log-container">
      <div
        v-for="line in logs"
        :key="line.line_number"
        :class="['log-line', line.stream === 'stderr' ? 'log-stderr' : 'log-stdout']"
      >
        <span class="line-num">{{ line.line_number }}</span>
        <span class="line-content">{{ line.content }}</span>
      </div>
      <div v-if="logs.length === 0" class="log-empty">{{ $t('install.noLogs') }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { taskAPI } from '../api/client'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  taskId: { type: Number, required: true },
  clusterId: { type: String, required: true }
})

const logs = ref([])
const logContainer = ref(null)
const wsConnected = ref(false)
const offset = ref(0)
let ws = null
let reconnectTimer = null

async function connectWS() {
  if (ws) {
    ws.close()
  }

  // Fetch historical logs from DB first (in case ring buffer was lost on backend restart)
  try {
    const res = await taskAPI.logs(props.clusterId, props.taskId, offset.value, 1000)
    if (res.data && res.data.length > 0) {
      for (const line of res.data) {
        logs.value.push(line)
        offset.value = Math.max(offset.value, line.line_number)
      }
      scrollToBottom()
    }
  } catch (e) {
    console.error('Failed to fetch historical logs', e)
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const url = `${protocol}//${host}/ws/tasks/${props.taskId}/logs?offset=${offset.value}`

  ws = new WebSocket(url)

  ws.onopen = () => {
    wsConnected.value = true
  }

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.type === 'ready') {
        offset.value = data.total_lines
        return
      }
      if (data.line_number !== undefined) {
        logs.value.push(data)
        offset.value = data.line_number
        scrollToBottom()
      }
    } catch (e) {
      // Non-JSON line
      logs.value.push({ line_number: offset.value + 1, content: event.data, stream: 'stdout' })
      offset.value++
      scrollToBottom()
    }
  }

  ws.onclose = () => {
    wsConnected.value = false
    reconnectTimer = setTimeout(connectWS, 3000)
  }

  ws.onerror = () => {
    wsConnected.value = false
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

function clearLogs() {
  logs.value = []
  offset.value = 0
}

async function copyLogs() {
  const text = logs.value.map(l => l.content).join('\n')
  if (!text) {
    ElMessage.warning(t('logs.noLogsToCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('logs.logsCopied'))
  } catch (e) {
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success(t('logs.logsCopied'))
  }
}

watch(() => props.taskId, () => {
  logs.value = []
  offset.value = 0
  connectWS()
})

onMounted(() => {
  connectWS()
})

onUnmounted(() => {
  if (ws) ws.close()
  if (reconnectTimer) clearTimeout(reconnectTimer)
})
</script>

<style scoped>
.log-viewer {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  background: #1a1a2e;
  color: #c0c0c0;
}
.log-header {
  padding: 8px 12px;
  border-bottom: 1px solid #2a2a3e;
  display: flex;
  align-items: center;
  gap: 12px;
}
.log-title {
  font-weight: bold;
  color: #fff;
}
.log-container {
  height: 400px;
  overflow-y: auto;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}
.log-line {
  display: flex;
  gap: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}
.line-num {
  color: #606266;
  min-width: 40px;
  text-align: right;
  user-select: none;
}
.log-stdout .line-content {
  color: #c0c0c0;
}
.log-stderr .line-content {
  color: #f56c6c;
}
.log-empty {
  color: #606266;
  text-align: center;
  padding: 40px;
}
</style>
