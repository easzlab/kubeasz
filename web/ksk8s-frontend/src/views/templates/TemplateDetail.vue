<template>
  <div>
    <div class="page-header">
      <h2>{{ templateData.name }}</h2>
      <el-tag v-if="templateData.is_default" type="success">{{ $t('cluster.default') }}</el-tag>
    </div>
    <p class="description">{{ templateData.description }}</p>

    <el-radio-group v-model="mode" style="margin-bottom: 16px">
      <el-radio-button value="form">{{ $t('common.formMode') }}</el-radio-button>
      <el-radio-button value="text">{{ $t('common.textMode') }}</el-radio-button>
    </el-radio-group>

    <div v-if="mode === 'form'">
      <el-tabs v-model="formTab" type="border-card">
        <el-tab-pane :label="$t('template.hosts')" name="hosts">
          <el-form label-width="160px">
            <h4>etcd</h4>
            <div v-for="(node, idx) in configData.nodes.etcd" :key="'etcd-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" disabled />
                <el-input :model-value="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" disabled />
              </el-form-item>
            </div>
            <div v-if="!configData.nodes.etcd || configData.nodes.etcd.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>

            <h4 style="margin-top: 24px">kube_master</h4>
            <div v-for="(node, idx) in configData.nodes.kube_master" :key="'master-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" disabled />
                <el-input :model-value="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" disabled />
              </el-form-item>
            </div>
            <div v-if="!configData.nodes.kube_master || configData.nodes.kube_master.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>

            <h4 style="margin-top: 24px">kube_node</h4>
            <div v-for="(node, idx) in configData.nodes.kube_node" :key="'node-' + idx" class="node-row">
              <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" disabled />
                <el-input :model-value="node.k8s_nodename" :placeholder="$t('config.placeholderNodeName')" style="width: 200px; margin-right: 8px" disabled />
              </el-form-item>
            </div>
            <div v-if="!configData.nodes.kube_node || configData.nodes.kube_node.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>

            <h4 style="margin-top: 24px">{{ $t('template.optionalGroups') }}</h4>
            <el-collapse>
              <el-collapse-item title="harbor" name="harbor">
                <div v-for="(node, idx) in configData.nodes.harbor" :key="'harbor-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" disabled />
                  </el-form-item>
                </div>
                <div v-if="!configData.nodes.harbor || configData.nodes.harbor.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>
              </el-collapse-item>
              <el-collapse-item title="ex_lb" name="ex_lb">
                <div v-for="(node, idx) in configData.nodes.ex_lb" :key="'exlb-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 160px; margin-right: 8px" disabled />
                    <el-input :model-value="node.lb_role" :placeholder="$t('config.placeholderLBRole')" style="width: 120px; margin-right: 8px" disabled />
                    <el-input :model-value="node.ex_apiserver_vip" :placeholder="$t('config.placeholderVIP')" style="width: 140px; margin-right: 8px" disabled />
                    <el-input :model-value="node.ex_apiserver_port" :placeholder="$t('config.placeholderPort')" style="width: 100px; margin-right: 8px" disabled />
                  </el-form-item>
                </div>
                <div v-if="!configData.nodes.ex_lb || configData.nodes.ex_lb.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>
              </el-collapse-item>
              <el-collapse-item title="chrony" name="chrony">
                <div v-for="(node, idx) in configData.nodes.chrony" :key="'chrony-' + idx" class="node-row">
                  <el-form-item :label="$t('config.nodeLabel', { n: idx + 1 })">
                    <el-input :model-value="node.ip_address" :placeholder="$t('config.placeholderIP')" style="width: 200px; margin-right: 8px" disabled />
                  </el-form-item>
                </div>
                <div v-if="!configData.nodes.chrony || configData.nodes.chrony.length === 0" class="empty-hint">{{ $t('template.noNodes') }}</div>
              </el-collapse-item>
            </el-collapse>

            <h4 style="margin-top: 24px">{{ $t('template.globalVars') }}</h4>
            <el-form-item :label="$t('config.securePort')">
              <el-input :model-value="configData.global_vars.secure_port" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.containerRuntime')">
              <el-input :model-value="configData.global_vars.container_runtime" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.networkPlugin')">
              <el-input :model-value="configData.global_vars.cluster_network" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.proxyMode')">
              <el-input :model-value="configData.global_vars.proxy_mode" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.serviceCIDR')">
              <el-input :model-value="configData.global_vars.service_cidr" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.clusterCIDR')">
              <el-input :model-value="configData.global_vars.cluster_cidr" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.nodePortRange')">
              <el-input :model-value="configData.global_vars.node_port_range" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.dnsDomain')">
              <el-input :model-value="configData.global_vars.cluster_dns_domain" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.binDir')">
              <el-input :model-value="configData.global_vars.bin_dir" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.baseDir')">
              <el-input :model-value="configData.global_vars.base_dir" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.clusterDir')">
              <el-input :model-value="configData.global_vars.cluster_dir" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.caDir')">
              <el-input :model-value="configData.global_vars.ca_dir" disabled />
            </el-form-item>
            <el-form-item label="k8s_nodename">
              <el-input :model-value="configData.global_vars.k8s_nodename" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.pythonInterpreter')">
              <el-input :model-value="configData.global_vars.ansible_python_interpreter" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.ansibleUser')">
              <el-input :model-value="configData.global_vars.ansible_user" disabled />
            </el-form-item>
            <el-form-item :label="$t('config.ansibleBecome')">
              <el-input :model-value="configData.global_vars.ansible_become" disabled />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('template.configYml')" name="config">
          <el-form label-width="200px">
            <h4>{{ $t('template.scalarParams') }}</h4>
            <div v-for="(val, key) in configData.params" :key="'param-' + key" class="param-row">
              <el-form-item :label="key">
                <el-input :model-value="val" style="width: 400px" disabled />
              </el-form-item>
            </div>
            <div v-if="!configData.params || Object.keys(configData.params).length === 0" class="empty-hint">{{ $t('template.noScalarParams') }}</div>

            <h4 style="margin-top: 24px">{{ $t('template.listParams') }}</h4>
            <div v-for="(items, key) in configData.param_lists" :key="'list-' + key" class="param-row">
              <el-form-item :label="key">
                <div v-for="(item, idx) in items" :key="key + '-' + idx" style="margin-bottom: 4px">
                  <el-input :model-value="item" style="width: 400px" disabled />
                </div>
              </el-form-item>
            </div>
            <div v-if="!configData.param_lists || Object.keys(configData.param_lists).length === 0" class="empty-hint">{{ $t('template.noListParams') }}</div>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div v-else>
      <el-tabs v-model="textTab">
        <el-tab-pane label="hosts" name="hosts">
          <el-input v-model="hostsText" type="textarea" :rows="25" disabled />
        </el-tab-pane>
        <el-tab-pane label="config.yml" name="config">
          <el-input v-model="configText" type="textarea" :rows="25" disabled />
        </el-tab-pane>
      </el-tabs>
    </div>

    <div style="margin-top: 16px">
      <el-button @click="$router.push('/templates')">{{ $t('template.back') }}</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { templateAPI } from '../../api/client'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()
const mode = ref('form')
const formTab = ref('hosts')
const textTab = ref('hosts')

const templateData = ref({
  name: '',
  description: '',
  is_default: false
})

const configData = ref({
  nodes: {
    etcd: [],
    kube_master: [],
    kube_node: [],
    harbor: [],
    ex_lb: [],
    chrony: []
  },
  global_vars: {},
  params: {},
  param_lists: {}
})

const hostsText = ref('')
const configText = ref('')

async function loadTemplate() {
  try {
    const res = await templateAPI.getParsed(route.params.id)
    const data = res.data

    templateData.value = {
      name: data.name,
      description: data.description,
      is_default: data.is_default
    }

    const groups = ['etcd', 'kube_master', 'kube_node', 'harbor', 'ex_lb', 'chrony']
    for (const g of groups) {
      configData.value.nodes[g] = (data.nodes[g] || []).map(n => ({
        ip_address: n.ip_address || '',
        k8s_nodename: n.k8s_nodename || '',
        new_install: n.new_install || false,
        lb_role: n.lb_role || '',
        ex_apiserver_vip: n.ex_apiserver_vip || '',
        ex_apiserver_port: n.ex_apiserver_port || ''
      }))
    }

    if (data.global_vars) {
      configData.value.global_vars = data.global_vars
    }
    if (data.params) {
      configData.value.params = { ...data.params }
    }
    if (data.param_lists) {
      configData.value.param_lists = {}
      for (const [k, v] of Object.entries(data.param_lists)) {
        configData.value.param_lists[k] = [...v]
      }
    }

    hostsText.value = data.hosts_content || ''
    configText.value = data.config_content || ''
  } catch (err) {
    console.error(t('template.failedToLoad'), err)
    ElMessage.error(t('template.failedToLoad'))
  }
}

onMounted(loadTemplate)
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.description {
  color: #606266;
  margin-bottom: 20px;
}
.node-row {
  margin-bottom: 4px;
}
.param-row {
  margin-bottom: 4px;
}
.empty-hint {
  color: #909399;
  font-style: italic;
  margin: 8px 0;
}
</style>
