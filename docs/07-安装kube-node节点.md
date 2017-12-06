## 07-安装kube-node节点.md

node 是集群中承载应用的节点，前置条件需要先部署好master节点(因为需要操作`用户角色绑定`、`批准kubelet TLS 证书请求`等)，它需要部署如下组件：

+ docker 运行容器
+ calico 配置容器网络
+ kubelet node上最主要的服务
+ kubeproxy 发布应用服务与负载均衡

``` bash
roles/kube-node
├── files
│   └── rbac.yaml
├── tasks
│   └── main.yml
└── templates
    ├── calico-kube-controllers.yaml.j2
    ├── kubelet.service.j2
    ├── kube-proxy-csr.json.j2
    └── kube-proxy.service.j2
```

请在另外窗口打开[roles/kube-node/tasks/main.yml](../roles/kube-node/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建角色绑定

kubelet 启动时向 kube-apiserver 发送 TLS bootstrapping 请求，需要先将 bootstrap token 文件中的 kubelet-bootstrap 用户赋予 system:node-bootstrapper 角色，然后 kubelet 才有权限创建认证请求

``` bash
# 增加15秒延时是为了等待上一步kube-master 启动完全
"sleep 15 && {{ bin_dir }}/kubectl create clusterrolebinding kubelet-bootstrap \
        --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap"
```

### 创建 bootstrapping kubeconfig 文件

``` bash
#设置集群参数
  shell: "{{ bin_dir }}/kubectl config set-cluster kubernetes \
        --certificate-authority={{ ca_dir }}/ca.pem \
        --embed-certs=true \
        --server={{ KUBE_APISERVER }} \
        --kubeconfig=bootstrap.kubeconfig"
#设置客户端认证参数
  shell: "{{ bin_dir }}/kubectl config set-credentials kubelet-bootstrap \
        --token={{ BOOTSTRAP_TOKEN }} \
        --kubeconfig=bootstrap.kubeconfig"
#设置上下文参数
  shell: "{{ bin_dir }}/kubectl config set-context default \
        --cluster=kubernetes \
        --user=kubelet-bootstrap \
        --kubeconfig=bootstrap.kubeconfig"
#选择默认上下文
  shell: "{{ bin_dir }}/kubectl config use-context default --kubeconfig=bootstrap.kubeconfig"
```
+ 注意 kubelet bootstrapping认证时是靠 token的，后续由 `master`为其生成证书和私钥
+ 以上生成的bootstrap.kubeconfig配置文件需要移动到/etc/kubernetes/目录下，后续在kubelet启动参数中指定该目录下的 bootstrap.kubeconfig

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
#--pod-infra-container-image=registry.access.redhat.com/rhel7/pod-infrastructure:latest
ExecStart={{ bin_dir }}/kubelet \
  --address={{ NODE_IP }} \
  --hostname-override={{ NODE_IP }} \
  --pod-infra-container-image=mirrorgooglecontainers/pause-amd64:3.0 \
  --experimental-bootstrap-kubeconfig=/etc/kubernetes/bootstrap.kubeconfig \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --cert-dir={{ ca_dir }} \
  --network-plugin=cni \
  --cni-conf-dir=/etc/cni/net.d \
  --cni-bin-dir={{ bin_dir }} \
  --cluster-dns={{ CLUSTER_DNS_SVC_IP }} \
  --cluster-domain={{ CLUSTER_DNS_DOMAIN }} \
  --cloud-provider='' \
  --hairpin-mode hairpin-veth \
  --allow-privileged=true \
  --fail-swap-on=false \
  --logtostderr=true \
  --v=2
#kubelet cAdvisor 默认在所有接口监听 4194 端口的请求, 以下iptables限制内网访问
ExecStartPost=/sbin/iptables -A INPUT -s 10.0.0.0/8 -p tcp --dport 4194 -j ACCEPT
ExecStartPost=/sbin/iptables -A INPUT -s 172.16.0.0/12 -p tcp --dport 4194 -j ACCEPT
ExecStartPost=/sbin/iptables -A INPUT -s 192.168.0.0/16 -p tcp --dport 4194 -j ACCEPT
ExecStartPost=/sbin/iptables -A INPUT -p tcp --dport 4194 -j DROP
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```
+ --pod-infra-container-image 指定`基础容器`的镜像，负责创建Pod 内部共享的网络、文件系统等
+ --experimental-bootstrap-kubeconfig 指向 bootstrap kubeconfig 文件，kubelet 使用该文件中的用户名和 token 向 kube-apiserver 发送 TLS Bootstrapping 请求
+ --cluster-dns 指定 kubedns 的 Service IP(可以先分配，后续创建 kubedns 服务时指定该 IP)，--cluster-domain 指定域名后缀，这两个参数同时指定后才会生效；
+ --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir={{ bin_dir }} 为使用cni 网络，并调用calico管理网络所需的配置
+ --fail-swap-on=false K8S 1.8需显示禁用这个，否则服务不能启动

### 批准kubelet 的 TLS 证书请求

``` bash
sleep 15 && {{ bin_dir }}/kubectl get csr|grep 'Pending' | awk 'NR>0{print $1}'| xargs {{ bin_dir }}/kubectl certificate approve
```
+ 增加15秒延时等待kubelet启动
+ `kubectl get csr |grep 'Pending'` 找出待批准的 TLS请求
+ `kubectl certificate approve` 批准请求

### 创建 kube-proxy 证书请求

``` bash
{
  "CN": "system:kube-proxy",
  "hosts": [],
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
+ CN 指定该证书的 User 为 system:kube-proxy，预定义的 ClusterRoleBinding system:node-proxier 将User system:kube-proxy 与 Role system:node-proxier 绑定，授予了调用 kube-apiserver Proxy 相关 API 的权限；
+ kube-proxy 使用客户端证书可以不指定hosts 字段

### 创建 kube-proxy kubeconfig 文件

``` bash
#设置集群参数
  shell: "{{ bin_dir }}/kubectl config set-cluster kubernetes \
        --certificate-authority={{ ca_dir }}/ca.pem \
        --embed-certs=true \
        --server={{ KUBE_APISERVER }} \
        --kubeconfig=kube-proxy.kubeconfig"
#设置客户端认证参数
  shell: "{{ bin_dir }}/kubectl config set-credentials kube-proxy \
        --client-certificate={{ ca_dir }}/kube-proxy.pem \
        --client-key={{ ca_dir }}/kube-proxy-key.pem \
        --embed-certs=true \
        --kubeconfig=kube-proxy.kubeconfig"
#设置上下文参数
  shell: "{{ bin_dir }}/kubectl config set-context default \
        --cluster=kubernetes \
        --user=kube-proxy \
        --kubeconfig=kube-proxy.kubeconfig"
#选择默认上下文
  shell: "{{ bin_dir }}/kubectl config use-context default --kubeconfig=kube-proxy.kubeconfig"
```
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
  --bind-address={{ NODE_IP }} \
  --hostname-override={{ NODE_IP }} \
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

### 部署calico-kube-controllers 

calico networkpolicy正常工作需要3个组件：

+ `master/node` 节点需要运行的 docker 容器 `calico/node`
+ `cni-plugin` 所需的插件二进制和配置文件
+ `calico kubernetes controllers` 负责监听Network Policy的变化，并将Policy应用到相应的网络接口

#### 准备RBAC和calico-kube-controllers.yaml 文件

- [RBAC](../roles/kube-node/files/rbac.yaml)
  - 最小化权限使用
- [Controllers](../roles/kube-node/templates/calico-kube-controllers.yaml.j2)
  - 注意只能跑一个 controller实例
  - 注意该 controller实例需要使用宿主机网络 `hostNetwork: true`

#### 创建calico-kube-controllers

``` bash
"sleep 15 && {{ bin_dir }}/kubectl create -f /root/local/kube-system/calico/rbac.yaml && \ 
        {{ bin_dir }}/kubectl create -f /root/local/kube-system/calico/calico-kube-controllers.yaml"
```
+ 增加15s等待集群node ready
