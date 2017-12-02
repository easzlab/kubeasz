## 05-安装calico网络组件.md

推荐阅读[feiskyer-kubernetes指南](https://github.com/feiskyer/kubernetes-handbook) 网络相关内容

首先回顾下K8S网络设计原则，在配置集群网络插件或者实践K8S 应用/服务部署请时刻想到这些原则：

- 1.每个Pod都拥有一个独立IP地址，Pod内所有容器共享一个网络命名空间
- 2.集群内所有Pod都在一个直接连通的扁平网络中，可通过IP直接访问
  - 所有容器之间无需NAT就可以直接互相访问
  - 所有Node和所有容器之间无需NAT就可以直接互相访问
  - 容器自己看到的IP跟其他容器看到的一样
- 3.Service cluster IP尽可在集群内部访问，外部请求需要通过NodePort、LoadBalance或者Ingress来访问

**Container Network Interface (CNI)**是目前CNCF主推的网络模型，它由两部分组成：

- CNI Plugin负责给容器配置网络，它包括两个基本的接口
  - 配置网络: AddNetwork(net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
  - 清理网络: DelNetwork(net *NetworkConfig, rt *RuntimeConf) error
- IPAM Plugin负责给容器分配IP地址

Kubernetes Pod的网络是这样创建的：
- 0.每个Pod除了创建时指定的容器外，都有一个kubelet启动时指定的`基础容器`，比如：[mirrorgooglecontainers/pause-amd64]() [registry.access.redhat.com/rhel7/pod-infrastructure]()
- 1.首先 kubelet创建`基础容器`生成network namespace
- 2.然后 kubelet调用网络CNI driver，由它根据配置调用具体的CNI 插件
- 3.然后 CNI 插件给`基础容器`配置网络
- 4.最后 Pod 中其他的容器共享使用`基础容器`的网络

本文档基于CNI driver 调用calico 插件来配置kubernetes的网络，常用CNI插件有 `flannel` `calico` `weave`等等，这些插件各有优势，也在互相借鉴学习优点，比如：在所有node节点都在一个二层网络时候，flannel提供hostgw实现，避免vxlan实现的udp封装开销，估计是目前最高效的；calico也针对L3 Fabric，推出了IPinIP的选项，利用了GRE隧道封装；因此这些插件都能适合很多实际应用场景，这里选择calico，主要考虑它支持 `kubernetes network policy`。

推荐阅读[calico kubernetes Integration Guide](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/installation/integration)

calico-node需要在所有master节点和node节点安装 

``` bash
roles/calico/
├── files
│   ├── ca.pem
│   ├── etcd-key.pem
│   └── etcd.pem
├── tasks
│   └── main.yml
└── templates
    ├── calicoctl.cfg.j2
    ├── calico-node.service.j2
    └── cni-calico.conf.j2
```
请在另外窗口打开[roles/calico/tasks/main.yml](../roles/calico/tasks/main.yml) 文件，对照看以下讲解内容。

### 准备与etcd集群交互的证书

这里为了方便直接复制使用了etcd1节点的证书

### 创建 calico-node 的服务文件 [calico-node.service.j2](../roles/calico/templates/calico-node.service.j2)

``` bash
[Unit]
Description=calico node
After=docker.service
Requires=docker.service

[Service]
User=root
PermissionsStartOnly=true
ExecStart={{ bin_dir }}/docker run --net=host --privileged --name=calico-node \
  -e ETCD_ENDPOINTS={{ ETCD_ENDPOINTS }} \
  -e ETCD_CA_CERT_FILE=/etc/calico/ssl/ca.pem \
  -e ETCD_CERT_FILE=/etc/calico/ssl//etcd.pem \
  -e ETCD_KEY_FILE=/etc/calico/ssl/etcd-key.pem \
  -e CALICO_LIBNETWORK_ENABLED=true \
  -e CALICO_NETWORKING_BACKEND=bird \
  -e CALICO_DISABLE_FILE_LOGGING=true \
  -e CALICO_IPV4POOL_CIDR={{ CLUSTER_CIDR }} \
  -e CALICO_IPV4POOL_IPIP=off \
  -e FELIX_DEFAULTENDPOINTTOHOSTACTION=ACCEPT \
  -e FELIX_IPV6SUPPORT=false \
  -e FELIX_LOGSEVERITYSCREEN=info \
  -e FELIX_IPINIPMTU=1440 \
  -e FELIX_HEALTHENABLED=true \
  -e IP= \
  -v /etc/calico/ssl:/etc/calico/ssl \
  -v /var/run/calico:/var/run/calico \
  -v /lib/modules:/lib/modules \
  -v /run/docker/plugins:/run/docker/plugins \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/log/calico:/var/log/calico \
  calico/node:v2.6.2
ExecStop={{ bin_dir }}/docker rm -f calico-node
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
+ 详细配置参数请参考[calico官方文档](https://docs.projectcalico.org/v2.6/reference/node/configuration)
+ 配置ETCD_ENDPOINTS 、CA、证书等，所有{{ }}变量与ansible hosts文件中设置对应
+ 配置集群POD网络 CALICO_IPV4POOL_CIDR={{ CLUSTER_CIDR }}
+ 本K8S集群运行在自有kvm虚机上，虚机间没有网络ACL限制，因此可以设置CALICO_IPV4POOL_IPIP=off，如果运行在公有云虚机上可能需要打开这个选项
+ 配置FELIX_DEFAULTENDPOINTTOHOSTACTION=ACCEPT 默认允许Pod到Node的网络流量，更多[felix配置选项](https://docs.projectcalico.org/v2.6/reference/felix/configuration)
+ calico-node是以docker容器运行在host上的，因此需要把之前的证书目录 /etc/calico/ssl挂载到容器中

### 启动calico-node

### 准备cni-calico配置文件 [cni-calico.conf.j2](../roles/calico/templates/cni-calico.conf.j2)

``` bash
{
    "name": "calico-k8s-network",
    "cniVersion": "0.1.0",
    "type": "calico",
    "etcd_endpoints": "{{ ETCD_ENDPOINTS }}",
    "etcd_key_file": "/etc/calico/ssl/etcd-key.pem",
    "etcd_cert_file": "/etc/calico/ssl//etcd.pem",
    "etcd_ca_cert_file": "/etc/calico/ssl/ca.pem",
    "log_level": "info",
    "mtu": 1500,
    "ipam": {
        "type": "calico-ipam"
    },
    "policy": {
        "type": "k8s"
    },
    "kubernetes": {
        "kubeconfig": "/root/.kube/config"
    }
}

```
+ 主要配置etcd相关、ipam、policy等，配置选项[参考](https://docs.projectcalico.org/v2.6/reference/cni-plugin/configuration)

### [可选]配置calicoctl工具 [calicoctl.cfg.j2](roles/calico/templates/calicoctl.cfg.j2)

``` bash
apiVersion: v1
kind: calicoApiConfig
metadata:
spec:
  datastoreType: "etcdv2"
  etcdEndpoints: {{ ETCD_ENDPOINTS }}
  etcdKeyFile: /etc/calico/ssl/etcd-key.pem
  etcdCertFile: /etc/calico/ssl/etcd.pem
  etcdCACertFile: /etc/calico/ssl/ca.pem
```

### 验证calico网络

执行calico安装 `ansible-playbook 05.calico.yml` 成功后

``` bash
 # 查看所有calico节点状态
calicoctl node status
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+------------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |   SINCE    |    INFO     |
+--------------+-------------------+-------+------------+-------------+
| 192.168.1.2  | node-to-node mesh | up    | 2017-11-30 | Established |
| 192.168.1.3  | node-to-node mesh | up    | 2017-11-30 | Established |
| 192.168.1.4  | node-to-node mesh | up    | 2017-11-30 | Established |
+--------------+-------------------+-------+------------+-------------+

IPv6 BGP status
No IPv6 peers found.

# 查看集群ipPool情况
calicoctl get ipPool -o yaml
- apiVersion: v1
  kind: ipPool
  metadata:
    cidr: 172.20.0.0/16
  spec:
    nat-outgoing: true
```

