# 04-安装kube_master节点

部署master节点主要包含三个组件`apiserver` `scheduler` `controller-manager`，其中：

- apiserver提供集群管理的REST API接口，包括认证授权、数据校验以及集群状态变更等
  - 只有API Server才直接操作etcd
  - 其他模块通过API Server查询或修改数据
  - 提供其他模块之间的数据交互和通信的枢纽
- scheduler负责分配调度Pod到集群内的node节点
  - 监听kube-apiserver，查询还未分配Node的Pod
  - 根据调度策略为这些Pod分配节点
- controller-manager由一系列的控制器组成，它通过apiserver监控整个集群的状态，并确保集群处于预期的工作状态

## 高可用机制

- apiserver 无状态服务，可以通过外部负载均衡实现高可用，如项目采用的两种高可用架构：HA-1x (#584)和 HA-2x (#585)
- controller-manager 组件启动时会进行类似选举（leader）；当多副本存在时，如果原先leader挂掉，那么会选举出新的leader，从而保证高可用；
- scheduler 类似选举机制

## 安装流程

``` bash
cat playbooks/04.kube-master.yml
- hosts: kube_master
  roles:
  - kube-lb        # 四层负载均衡，监听在127.0.0.1:6443，转发到真实master节点apiserver服务
  - kube-master    #
  - kube-node      # 因为网络、监控等daemonset组件，master节点也推荐安装kubelet和kube-proxy服务
  ... 
```

### 创建 kubernetes 证书签名请求

``` bash
{
  "CN": "kubernetes",
  "hosts": [
    "127.0.0.1",
{% if groups['ex_lb']|length > 0 %}
    "{{ hostvars[groups['ex_lb'][0]]['EX_APISERVER_VIP'] }}",
{% endif %}
{% for host in groups['kube_master'] %}
    "{{ host }}",
{% endfor %}
    "{{ CLUSTER_KUBERNETES_SVC_IP }}",
{% for host in MASTER_CERT_HOSTS %}
    "{{ host }}",
{% endfor %}
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "HangZhou",
      "L": "XS",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
```
kubernetes apiserver 使用对等证书，创建时hosts字段需要配置：
- 如果配置 ex_lb，需要把 EX_APISERVER_VIP 也配置进去
- 如果需要外部访问 apiserver，可选在config.yml配置 MASTER_CERT_HOSTS
- `kubectl get svc` 将看到集群中由api-server 创建的默认服务 `kubernetes`，因此也要把 `kubernetes` 服务名和各个服务域名也添加进去

### 创建apiserver的服务配置文件

``` bash
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target

[Service]
ExecStart={{ bin_dir }}/kube-apiserver \
  --allow-privileged=true \
  --anonymous-auth=false \
  --api-audiences=api,istio-ca \
  --authorization-mode=Node,RBAC \
  --bind-address={{ inventory_hostname }} \
  --client-ca-file={{ ca_dir }}/ca.pem \
  --endpoint-reconciler-type=lease \
  --etcd-cafile={{ ca_dir }}/ca.pem \
  --etcd-certfile={{ ca_dir }}/kubernetes.pem \
  --etcd-keyfile={{ ca_dir }}/kubernetes-key.pem \
  --etcd-servers={{ ETCD_ENDPOINTS }} \
  --kubelet-certificate-authority={{ ca_dir }}/ca.pem \
  --kubelet-client-certificate={{ ca_dir }}/kubernetes.pem \
  --kubelet-client-key={{ ca_dir }}/kubernetes-key.pem \
  --secure-port={{ SECURE_PORT }} \
  --service-account-issuer=https://kubernetes.default.svc \
  --service-account-signing-key-file={{ ca_dir }}/ca-key.pem \
  --service-account-key-file={{ ca_dir }}/ca.pem \
  --service-cluster-ip-range={{ SERVICE_CIDR }} \
  --service-node-port-range={{ NODE_PORT_RANGE }} \
  --tls-cert-file={{ ca_dir }}/kubernetes.pem \
  --tls-private-key-file={{ ca_dir }}/kubernetes-key.pem \
  --requestheader-client-ca-file={{ ca_dir }}/ca.pem \
  --requestheader-allowed-names= \
  --requestheader-extra-headers-prefix=X-Remote-Extra- \
  --requestheader-group-headers=X-Remote-Group \
  --requestheader-username-headers=X-Remote-User \
  --proxy-client-cert-file={{ ca_dir }}/aggregator-proxy.pem \
  --proxy-client-key-file={{ ca_dir }}/aggregator-proxy-key.pem \
  --enable-aggregator-routing=true \
  --v=2
Restart=always
RestartSec=5
Type=notify
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
+ Kubernetes 对 API 访问需要依次经过认证、授权和准入控制(admission controll)，认证解决用户是谁的问题，授权解决用户能做什么的问题，Admission Control则是资源管理方面的作用。
+ 关于authorization-mode=Node,RBAC v1.7+支持Node授权，配合NodeRestriction准入控制来限制kubelet仅可访问node、endpoint、pod、service以及secret、configmap、PV和PVC等相关的资源；需要注意的是v1.7中Node 授权是默认开启的，v1.8中需要显式配置开启，否则 Node无法正常工作
+ 详细参数配置请参考`kube-apiserver --help`，关于认证、授权和准入控制请[阅读](https://github.com/feiskyer/kubernetes-handbook/blob/master/components/apiserver.md)
+ 增加了访问kubelet使用的证书配置，防止匿名访问kubelet的安全漏洞，详见[漏洞说明](../mixes/01.fix_kubelet_annoymous_access.md)

### 创建controller-manager 的服务文件

``` bash
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
ExecStart={{ bin_dir }}/kube-controller-manager \
  --allocate-node-cidrs=true \
  --authentication-kubeconfig=/etc/kubernetes/kube-controller-manager.kubeconfig \
  --authorization-kubeconfig=/etc/kubernetes/kube-controller-manager.kubeconfig \
  --bind-address=0.0.0.0 \
  --cluster-cidr={{ CLUSTER_CIDR }} \
  --cluster-name=kubernetes \
  --cluster-signing-cert-file={{ ca_dir }}/ca.pem \
  --cluster-signing-key-file={{ ca_dir }}/ca-key.pem \
  --kubeconfig=/etc/kubernetes/kube-controller-manager.kubeconfig \
  --leader-elect=true \
  --node-cidr-mask-size={{ NODE_CIDR_LEN }} \
  --root-ca-file={{ ca_dir }}/ca.pem \
  --service-account-private-key-file={{ ca_dir }}/ca-key.pem \
  --service-cluster-ip-range={{ SERVICE_CIDR }} \
  --use-service-account-credentials=true \
  --v=2
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
+ --cluster-cidr 指定 Cluster 中 Pod 的 CIDR 范围，该网段在各 Node 间必须路由可达(flannel/calico 等网络插件实现)
+ --service-cluster-ip-range 参数指定 Cluster 中 Service 的CIDR范围，必须和 kube-apiserver 中的参数一致
+ --cluster-signing-* 指定的证书和私钥文件用来签名为 TLS BootStrap 创建的证书和私钥
+ --root-ca-file 用来对 kube-apiserver 证书进行校验，指定该参数后，才会在Pod 容器的 ServiceAccount 中放置该 CA 证书文件
+ --leader-elect=true 使用多节点选主的方式选择主节点。只有主节点才会启动所有控制器，而其他从节点则仅执行选主算法

### 创建scheduler 的服务文件

``` bash
[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
ExecStart={{ bin_dir }}/kube-scheduler \
  --authentication-kubeconfig=/etc/kubernetes/kube-scheduler.kubeconfig \
  --authorization-kubeconfig=/etc/kubernetes/kube-scheduler.kubeconfig \
  --bind-address=0.0.0.0 \
  --kubeconfig=/etc/kubernetes/kube-scheduler.kubeconfig \
  --leader-elect=true \
  --v=2
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

+ --leader-elect=true 部署多台机器组成的 master 集群时选举产生一个处于工作状态的 kube-controller-manager 进程

### 在master 节点安装 node 服务: kubelet kube-proxy 

项目master 分支使用 DaemonSet 方式安装网络插件，如果master 节点不安装 kubelet 服务是无法安装网络插件的，如果 master 节点不安装网络插件，那么通过`apiserver` 方式无法访问 `dashboard` `kibana`等管理界面，[ISSUES #130](https://github.com/easzlab/kubeasz/issues/130)

在master 节点也同时成为 node 节点后，默认业务 POD也会调度到 master节点；可以使用 `kubectl cordon`命令禁止业务 POD调度到 master节点。

### master 集群的验证

运行 `ansible-playbook 04.kube-master.yml` 成功后，验证 master节点的主要组件：

``` bash
# 查看进程状态
systemctl status kube-apiserver
systemctl status kube-controller-manager
systemctl status kube-scheduler
# 查看进程运行日志
journalctl -u kube-apiserver
journalctl -u kube-controller-manager
journalctl -u kube-scheduler
```
执行 `kubectl get componentstatus` 可以看到

``` bash
NAME                 STATUS    MESSAGE              ERROR
scheduler            Healthy   ok                   
controller-manager   Healthy   ok                   
etcd-0               Healthy   {"health": "true"}   
etcd-2               Healthy   {"health": "true"}   
etcd-1               Healthy   {"health": "true"} 
```

[后一篇](05-install_kube_node.md)
