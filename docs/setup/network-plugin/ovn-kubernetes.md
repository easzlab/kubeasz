## 06-安装ovn-kubernetes网络组件.md

ovn-kubernetes是一款开源的OVN的网络插件

- 项目地址 https://github.com/ovn-org/ovn-kubernetes

### kubeasz 集成安装 ovn-kubernetes

ovn-kubernetes 的安装十分简单，详见项目的安装文档；基于 kubeasz，以下步骤将安装一个集成了 ovn-kubernetes 网络的 k8s 集群；

- 在 ansible hosts 中设置变量 `CLUSTER_NETWORK="ovn-kubernetes"`
- 执行安装 `ansible-playbook 90.setup.yml` 或者 `ezctl setup`

kubeasz 项目为`ovn-kubernetes`网络生成的 ansible role 如下：

``` bash
roles/ovn-kubernetes
├── files
│   ├── egressfirewall.yaml
│   └── egressip.yaml
├── tasks
│   └── main.yml # 安装执行文件
├── templates
│   ├── ovn-setup.yaml.j2                    # ovn-kubernetes yaml 模板   
│   ├── ovnkube-db-raft.yaml.j2              # ovn-kubernetes yaml 模板   
│   ├── ovnkube-db.yaml.j2                   # ovn-kubernetes yaml 模板   
│   ├── ovnkube-master.yaml.j2               # ovn-kubernetes yaml 模板   
│   ├── ovnkube-node-smart-nic-host.yaml.j2  # ovn-kubernetes yaml 模板   
│   ├── ovnkube-node.yaml.j2                 # ovn-kubernetes yaml 模板   
│   └── ovs-node.yaml.j2                     # ovs yaml 模板
└── vars
    └── main.yml

```

安装成功后，可以验证所有 k8s 集群功能正常，查看集群的 pod 网络如下：

```
$ kubectl get pods --all-namespaces -o wide
NAMESPACE        NAME                                         READY   STATUS    RESTARTS   AGE     IP               NODE             NOMINATED NODE   READINESS GATES
kube-system      coredns-74c56d8f8d-j4mvs                     1/1     Running   0          9m33s   172.20.0.3       192.168.200.62   <none>           <none>
kube-system      dashboard-metrics-scraper-856586f554-bc9mj   1/1     Running   0          9m12s   172.20.0.6       192.168.200.62   <none>           <none>
kube-system      kubernetes-dashboard-c4ff5556c-rg9wf         1/1     Running   0          9m13s   172.20.0.5       192.168.200.62   <none>           <none>
kube-system      metrics-server-8568cf894b-rxmrw              1/1     Running   0          9m26s   172.20.0.4       192.168.200.62   <none>           <none>
kube-system      node-local-dns-ndn4n                         1/1     Running   0          9m32s   192.168.200.62   192.168.200.62   <none>           <none>
ovn-kubernetes   ovnkube-db-7f54794f95-9rrq9                  2/2     Running   0          11m     192.168.200.62   192.168.200.62   <none>           <none>
ovn-kubernetes   ovnkube-master-76d755dbc9-txjdq              3/3     Running   0          11m     192.168.200.62   192.168.200.62   <none>           <none>
ovn-kubernetes   ovnkube-node-qpfk7                           3/3     Running   0          11m     192.168.200.62   192.168.200.62   <none>           <none>
ovn-kubernetes   ovs-node-gvgzw                               1/1     Running   0          11m     192.168.200.62   192.168.200.62   <none>           <none>
```

### 启用 ovn-kubernetes 特性

编辑集群配置文件`config.yml`

#### 启用高可用模式

OVN DB将以集群模式部署

```
OVN_DB_RAFT_ENABLE: "true"
OVN_DB_NODES:
  - "{{ groups['kube_master'][0] }}"
  - "192.168.200.10"
```

#### 启用EgressIP功能

```
EGRESS_IP_ENABLE: "true"
EGRESS_IP_NODES:
  - "{{ groups['kube_node'][0] }}"
  - "192.168.200.10"
```

#### 启用EgressFirewall功能

```
EGRESS_FIREWALL_ENABLE: "true"
```

#### 启用智能网卡功能

使用智能网卡替代OVS

```
SMART_NIC_HOST_ENABLE: "true"
MGMT_PORT_NETDEV: "ens3"
SMART_NIC_HOSTS:
  - "{{ groups['kube_node'][0] }}"
  - "192.168.200.10"
```

#### 指定网关网卡

当节点存在多张网卡，且未配置默认路由，则可以指定网络插件出外网的网卡

```
GATEWAY_INTERFACE_ENABLE: "true"
GATEWAY_INTERFACE: "ens4"
GATEWAY_NEXTHOP: "192.168.200.1"
```