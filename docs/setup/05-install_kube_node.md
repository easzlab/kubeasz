## 05-安装kube-node节点

`kube-node` 是集群中运行工作负载的节点，前置条件需要先部署好`kube-master`节点，它需要部署如下组件：

+ docker：运行容器
+ kubelet： kube-node上最主要的组件
+ kube-proxy： 发布应用服务与负载均衡
+ haproxy：用于请求转发到多个 apiserver，详见[HA-2x 架构](00-planning_and_overall_intro.md#ha-architecture)
+ calico： 配置容器网络 (或者其他网络组件)

``` bash
roles/kube-node/
├── defaults
│   └── main.yml		# 变量配置文件
├── tasks
│   ├── main.yml		# 主执行文件
│   ├── node_lb.yml		# haproxy 安装文件
│   └── offline.yml             # 离线安装 haproxy
└── templates
    ├── cni-default.conf.j2	# 默认cni插件配置模板
    ├── haproxy.cfg.j2		# haproxy 配置模板
    ├── haproxy.service.j2	# haproxy 服务模板
    ├── kubelet-config.yaml.j2  # kubelet 独立配置文件
    ├── kubelet-csr.json.j2	# 证书请求模板
    ├── kubelet.service.j2	# kubelet 服务模板
    └── kube-proxy.service.j2	# kube-proxy 服务模板
```

请在另外窗口打开`roles/kube-node/tasks/main.yml`文件，对照看以下讲解内容。

### 变量配置文件

详见 roles/kube-node/defaults/main.yml，举例以下3个变量配置说明
- 变量`KUBE_APISERVER`，根据不同的节点情况，它有三种取值方式
- 变量`MASTER_CHG`，变更 master 节点时会根据它来重新配置 haproxy

### 创建cni 基础网络插件配置文件

因为后续需要用 `DaemonSet Pod`方式运行k8s网络插件，所以kubelet.server服务必须开启cni相关参数，并且提供cni网络配置文件

### 创建 kubelet 的服务文件

+ 根据官方建议独立使用 kubelet 配置文件，详见roles/kube-node/templates/kubelet-config.yaml.j2
+ 必须先创建工作目录 `/var/lib/kubelet`

``` bash
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
WorkingDirectory=/var/lib/kubelet
{% if KUBE_RESERVED_ENABLED == "yes" or SYS_RESERVED_ENABLED == "yes" %}
ExecStartPre=/bin/mount -o remount,rw '/sys/fs/cgroup'
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuset/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/hugetlb/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/memory/system.slice/kubelet.service
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/pids/system.slice/kubelet.service
{% endif %}
ExecStart={{ bin_dir }}/kubelet \
  --config=/var/lib/kubelet/config.yaml \
{% if KUBE_VER|float < 1.13 %}
  --allow-privileged=true \
{% endif %}
  --cni-bin-dir={{ bin_dir }} \
  --cni-conf-dir=/etc/cni/net.d \
{% if CONTAINER_RUNTIME == "containerd" %}
  --container-runtime=remote \
  --container-runtime-endpoint=unix:///run/containerd/containerd.sock \
{% endif %}
  --hostname-override={{ inventory_hostname }} \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --network-plugin=cni \
  --pod-infra-container-image={{ SANDBOX_IMAGE }} \
  --root-dir={{ KUBELET_ROOT_DIR }} \
  --v=2
Restart=always
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
  --cluster-cidr={{ CLUSTER_CIDR }} \
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
+ 特别注意：kube-proxy 根据 --cluster-cidr 判断集群内部和外部流量，指定 --cluster-cidr 或 --masquerade-all 选项后 kube-proxy 才会对访问 Service IP 的请求做 SNAT；

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
