## 04-安装kube-master节点.md

部署master节点主要包含三个组件`apiserver` `scheduler` `controller-manager`，其中：

- apiserver提供集群管理的REST API接口，包括认证授权、数据校验以及集群状态变更等
  - 只有API Server才直接操作etcd
  - 其他模块通过API Server查询或修改数据
  - 提供其他模块之间的数据交互和通信的枢纽
- scheduler负责分配调度Pod到集群内的node节点
  - 监听kube-apiserver，查询还未分配Node的Pod
  - 根据调度策略为这些Pod分配节点
- controller-manager由一系列的控制器组成，它通过apiserver监控整个集群的状态，并确保集群处于预期的工作状态

master节点的高可用主要就是实现apiserver组件的高可用，在之前部署lb节点时候已经配置haproxy对它进行负载均衡。

``` text
roles/kube-master/
├── tasks
│   └── main.yml
└── templates
    ├── basic-auth.csv.j2
    ├── kube-apiserver.service.j2
    ├── kube-controller-manager.service.j2
    ├── kubernetes-csr.json.j2
    ├── kube-scheduler.service.j2
    └── token.csv.j2
```

请在另外窗口打开[roles/kube-master/tasks/main.yml](../roles/kube-master/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建 kubernetes 证书签名请求

增加判断是否已经有kubernetes证书，如果是就使用原证书，跳过生成证书步骤

``` bash
{
  "CN": "kubernetes",
  "hosts": [
    "127.0.0.1",
    "{{ MASTER_IP }}",
    "{{ inventory_hostname }}",
    "{{ CLUSTER_KUBERNETES_SVC_IP }}",
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
- kubernetes 证书既是服务器证书，同时apiserver又作为客户端证书去访问etcd 集群；作为服务器证书需要设置hosts 指定使用该证书的IP 或域名列表，需要注意的是：
  - 多主高可用集群需要把master VIP地址 {{ MASTER_IP }} 也添加进去
  - `kubectl get svc` 将看到集群中由api-server 创建的默认服务 `kubernetes`，因此也要把 `kubernetes` 服务名和各个服务域名也添加进去
- 注意所有{{ }}变量与ansible hosts中设置的对应关系

### 创建 token 认证配置

因为手动为每个node节点配置TLS认证比较麻烦，后续apiserver会开启 experimental-bootstrap-token-auth 特性，利用 kubelet启动时的 token信息与此处token认证匹配认证，然后自动为 node颁发证书

``` bash
{{ BOOTSTRAP_TOKEN }},kubelet-bootstrap,10001,"system:kubelet-bootstrap"
```

### 创建基础用户名/密码认证配置

可选，为后续使用基础认证的场景做准备，如实现dashboard 用不同用户名登陆绑定不同的权限，后续更新dashboard的实践文档。

### 创建apiserver的服务配置文件

``` bash
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target

[Service]
ExecStart={{ bin_dir }}/kube-apiserver \
  --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,NodeRestriction \
  --bind-address={{ inventory_hostname }} \
  --insecure-bind-address=127.0.0.1 \
  --authorization-mode=Node,RBAC \
  --runtime-config=rbac.authorization.k8s.io/v1 \
  --kubelet-https=true \
  --kubelet-client-certificate={{ ca_dir }}/kubernetes.pem \
  --kubelet-client-key={{ ca_dir }}/kubernetes-key.pem \
  --anonymous-auth=false \
  --basic-auth-file={{ ca_dir }}/basic-auth.csv \
  --enable-bootstrap-token-auth \
  --token-auth-file={{ ca_dir }}/token.csv \
  --service-cluster-ip-range={{ SERVICE_CIDR }} \
  --service-node-port-range={{ NODE_PORT_RANGE }} \
  --tls-cert-file={{ ca_dir }}/kubernetes.pem \
  --tls-private-key-file={{ ca_dir }}/kubernetes-key.pem \
  --client-ca-file={{ ca_dir }}/ca.pem \
  --service-account-key-file={{ ca_dir }}/ca-key.pem \
  --etcd-cafile={{ ca_dir }}/ca.pem \
  --etcd-certfile={{ ca_dir }}/kubernetes.pem \
  --etcd-keyfile={{ ca_dir }}/kubernetes-key.pem \
  --etcd-servers={{ ETCD_ENDPOINTS }} \
  --enable-swagger-ui=true \
  --allow-privileged=true \
  --audit-log-maxage=30 \
  --audit-log-maxbackup=3 \
  --audit-log-maxsize=100 \
  --audit-log-path=/var/lib/audit.log \
  --event-ttl=1h \
  --v=2
Restart=on-failure
RestartSec=5
Type=notify
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
+ Kubernetes 对 API 访问需要依次经过认证、授权和准入控制(admission controll)，认证解决用户是谁的问题，授权解决用户能做什么的问题，Admission Control则是资源管理方面的作用。
+ 支持同时提供https（默认监听在6443端口）和http API（默认监听在127.0.0.1的8080端口），其中http API是非安全接口，不做任何认证授权机制，kube-scheduler、kube-controller-manager 一般和 kube-apiserver 部署在同一台机器上，它们使用非安全端口和 kube-apiserver通信; 其他集群外部就使用HTTPS访问 apiserver
+ 关于authorization-mode=Node,RBAC v1.7+支持Node授权，配合NodeRestriction准入控制来限制kubelet仅可访问node、endpoint、pod、service以及secret、configmap、PV和PVC等相关的资源；需要注意的是v1.7中Node 授权是默认开启的，v1.8中需要显式配置开启，否则 Node无法正常工作
+ 缺省情况下 kubernetes 对象保存在 etcd /registry 路径下，可以通过 --etcd-prefix 参数进行调整
+ 详细参数配置请参考`kube-apiserver --help`，关于认证、授权和准入控制请[阅读](https://github.com/feiskyer/kubernetes-handbook/blob/master/components/apiserver.md)
+ 增加了访问kubelet使用的证书配置，防止匿名访问kubelet的安全漏洞，详见[漏洞说明](mixes/01.fix_kubelet_annoymous_access.md)

### 创建controller-manager 的服务文件

``` bash
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
ExecStart={{ bin_dir }}/kube-controller-manager \
  --address=127.0.0.1 \
  --master=http://127.0.0.1:8080 \
  --allocate-node-cidrs=true \
  --service-cluster-ip-range={{ SERVICE_CIDR }} \
  --cluster-cidr={{ CLUSTER_CIDR }} \
  --cluster-name=kubernetes \
  --cluster-signing-cert-file={{ ca_dir }}/ca.pem \
  --cluster-signing-key-file={{ ca_dir }}/ca-key.pem \
  --service-account-private-key-file={{ ca_dir }}/ca-key.pem \
  --root-ca-file={{ ca_dir }}/ca.pem \
  --leader-elect=true \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```
+ --address 值必须为 127.0.0.1，因为当前 kube-apiserver 期望 scheduler 和 controller-manager 在同一台机器
+ --master=http://127.0.0.1:8080 使用非安全 8080 端口与 kube-apiserver 通信
+ --cluster-cidr 指定 Cluster 中 Pod 的 CIDR 范围，该网段在各 Node 间必须路由可达(calico 实现)
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
  --address=127.0.0.1 \
  --master=http://127.0.0.1:8080 \
  --leader-elect=true \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

+ --address 同样值必须为 127.0.0.1
+ --master=http://127.0.0.1:8080 使用非安全 8080 端口与 kube-apiserver 通信
+ --leader-elect=true 部署多台机器组成的 master 集群时选举产生一个处于工作状态的 kube-controller-manager 进程

### 在master 节点安装 node 服务: kubelet kube-proxy 

项目master 分支使用 DaemonSet 方式安装网络插件，如果master 节点不安装 kubelet 服务是无法安装网络插件的，如果 master 节点不安装网络插件，那么通过`apiserver` 方式无法访问 `dashboard` `kibana`等管理界面，[ISSUES #130](https://github.com/gjmzj/kubeasz/issues/130)

项目v1.8 分支使用二进制方式安装网络插件，所以没有这个问题

``` bash
# vi 04.kube-master.yml
- hosts: kube-master
  roles:
  - kube-master
  - kube-node
  # 禁止业务 pod调度到 master节点
  tasks:
  - name: 禁止业务 pod调度到 master节点
    shell: "{{ bin_dir }}/kubectl cordon {{ inventory_hostname }} "
    when: DEPLOY_MODE != "allinone"
    ignore_errors: true
```
在master 节点也同时成为 node 节点后，默认业务 POD也会调度到 master节点，多主模式下这显然增加了 master节点的负载，因此可以使用 `kubectl cordon`命令禁止业务 POD调度到 master节点


### master 集群的验证

运行 `ansible-playbook 06.kube-master.yml` 成功后，验证 master节点的主要组件：

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

[前一篇](03-安装docker服务.md) -- [后一篇](05-安装kube-node节点.md)
