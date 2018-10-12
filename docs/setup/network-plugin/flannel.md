## 06-安装flannel网络组件.md

本项目提供多种网络插件可选，如果需要安装flannel，请在/etc/ansible/hosts文件中设置变量 `CLUSTER_NETWORK="flannel"`，更多设置请查看`roles/flannel/defaults/main.yml`

`Flannel`是最早应用到k8s集群的网络插件之一，简单高效，且提供多个后端`backend`模式供选择；本文介绍以`DaemonSet Pod`方式集成到k8s集群，需要在所有master节点和node节点安装。

``` text
roles/flannel/
├── tasks
│   └── main.yml
└── templates
    └── kube-flannel.yaml.j2
```

请在另外窗口打开[roles/flannel/tasks/main.yml](../roles/flannel/tasks/main.yml) 文件，对照看以下讲解内容。

### 下载基础cni 插件

请到CNI 插件最新[release](https://github.com/containernetworking/plugins/releases)页面下载[cni-v0.6.0.tgz](https://github.com/containernetworking/plugins/releases/download/v0.6.0/cni-v0.6.0.tgz)，解压后里面有很多插件，选择如下几个复制到项目 `bin`目录下

- flannel用到的插件
  - bridge
  - flannel
  - host-local
  - loopback
  - portmap

Flannel CNI 插件的配置文件可以包含多个`plugin` 或由其调用其他`plugin`；`Flannel DaemonSet Pod`运行以后会生成`/run/flannel/subnet.env `文件，例如：

``` bash
FLANNEL_NETWORK=10.1.0.0/16
FLANNEL_SUBNET=10.1.17.1/24
FLANNEL_MTU=1472
FLANNEL_IPMASQ=true
```
然后它利用这个文件信息去配置和调用`bridge`插件来生成容器网络，调用`host-local`来管理`IP`地址，例如：

``` bash
{
	"name": "mynet",
	"type": "bridge",
	"mtu": 1472,
	"ipMasq": false,
	"isGateway": true,
	"ipam": {
		"type": "host-local",
		"subnet": "10.1.17.0/24"
	}
}
```
- 更多相关介绍请阅读：
  - [flannel kubernetes 集成](https://github.com/coreos/flannel/blob/master/Documentation/kubernetes.md)
  - [flannel cni 插件](https://github.com/containernetworking/plugins/tree/master/plugins/meta/flannel)
  - [更多 cni 插件](https://github.com/containernetworking/plugins)

### 准备`Flannel DaemonSet` yaml配置文件

请阅读 `roles/flannel/templates/kube-flannel.yaml.j2` 内容，注意：

+ 本安装方式，flannel使用apiserver 存储数据
+ 配置相关RBAC 权限和 `service account`
+ 配置`ConfigMap`包含 CNI配置和 flannel配置(指定backend等)，和`hosts`文件中相关设置对应
+ `DaemonSet Pod`包含两个容器，一个容器运行flannel本身，另一个init容器部署cni 配置文件
+ 为方便国内加速使用镜像 `jmgao1983/flannel:v0.10.0-amd64` (官方镜像在docker-hub上的转存)
+ 特别注意：如果服务器是多网卡（例如vagrant环境），则需要在`roles/flannel/templates/kube-flannel.yaml.j2 `中增加指定环境变量，详见 [kubernetes ISSUE 39701](https://github.com/kubernetes/kubernetes/issues/39701)

``` bash
      ...
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: KUBERNETES_SERVICE_HOST   # 指定apiserver的主机地址
          value: {{ MASTER_IP }}
        - name: KUBERNETES_SERVICE_PORT   # 指定apiserver的服务端口
          value: {{ KUBE_APISERVER.split(':')[2] }}      
       ...
```
### 安装 flannel网络

+ 安装之前必须确保kube-master和kube-node节点已经成功部署
+ 只需要在任意装有kubectl客户端的节点运行 kubectl create安装即可
+ 等待15s后(视网络拉取相关镜像速度)，flannel 网络插件安装完成，删除之前kube-node安装时默认cni网络配置

### 验证flannel网络

执行flannel安装成功后可以验证如下：(需要等待镜像下载完成，有时候即便上一步已经配置了docker国内加速，还是可能比较慢，请确认以下容器运行起来以后，再执行后续验证步骤)

``` bash
# kubectl get pod --all-namespaces
NAMESPACE     NAME                    READY     STATUS    RESTARTS   AGE
kube-system   kube-flannel-ds-m8mzm   1/1       Running   0          3m
kube-system   kube-flannel-ds-mnj6j   1/1       Running   0          3m
kube-system   kube-flannel-ds-mxn6k   1/1       Running   0          3m
```
在集群创建几个测试pod:  `kubectl run test --image=busybox --replicas=3 sleep 30000`

``` bash
# kubectl get pod --all-namespaces -o wide|head -n 4
NAMESPACE     NAME                    READY     STATUS    RESTARTS   AGE       IP             NODE
default       busy-5956b54c8b-ld4gb   1/1       Running   0          9m        172.20.2.7     192.168.1.1
default       busy-5956b54c8b-lj9l9   1/1       Running   0          9m        172.20.1.5     192.168.1.2
default       busy-5956b54c8b-wwpkz   1/1       Running   0          9m        172.20.0.6     192.168.1.3

# 查看路由
# ip route
default via 192.168.1.254 dev ens3 onlink 
192.168.1.0/24 dev ens3  proto kernel  scope link  src 192.168.1.1 
172.17.0.0/16 dev docker0  proto kernel  scope link  src 172.17.0.1 linkdown 
172.20.0.0/24 via 192.168.1.3 dev ens3 
172.20.1.0/24 via 192.168.1.2 dev ens3 
172.20.2.0/24 dev cni0  proto kernel  scope link  src 172.20.2.1 
```
在各节点上分别 ping 这三个POD IP地址，确保能通：

``` bash
ping 172.20.2.7
ping 172.20.1.5
ping 172.20.0.6
```

