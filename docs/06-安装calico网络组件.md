## 07-安装calico网络组件.md

推荐阅读[feiskyer-kubernetes指南](https://github.com/feiskyer/kubernetes-handbook) 网络相关内容

首先回顾下K8S网络设计原则，在配置集群网络插件或者实践K8S 应用/服务部署请时刻想到这些原则：

- 1.每个Pod都拥有一个独立IP地址，Pod内所有容器共享一个网络命名空间
- 2.集群内所有Pod都在一个直接连通的扁平网络中，可通过IP直接访问
  - 所有容器之间无需NAT就可以直接互相访问
  - 所有Node和所有容器之间无需NAT就可以直接互相访问
  - 容器自己看到的IP跟其他容器看到的一样
- 3.Service cluster IP尽可在集群内部访问，外部请求需要通过NodePort、LoadBalance或者Ingress来访问

`Container Network Interface (CNI)`是目前CNCF主推的网络模型，它由两部分组成：

- CNI Plugin负责给容器配置网络，它包括两个基本的接口
  - 配置网络: AddNetwork(net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
  - 清理网络: DelNetwork(net *NetworkConfig, rt *RuntimeConf) error
- IPAM Plugin负责给容器分配IP地址

Kubernetes Pod的网络是这样创建的：
- 0.每个Pod除了创建时指定的容器外，都有一个kubelet启动时指定的`基础容器`，比如：`mirrorgooglecontainers/pause-amd64` `registry.access.redhat.com/rhel7/pod-infrastructure`
- 1.首先 kubelet创建`基础容器`生成network namespace
- 2.然后 kubelet调用网络CNI driver，由它根据配置调用具体的CNI 插件
- 3.然后 CNI 插件给`基础容器`配置网络
- 4.最后 Pod 中其他的容器共享使用`基础容器`的网络

本文档基于CNI driver 调用calico 插件来配置kubernetes的网络，常用CNI插件有 `flannel` `calico` `weave`等等，这些插件各有优势，也在互相借鉴学习优点，比如：在所有node节点都在一个二层网络时候，flannel提供hostgw实现，避免vxlan实现的udp封装开销，估计是目前最高效的；calico也针对L3 Fabric，推出了IPinIP的选项，利用了GRE隧道封装；因此这些插件都能适合很多实际应用场景，这里选择calico，主要考虑它支持 `kubernetes network policy`。

推荐阅读[calico kubernetes guide](https://docs.projectcalico.org/v3.0/getting-started/kubernetes/)

calico-node需要在所有master节点和node节点安装 

``` bash
roles/calico/
├── tasks
│   └── main.yml
└── templates
    ├── calico-csr.json.j2
    ├── calicoctl.cfg.j2
    ├── calico-rbac.yaml.j2
    └── calico.yaml.j2
```
请在另外窗口打开[roles/calico/tasks/main.yml](../roles/calico/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建calico 证书申请

``` bash
{
  "CN": "calico",
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
- calico 使用客户端证书，所以hosts字段可以为空；后续可以看到calico证书用在四个地方：
  - calico/node 这个docker 容器运行时访问 etcd 使用证书
  - cni 配置文件中，cni 插件需要访问 etcd 使用证书
  - calicoctl 操作集群网络时访问 etcd 使用证书
  - calico/kube-controllers 同步集群网络策略时访问 etcd 使用证书

### 创建 calico DaemonSet yaml文件和rbac 文件

请对照 roles/calico/templates/calico.yaml.j2文件注释和以下注意内容

+ 详细配置参数请参考[calico官方文档](https://docs.projectcalico.org/v2.6/reference/node/configuration)
+ calico-node是以docker容器运行在host上的，因此需要把之前的证书目录 /etc/calico/ssl挂载到容器中
+ 配置ETCD_ENDPOINTS 、CA、证书等，所有{{ }}变量与ansible hosts文件中设置对应
+ 配置集群POD网络 CALICO_IPV4POOL_CIDR={{ CLUSTER_CIDR }}
+ **重要**本K8S集群运行在同网段kvm虚机上，虚机间没有网络ACL限制，因此可以设置`CALICO_IPV4POOL_IPIP=off`，如果你的主机位于不同网段，或者运行在公有云上需要打开这个选项 `CALICO_IPV4POOL_IPIP=always`
+ 配置FELIX_DEFAULTENDPOINTTOHOSTACTION=ACCEPT 默认允许Pod到Node的网络流量，更多[felix配置选项](https://docs.projectcalico.org/v2.6/reference/felix/configuration)

### 安装calico 网络

+ 安装前检查主机名不能有大写字母，只能由`小写字母` `-` `.` 组成 (name must consist of lower case alphanumeric characters, '-' or '.' (regex: [a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*))
+ 安装之前必须确保`kube-master`和`kube-node`节点已经成功部署
+ 只需要在任意装有kubectl客户端的节点运行 `kubectl create `安装即可
+ 等待15s后(视网络拉取calico相关镜像速度)，calico 网络插件安装完成，删除之前kube-node安装时默认cni网络配置

### [可选]配置calicoctl工具 [calicoctl.cfg.j2](roles/calico/templates/calicoctl.cfg.j2)

``` bash
apiVersion: v1
kind: calicoApiConfig
metadata:
spec:
  datastoreType: "etcdv2"
  etcdEndpoints: {{ ETCD_ENDPOINTS }}
  etcdKeyFile: /etc/calico/ssl/calico-key.pem
  etcdCertFile: /etc/calico/ssl/calico.pem
  etcdCACertFile: /etc/calico/ssl/ca.pem
```

### 验证calico网络

执行calico安装成功后可以验证如下：(需要等待镜像下载完成，有时候即便上一步已经配置了docker国内加速，还是可能比较慢，请确认以下容器运行起来以后，再执行后续验证步骤)

``` bash
kubectl get pod --all-namespaces
NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-5c6b98d9df-xj2n4   1/1       Running   0          1m
kube-system   calico-node-4hr52                          2/2       Running   0          1m
kube-system   calico-node-8ctc2                          2/2       Running   0          1m
kube-system   calico-node-9t8md                          2/2       Running   0          1m
```

**查看网卡和路由信息**

先在集群创建几个测试pod:  `kubectl run test --image=busybox --replicas=3 sleep 30000`

``` bash
# 查看网卡信息
ip a
```

+ 可以看到包含类似cali1cxxx的网卡，是calico为测试pod生成的
+ tunl0网卡现在不用管，是默认生成的，当开启IPIP 特性时使用的隧道

``` bash
# 查看路由
route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         192.168.1.1     0.0.0.0         UG    0      0        0 ens3
192.168.1.0     0.0.0.0         255.255.255.0   U     0      0        0 ens3
172.17.0.0      0.0.0.0         255.255.0.0     U     0      0        0 docker0
172.20.3.64     192.168.1.34    255.255.255.192 UG    0      0        0 ens3
172.20.33.128   0.0.0.0         255.255.255.192 U     0      0        0 *
172.20.33.129   0.0.0.0         255.255.255.255 UH    0      0        0 caliccc295a6d4f
172.20.104.0    192.168.1.35    255.255.255.192 UG    0      0        0 ens3
172.20.166.128  192.168.1.63    255.255.255.192 UG    0      0        0 ens3
```

**查看所有calico节点状态**

``` bash
calicoctl node status
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+--------------+-------------------+-------+----------+-------------+
| 192.168.1.34 | node-to-node mesh | up    | 12:34:00 | Established |
| 192.168.1.35 | node-to-node mesh | up    | 12:34:00 | Established |
| 192.168.1.63 | node-to-node mesh | up    | 12:34:01 | Established |
+--------------+-------------------+-------+----------+-------------+
```

**BGP 协议是通过TCP 连接来建立邻居的，因此可以用netstat 命令验证 BGP Peer**

``` bash
netstat -antlp|grep ESTABLISHED|grep 179
tcp        0      0 192.168.1.66:179        192.168.1.35:41316      ESTABLISHED 28479/bird      
tcp        0      0 192.168.1.66:179        192.168.1.34:40243      ESTABLISHED 28479/bird      
tcp        0      0 192.168.1.66:179        192.168.1.63:48979      ESTABLISHED 28479/bird
```

**查看集群ipPool情况**

``` bash
calicoctl get ipPool -o yaml
- apiVersion: v1
  kind: ipPool
  metadata:
    cidr: 172.20.0.0/16
  spec:
    nat-outgoing: true
```

**查看etcd中calico相关信息**

因为这里calico网络使用etcd存储数据，所以可以在etcd集群中查看数据

+ calico 3.x 版本默认使用 etcd v3存储，**登陆集群的一个etcd 节点**，查看命令：

``` bash
# 查看所有calico相关数据
ETCDCTL_API=3 etcdctl --endpoints="http://127.0.0.1:2379" get --prefix /calico
# 查看 calico网络为各节点分配的网段
ETCDCTL_API=3 etcdctl --endpoints="http://127.0.0.1:2379" get --prefix /calico/ipam/v2/host
```

+ calico 2.x 版本默认使用 etcd v2存储，**登陆集群的一个etcd 节点**，查看命令：

``` bash
# 查看所有calico相关数据
etcdctl --endpoints=http://127.0.0.1:2379 --ca-file=/etc/kubernetes/ssl/ca.pem ls /calico
```

[前一篇](06-安装kube-node节点.md) -- [后一篇]()
