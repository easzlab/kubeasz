<template>
  <div>
    <el-radio-group v-model="mode" style="margin-bottom: 16px">
      <el-radio-button value="form">{{ $t('config.formMode') }}</el-radio-button>
      <el-radio-button value="text">{{ $t('config.textMode') }}</el-radio-button>
    </el-radio-group>

    <el-alert
      v-if="dirty"
      :title="$t('common.unsavedChanges')"
      type="warning"
      :closable="false"
      style="margin-bottom: 16px"
    />

    <div v-if="mode === 'form'">
      <el-tabs v-model="formTab" type="border-card">
        <el-tab-pane :label="$t('config.hosts')" name="hosts">
          <el-form :model="configData" label-width="160px">
            <!-- etcd -->
            <h4>{{ $t('config.etcd') }} <span style="color: #f56c6c">*</span></h4>
            <div v-for="(node, idx) in configData.nodes.etcd" :key="'etcd-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" />
                <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" />
                <el-button type="danger" size="small" @click="removeNode('etcd', idx)">{{ $t('common.delete') }}</el-button>
              </el-form-item>
            </div>
            <el-button type="primary" size="small" @click="addNode('etcd')">{{ $t('config.addNode', { group: $t('config.etcd') }) }}</el-button>

            <!-- kube_master -->
            <h4 style="margin-top: 24px">{{ $t('config.kubeMaster') }} <span style="color: #f56c6c">*</span></h4>
            <div v-for="(node, idx) in configData.nodes.kube_master" :key="'master-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" />
                <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" />
                <el-button type="danger" size="small" @click="removeNode('kube_master', idx)">{{ $t('common.delete') }}</el-button>
              </el-form-item>
            </div>
            <el-button type="primary" size="small" @click="addNode('kube_master')">{{ $t('config.addNode', { group: $t('config.kubeMaster') }) }}</el-button>

            <!-- kube_node -->
            <h4 style="margin-top: 24px">{{ $t('config.kubeNode') }}</h4>
            <div v-for="(node, idx) in configData.nodes.kube_node" :key="'node-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" />
                <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" />
                <el-button type="danger" size="small" @click="removeNode('kube_node', idx)">{{ $t('common.delete') }}</el-button>
              </el-form-item>
            </div>
            <el-button type="primary" size="small" @click="addNode('kube_node')">{{ $t('config.addNode', { group: $t('config.kubeNode') }) }}</el-button>

            <!-- Optional groups -->
            <h4 style="margin-top: 24px">{{ $t('config.optionalGroups') }}</h4>
            <el-collapse>
              <el-collapse-item title="harbor" name="harbor">
                <div v-for="(node, idx) in configData.nodes.harbor" :key="'harbor-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" />
                    <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" />
                    <el-button type="danger" size="small" @click="removeNode('harbor', idx)">{{ $t('common.delete') }}</el-button>
                  </el-form-item>
                </div>
                <el-button type="primary" size="small" @click="addNode('harbor')">{{ $t('config.addNode', { group: 'harbor' }) }}</el-button>
              </el-collapse-item>
              <el-collapse-item title="ex_lb" name="ex_lb">
                <div v-for="(node, idx) in configData.nodes.ex_lb" :key="'exlb-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 160px; margin-right: 8px" />
                    <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 160px; margin-right: 8px" />
                    <el-input v-model="node.lb_role" :placeholder="$t('config.placeholderLBRole')" style="width: 120px; margin-right: 8px" />
                    <el-input v-model="node.ex_apiserver_vip" :placeholder="$t('config.placeholderVIP')" style="width: 140px; margin-right: 8px" />
                    <el-input v-model="node.ex_apiserver_port" :placeholder="$t('config.placeholderPort')" style="width: 100px; margin-right: 8px" />
                    <el-button type="danger" size="small" @click="removeNode('ex_lb', idx)">{{ $t('common.delete') }}</el-button>
                  </el-form-item>
                </div>
                <el-button type="primary" size="small" @click="addNode('ex_lb')">{{ $t('config.addNode', { group: 'ex_lb' }) }}</el-button>
              </el-collapse-item>
              <el-collapse-item title="chrony" name="chrony">
                <div v-for="(node, idx) in configData.nodes.chrony" :key="'chrony-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input v-model="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" />
                    <el-input v-model="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" />
                    <el-button type="danger" size="small" @click="removeNode('chrony', idx)">{{ $t('common.delete') }}</el-button>
                  </el-form-item>
                </div>
                <el-button type="primary" size="small" @click="addNode('chrony')">{{ $t('config.addNode', { group: 'chrony' }) }}</el-button>
              </el-collapse-item>
            </el-collapse>

            <!-- Global Variables -->
            <h4 style="margin-top: 24px">{{ $t('config.globalVars') }}</h4>
            <div v-for="meta in globalVarMeta" :key="'gvar-' + meta.key" class="param-row">
              <el-form-item>
                <template #label>
                  <span>{{ meta.key }}</span>
                  <el-tooltip v-if="meta.comment" :content="meta.comment" placement="top">
                    <el-icon style="margin-left: 4px; cursor: pointer; color: #909399"><QuestionFilled /></el-icon>
                  </el-tooltip>
                </template>
                <el-input v-model="configData.global_vars[hostsVarKeyToProp(meta.key)]" />
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('config.configYml')" name="config">
          <el-form :model="configData" label-width="240px">
            <!-- Parameters in template order with hints -->
            <div v-for="meta in paramMeta" :key="'meta-' + meta.key" class="param-row">
              <!-- Scalar param -->
              <el-form-item v-if="meta.type === 'scalar'">
                <template #label>
                  <span>{{ meta.key }}{{ isRequiredParam(meta.key) ? ' *' : '' }}</span>
                  <el-tooltip v-if="meta.comment" :content="meta.comment" placement="top">
                    <el-icon style="margin-left: 4px; cursor: pointer; color: #909399"><QuestionFilled /></el-icon>
                  </el-tooltip>
                </template>
                <el-input v-model="configData.params[meta.key]" style="width: 400px; margin-right: 8px" />
                <el-button type="danger" size="small" @click="deleteParam(meta.key)">{{ $t('common.delete') }}</el-button>
              </el-form-item>

              <!-- List param -->
              <el-form-item v-else>
                <template #label>
                  <span>{{ meta.key }}{{ isRequiredList(meta.key) ? ' *' : '' }}</span>
                  <el-tooltip v-if="meta.comment" :content="meta.comment" placement="top">
                    <el-icon style="margin-left: 4px; cursor: pointer; color: #909399"><QuestionFilled /></el-icon>
                  </el-tooltip>
                </template>
                <div v-for="(item, idx) in configData.param_lists[meta.key]" :key="meta.key + '-' + idx" style="margin-bottom: 4px">
                  <el-input v-model="configData.param_lists[meta.key][idx]" style="width: 400px; margin-right: 8px" />
                  <el-button type="danger" size="small" @click="removeListItem(meta.key, idx)">{{ $t('common.delete') }}</el-button>
                </div>
                <el-button type="primary" size="small" @click="addListItem(meta.key)">{{ $t('config.addNode', { group: $t('config.paramTypeList') }) }}</el-button>
                <el-button type="danger" size="small" @click="deleteListParam(meta.key)">{{ $t('common.delete') }} {{ $t('config.paramName') }}</el-button>
              </el-form-item>
            </div>

            <!-- Extra params not in template meta (user-added) -->
            <div v-for="(val, key) in extraScalarParams" :key="'extra-param-' + key" class="param-row">
              <el-form-item>
                <template #label>
                  <span>{{ key }}{{ isRequiredParam(key) ? ' *' : '' }}</span>
                </template>
                <el-input v-model="configData.params[key]" style="width: 400px; margin-right: 8px" />
                <el-button type="danger" size="small" @click="deleteParam(key)">{{ $t('common.delete') }}</el-button>
              </el-form-item>
            </div>
            <div v-for="(items, key) in extraListParams" :key="'extra-list-' + key" class="param-row">
              <el-form-item>
                <template #label>
                  <span>{{ key }}{{ isRequiredList(key) ? ' *' : '' }}</span>
                </template>
                <div v-for="(item, idx) in items" :key="key + '-extra-' + idx" style="margin-bottom: 4px">
                  <el-input v-model="configData.param_lists[key][idx]" style="width: 400px; margin-right: 8px" />
                  <el-button type="danger" size="small" @click="removeListItem(key, idx)">{{ $t('common.delete') }}</el-button>
                </div>
                <el-button type="primary" size="small" @click="addListItem(key)">{{ $t('config.addNode', { group: $t('config.paramTypeList') }) }}</el-button>
                <el-button type="danger" size="small" @click="deleteListParam(key)">{{ $t('common.delete') }} {{ $t('config.paramName') }}</el-button>
              </el-form-item>
            </div>

            <!-- Add new param -->
            <h4 style="margin-top: 24px">{{ $t('config.addParam') }}</h4>
            <el-form-item :label="$t('config.paramName')">
              <el-input v-model="newParamName" placeholder="e.g. K8S_VER" style="width: 200px; margin-right: 8px" />
              <el-radio-group v-model="newParamType">
                <el-radio-button value="scalar">{{ $t('config.paramTypeScalar') }}</el-radio-button>
                <el-radio-button value="list">{{ $t('config.paramTypeList') }}</el-radio-button>
              </el-radio-group>
              <el-button type="primary" size="small" @click="addNewParam" style="margin-left: 8px">{{ $t('config.addParam') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div style="margin-top: 16px">
        <el-button type="primary" @click="saveConfig" :loading="saving">{{ $t('common.save') }}</el-button>
        <el-button @click="previewConfig">{{ $t('config.preview') }}</el-button>
      </div>
    </div>

    <div v-else>
      <el-tabs v-model="textTab">
        <el-tab-pane label="hosts" name="hosts">
          <el-input v-model="hostsText" type="textarea" :rows="25" />
        </el-tab-pane>
        <el-tab-pane label="config.yml" name="config">
          <el-input v-model="configText" type="textarea" :rows="25" />
        </el-tab-pane>
      </el-tabs>
      <el-button type="primary" @click="saveTextConfig" :loading="saving" style="margin-top: 12px">{{ $t('common.save') }}</el-button>
    </div>

    <el-dialog v-model="previewVisible" :title="$t('config.preview')" width="800px">
      <el-tabs>
        <el-tab-pane label="hosts">
          <pre class="config-preview">{{ preview.hosts }}</pre>
        </el-tab-pane>
        <el-tab-pane label="config.yml">
          <pre class="config-preview">{{ preview.config }}</pre>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { clusterAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const mode = ref('form')
const formTab = ref('hosts')
const textTab = ref('hosts')
const dirty = ref(false)
const saving = ref(false)
const previewVisible = ref(false)
const preview = ref({ hosts: '', config: '' })

const newParamName = ref('')
const newParamType = ref('scalar')

// Metadata from template for display order and hints
const paramMeta = ref([])
const globalVarMeta = ref([])

function makeEmptyNode() {
  return { ip_address: '', k8s_nodename: '', new_install: false, lb_role: '', ex_apiserver_vip: '', ex_apiserver_port: '' }
}

function makeDefaultGlobalVars() {
  return {
    secure_port: '6443',
    container_runtime: 'containerd',
    cluster_network: 'calico',
    proxy_mode: 'ipvs',
    service_cidr: '10.68.0.0/16',
    cluster_cidr: '172.20.0.0/16',
    node_port_range: '30000-32767',
    cluster_dns_domain: 'cluster.local',
    bin_dir: '/opt/kube/bin',
    base_dir: '/etc/kubeasz',
    cluster_dir: '',
    ca_dir: '/etc/kubernetes/ssl',
    k8s_nodename: '',
    ansible_python_interpreter: '/usr/bin/python3',
    ansible_user: 'root',
    ansible_become: 'no'
  }
}

const configData = ref({
  nodes: {
    etcd: [],
    kube_master: [],
    kube_node: [],
    harbor: [],
    ex_lb: [],
    chrony: []
  },
  global_vars: makeDefaultGlobalVars(),
  params: {},
  param_lists: {}
})

const hostsText = ref('')
const configText = ref('')

const requiredParams = ['nfs_server', 'nfs_path']
const requiredLists = ['INSECURE_REG', 'MASTER_CERT_HOSTS']

function isRequiredParam(key) {
  return requiredParams.includes(key)
}

function isRequiredList(key) {
  return requiredLists.includes(key)
}

function hostsVarKeyToProp(key) {
  const map = {
    'SECURE_PORT': 'secure_port',
    'CONTAINER_RUNTIME': 'container_runtime',
    'CLUSTER_NETWORK': 'cluster_network',
    'PROXY_MODE': 'proxy_mode',
    'SERVICE_CIDR': 'service_cidr',
    'CLUSTER_CIDR': 'cluster_cidr',
    'NODE_PORT_RANGE': 'node_port_range',
    'CLUSTER_DNS_DOMAIN': 'cluster_dns_domain',
    'bin_dir': 'bin_dir',
    'base_dir': 'base_dir',
    'cluster_dir': 'cluster_dir',
    'ca_dir': 'ca_dir',
    'k8s_nodename': 'k8s_nodename',
    'ansible_python_interpreter': 'ansible_python_interpreter',
    'ansible_user': 'ansible_user',
    'ansible_become': 'ansible_become',
  }
  return map[key] || key.toLowerCase()
}

// Extra params not present in template meta
const extraScalarParams = computed(() => {
  const metaKeys = new Set(paramMeta.value.filter(m => m.type === 'scalar').map(m => m.key))
  const result = {}
  for (const [k, v] of Object.entries(configData.value.params)) {
    if (!metaKeys.has(k)) {
      result[k] = v
    }
  }
  return result
})

const extraListParams = computed(() => {
  const metaKeys = new Set(paramMeta.value.filter(m => m.type === 'list').map(m => m.key))
  const result = {}
  for (const [k, v] of Object.entries(configData.value.param_lists)) {
    if (!metaKeys.has(k)) {
      result[k] = v
    }
  }
  return result
})

function addNode(group) {
  configData.value.nodes[group].push(makeEmptyNode())
  dirty.value = true
}

function removeNode(group, idx) {
  configData.value.nodes[group].splice(idx, 1)
  dirty.value = true
}

function deleteParam(key) {
  delete configData.value.params[key]
  dirty.value = true
}

function deleteListParam(key) {
  delete configData.value.param_lists[key]
  dirty.value = true
}

function addListItem(key) {
  if (!configData.value.param_lists[key]) {
    configData.value.param_lists[key] = []
  }
  configData.value.param_lists[key].push('')
  dirty.value = true
}

function removeListItem(key, idx) {
  configData.value.param_lists[key].splice(idx, 1)
  dirty.value = true
}

function addNewParam() {
  const name = newParamName.value.trim()
  if (!name) {
    ElMessage.warning(t('config.pleaseEnterParamName'))
    return
  }
  if (newParamType.value === 'scalar') {
    configData.value.params[name] = ''
  } else {
    if (!configData.value.param_lists[name]) {
      configData.value.param_lists[name] = []
    }
    configData.value.param_lists[name].push('')
  }
  newParamName.value = ''
  dirty.value = true
}

async function loadConfig() {
  try {
    const res = await clusterAPI.getConfig(route.params.id)
    const data = res.data

    // Merge nodes - ignore new_install from template/DB, default to false
    const groups = ['etcd', 'kube_master', 'kube_node', 'harbor', 'ex_lb', 'chrony']
    for (const g of groups) {
      configData.value.nodes[g] = (data.nodes[g] || []).map(n => ({
        ip_address: n.ip_address || '',
        k8s_nodename: n.k8s_nodename || '',
        new_install: false,
        lb_role: n.lb_role || '',
        ex_apiserver_vip: n.ex_apiserver_vip || '',
        ex_apiserver_port: n.ex_apiserver_port || ''
      }))
    }

    // Merge global vars
    if (data.global_vars) {
      configData.value.global_vars = { ...makeDefaultGlobalVars(), ...data.global_vars }
    }

    // Merge params
    if (data.params) {
      configData.value.params = { ...data.params }
    }

    // Merge param lists
    if (data.param_lists) {
      configData.value.param_lists = {}
      for (const [k, v] of Object.entries(data.param_lists)) {
        configData.value.param_lists[k] = [...v]
      }
    }

    // Load metadata for display order and hints
    paramMeta.value = data.param_meta || []
    globalVarMeta.value = data.global_var_meta || []

    // Text mode content
    if (data.hosts_content) {
      hostsText.value = data.hosts_content
    }
    if (data.config_content) {
      configText.value = data.config_content
    }
  } catch (err) {
    console.error('Failed to load config', err)
    ElMessage.error(t('config.failedToLoad'))
  }
}

async function saveConfig() {
  // Client-side validation
  if (configData.value.nodes.etcd.length === 0 || configData.value.nodes.etcd.some(n => !n.ip_address)) {
    ElMessage.error(t('error.etcd nodes are required'))
    return
  }
  if (configData.value.nodes.kube_master.length === 0 || configData.value.nodes.kube_master.some(n => !n.ip_address)) {
    ElMessage.error(t('error.kube_master nodes are required'))
    return
  }
  if (!configData.value.param_lists.INSECURE_REG || configData.value.param_lists.INSECURE_REG.length === 0 || configData.value.param_lists.INSECURE_REG.some(v => !v)) {
    ElMessage.error(t('error.INSECURE_REG is required'))
    return
  }
  if (!configData.value.param_lists.MASTER_CERT_HOSTS || configData.value.param_lists.MASTER_CERT_HOSTS.length === 0 || configData.value.param_lists.MASTER_CERT_HOSTS.some(v => !v)) {
    ElMessage.error(t('error.MASTER_CERT_HOSTS is required'))
    return
  }
  if (!configData.value.params.nfs_server) {
    ElMessage.error(t('error.nfs_server is required'))
    return
  }
  if (!configData.value.params.nfs_path) {
    ElMessage.error(t('error.nfs_path is required'))
    return
  }

  saving.value = true
  try {
    // Build payload with only non-empty optional nodes, new_install always false
    const payload = {
      nodes: {},
      global_vars: configData.value.global_vars,
      params: configData.value.params,
      param_lists: configData.value.param_lists
    }

    const groups = ['etcd', 'kube_master', 'kube_node', 'harbor', 'ex_lb', 'chrony']
    for (const g of groups) {
      payload.nodes[g] = configData.value.nodes[g]
        .filter(n => n.ip_address.trim() !== '')
        .map(n => ({
          ip_address: n.ip_address,
          k8s_nodename: n.k8s_nodename,
          lb_role: n.lb_role || null,
          ex_apiserver_vip: n.ex_apiserver_vip || null,
          ex_apiserver_port: n.ex_apiserver_port || null,
          new_install: false
        }))
    }

    await clusterAPI.saveConfig(route.params.id, payload)
    ElMessage.success(t('config.configSaved'))

    // Reload to get regenerated text
    await loadConfig()
    dirty.value = false
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('config.failedToSave'))
  } finally {
    saving.value = false
  }
}

async function saveTextConfig() {
  saving.value = true
  try {
    await clusterAPI.update(route.params.id, {
      hosts_content: hostsText.value,
      config_content: configText.value
    })
    ElMessage.success(t('config.configSaved'))
    dirty.value = false
  } catch (err) {
    ElMessage.error(err.response?.data?.error ? t(`error.${err.response.data.error}`) : t('config.failedToSave'))
  } finally {
    saving.value = false
  }
}

async function previewConfig() {
  try {
    const res = await clusterAPI.generateConfig(route.params.id)
    preview.value = res.data
    previewVisible.value = true
  } catch (err) {
    ElMessage.error(t('config.failedToGenerate'))
  }
}

watch(configData, () => {
  dirty.value = true
}, { deep: true })

watch(hostsText, () => { dirty.value = true })
watch(configText, () => { dirty.value = true })

onMounted(loadConfig)
</script>

<style scoped>
.node-row {
  margin-bottom: 4px;
}
.param-row {
  margin-bottom: 4px;
}
.config-preview {
  background: #1a1a2e;
  color: #c0c0c0;
  padding: 16px;
  border-radius: 4px;
  overflow-x: auto;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}
</style>
