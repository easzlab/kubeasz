## 05-安装kube-node节点

`kube-node` 是集群中承载应用的节点，前置条件需要先部署好`kube-master`节点(因为需要操作`用户角色绑定`、`批准kubelet TLS 证书请求`等)，它需要部署如下组件：

+ docker：运行容器
+ calico： 配置容器网络 (或者 flannel)
+ kubelet： kube-node上最主要的组件
+ kube-proxy： 发布应用服务与负载均衡

``` bash
roles/kube-node
├── tasks
│   └── main.yml
└── templates
    ├── cni-default.conf.j2
    ├── kubelet.service.j2
    ├── kubelet-csr.json.j2
    └── kube-proxy.service.j2
```

请在另外窗口打开[roles/kube-node/tasks/main.yml](../../roles/kube-node/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建cni 基础网络插件配置文件

因为后续需要用 `DaemonSet Pod`方式运行k8s网络插件，所以kubelet.server服务必须开启cni相关参数，并且提供cni网络配置文件

### 创建 kubelet 的服务文件

+ 必须先创建工作目录 `/var/lib/kubelet`

``` bash
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=docker.service
Requires=docker.service

[Service]
WorkingDirectory=/var/lib/kubelet
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuset/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/hugetlb/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/memory/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/pids/system.slice/kubelet.service
ExecStart={{ bin_dir }}/kubelet \
  --address={{ inventory_hostname }} \
  --allow-privileged=true \
  --anonymous-auth=false \
  --authentication-token-webhook \
  --authorization-mode=Webhook \
  --pod-manifest-path=/etc/kubernetes/manifest \
  --client-ca-file={{ ca_dir }}/ca.pem \
  --cluster-dns={{ CLUSTER_DNS_SVC_IP }} \
  --cluster-domain={{ CLUSTER_DNS_DOMAIN }} \
  --cni-bin-dir={{ bin_dir }} \
  --cni-conf-dir=/etc/cni/net.d \
  --fail-swap-on=false \
  --hairpin-mode hairpin-veth \
  --hostname-override={{ inventory_hostname }} \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --max-pods={{ MAX_PODS }} \
  --network-plugin=cni \
  --pod-infra-container-image=mirrorgooglecontainers/pause-amd64:3.1 \
  --register-node=true \
  --root-dir={{ KUBELET_ROOT_DIR }} \
  --tls-cert-file={{ ca_dir }}/kubelet.pem \
  --tls-private-key-file={{ ca_dir }}/kubelet-key.pem \
  --cgroups-per-qos=true \
  --cgroup-driver=cgroupfs \
  --enforce-node-allocatable=pods,kube-reserved \
  --kube-reserved={{ KUBE_RESERVED }} \
  --kube-reserved-cgroup=/system.slice/kubelet.service \
  --eviction-hard={{ HARD_EVICTION }} \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```
+ --pod-infra-container-image 指定`基础容器`（负责创建Pod 内部共享的网络、文件系统等）镜像，**K8S每一个运行的 POD里面必然包含这个基础容器**，如果它没有运行起来那么你的POD 肯定创建不了，kubelet日志里面会看到类似 ` FailedCreatePodSandBox` 错误，可用`docker images` 查看节点是否已经下载到该镜像
+ --cluster-dns 指定 kubedns 的 Service IP(可以先分配，后续创建 kubedns 服务时指定该 IP)，--cluster-domain 指定域名后缀，这两个参数同时指定后才会生效；
+ --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir={{ bin_dir }} 为使用cni 网络，并调用calico管理网络所需的配置
+ --fail-swap-on=false K8S 1.8+需显示禁用这个，否则服务不能启动
+ --client-ca-file={{ ca_dir }}/ca.pem 和 --anonymous-auth=false 关闭kubelet的匿名访问，详见[匿名访问漏洞说明](mixes/01.fix_kubelet_annoymous_access.md)
+ --ExecStartPre=/bin/mkdir -p xxx 对于某些系统（centos7）cpuset和hugetlb 是默认没有初始化system.slice 的，需要手动创建，否则在启用--kube-reserved-cgroup 时会报错Failed to start ContainerManager Failed to enforce System Reserved Cgroup Limits
+ 关于kubelet资源预留相关配置请参考 https://kubernetes.io/docs/tasks/administer-cluster/reserve-compute-resources/

### 创建 kube-proxy kubeconfig 文件

该步骤已经在 deploy节点完成，[roles/deploy/tasks/main.yml](../../roles/deploy/tasks/main.yml)

+ 生成的kube-proxy.kubeconfig 配置文件需要移动到/etc/kubernetes/目录，后续kube-proxy服务启动参数里面需要指定

### 创建 kube-proxy服务文件

``` bash
[Unit]
Description=Kubernetes Kube-Proxy Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target

[Service]
WorkingDirectory=/var/lib/kube-proxy
ExecStart={{ bin_dir }}/kube-proxy \
  --bind-address={{ inventory_hostname }} \
  --hostname-override={{ inventory_hostname }} \
  --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig \
  --logtostderr=true \
  --v=2
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

+ --hostname-override 参数值必须与 kubelet 的值一致，否则 kube-proxy 启动后会找不到该 Node，从而不会创建任何 iptables 规则
+ 特别注意：kube-proxy 根据 --cluster-cidr 判断集群内部和外部流量，指定 --cluster-cidr 或 --masquerade-all 选项后 kube-proxy 才会对访问 Service IP 的请求做 SNAT；但是这个特性与calico 实现 network policy冲突，所以如果要用 network policy，这两个选项都不要指定。

### 批准kubelet 的 TLS 证书请求

``` bash
sleep 15 && {{ bin_dir }}/kubectl get csr|grep 'Pending' | awk 'NR>0{print $1}'| xargs {{ bin_dir }}/kubectl certificate approve
```
+ 增加15秒延时等待kubelet启动
+ `kubectl get csr |grep 'Pending'` 找出待批准的 TLS请求
+ `kubectl certificate approve` 批准请求

### 验证 node 状态

``` bash
systemctl status kubelet	# 查看状态
systemctl status kube-proxy
journalctl -u kubelet		# 查看日志
journalctl -u kube-proxy 
```
运行 `kubectl get node` 可以看到类似

``` bash
NAME           STATUS    ROLES     AGE       VERSION
192.168.1.42   Ready     <none>    2d        v1.9.0
192.168.1.43   Ready     <none>    2d        v1.9.0
192.168.1.44   Ready     <none>    2d        v1.9.0
```


[后一篇](06-install_network_plugin.md)
