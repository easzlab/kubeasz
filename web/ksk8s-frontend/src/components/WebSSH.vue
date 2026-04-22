<template>
  <div class="webssh-container">
    <div class="ssh-header">
      <span class="ssh-title">{{ $t('ssh.webssh') }}</span>
      <el-tag size="small" :type="connected ? 'success' : 'danger'">
        {{ connected ? $t('ssh.connected') : $t('ssh.disconnected') }}
      </el-tag>
      <el-button size="small" @click="connect" v-if="!connected">{{ $t('ssh.connect') }}</el-button>
      <el-button size="small" @click="disconnect" v-else>{{ $t('ssh.disconnect') }}</el-button>
    </div>
    <div ref="terminalRef" class="terminal"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

const props = defineProps({
  addr: { type: String, default: '127.0.0.1:22' },
  user: { type: String, default: 'root' }
})

const terminalRef = ref(null)
const connected = ref(false)
let term = null
let fitAddon = null
let ws = null

function initTerminal() {
  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'monospace',
    theme: {
      background: '#1a1a2e',
      foreground: '#c0c0c0'
    }
  })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(terminalRef.value)
  fitAddon.fit()

  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })
}

function connect() {
  if (ws) ws.close()

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const url = `${protocol}//${host}/ws/ssh?addr=${encodeURIComponent(props.addr)}&user=${encodeURIComponent(props.user)}`

  ws = new WebSocket(url)
  ws.binaryType = 'arraybuffer'

  ws.onopen = () => {
    connected.value = true
    term.writeln('\r\n\x1b[32mConnected\x1b[0m')
    sendResize()
  }

  ws.onmessage = (event) => {
    if (typeof event.data === 'string') {
      term.writeln(event.data)
    } else {
      const decoder = new TextDecoder()
      term.write(decoder.decode(event.data))
    }
  }

  ws.onclose = () => {
    connected.value = false
    term.writeln('\r\n\x1b[31mDisconnected\x1b[0m')
  }

  ws.onerror = () => {
    connected.value = false
    term.writeln('\r\n\x1b[31mConnection error\x1b[0m')
  }
}

function disconnect() {
  if (ws) ws.close()
}

function sendResize() {
  if (!ws || !fitAddon) return
  const dims = fitAddon.proposeDimensions()
  if (dims) {
    ws.send(JSON.stringify({ type: 'resize', cols: dims.cols, rows: dims.rows }))
  }
}

onMounted(() => {
  nextTick(() => {
    initTerminal()
    connect()
    window.addEventListener('resize', () => {
      if (fitAddon) fitAddon.fit()
      sendResize()
    })
  })
})

onUnmounted(() => {
  if (ws) ws.close()
  if (term) term.dispose()
})
</script>

<style scoped>
.webssh-container {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}
.ssh-header {
  padding: 8px 12px;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fff;
}
.ssh-title {
  font-weight: bold;
}
.terminal {
  height: 500px;
  padding: 4px;
  background: #1a1a2e;
}
</style>
