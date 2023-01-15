## 05-安装kube_node节点

`kube_node` 是集群中运行工作负载的节点，前置条件需要先部署好`kube_master`节点，它需要部署如下组件：

``` bash
cat playbooks/05.kube-node.yml
- hosts: kube_node
  roles:
  - { role: kube-lb, when: "inventory_hostname not in groups['kube_master']" }
  - { role: kube-node, when: "inventory_hostname not in groups['kube_master']" }
```

+ kube-lb：由nginx裁剪编译的四层负载均衡，用于将请求转发到主节点的 apiserver服务
+ kubelet：kube_node上最主要的组件
+ kube-proxy： 发布应用服务与负载均衡

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
{% if ansible_distribution == "Debian" and ansible_distribution_version|int >= 10 %}
ExecStartPre=/bin/mount -o remount,rw '/sys/fs/cgroup'
{% endif %}
{% if KUBE_RESERVED_ENABLED == "yes" or SYS_RESERVED_ENABLED == "yes" %}
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpu/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuacct/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuset/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/memory/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/pids/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/systemd/podruntime.slice

ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpu/system.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuacct/system.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/cpuset/system.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/memory/system.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/pids/system.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/systemd/system.slice

{% if ansible_distribution != "Debian" %}
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/hugetlb/podruntime.slice
ExecStartPre=/bin/mkdir -p /sys/fs/cgroup/hugetlb/system.slice
{% endif %}
{% endif %}
ExecStart={{ bin_dir }}/kubelet \
  --config=/var/lib/kubelet/config.yaml \
  --container-runtime-endpoint=unix:///run/containerd/containerd.sock \
  --hostname-override={{ K8S_NODENAME }} \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --root-dir={{ KUBELET_ROOT_DIR }} \
  --v=2
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
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
  --config=/var/lib/kube-proxy/kube-proxy-config.yaml
Restart=always
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

请注意 [kube-proxy-config](../../roles/kube-node/templates/kube-proxy-config.yaml.j2) 文件的注释说明

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
