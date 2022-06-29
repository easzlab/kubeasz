## 06-安装calico网络组件.md

calico 是k8s社区最流行的网络插件之一，也是k8s-conformance test 默认使用的网络插件，功能丰富，支持network policy；是当前kubeasz项目的默认网络插件。

如果需要安装calico，请在`clusters/xxxx/hosts`文件中设置变量 `CLUSTER_NETWORK="calico"`，参考[这里](../config_guide.md)

``` bash
roles/calico/
├── tasks
│   └── main.yml
├── templates
│   ├── calico-csr.json.j2
│   ├── calicoctl.cfg.j2
│   ├── calico-v3.15.yaml.j2
│   ├── calico-v3.19.yaml.j2
│   └── calico-v3.8.yaml.j2
└── vars
    └── main.yml
```
请在另外窗口打开`roles/calico/tasks/main.yml`文件，对照看以下讲解内容。

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
calico 使用客户端证书，所以hosts字段可以为空；后续可以看到calico证书用在四个地方：

- calico/node 这个docker 容器运行时访问 etcd 使用证书
- cni 配置文件中，cni 插件需要访问 etcd 使用证书
- calicoctl 操作集群网络时访问 etcd 使用证书
- calico/kube-controllers 同步集群网络策略时访问 etcd 使用证书

### 创建 calico DaemonSet yaml文件和rbac 文件

请对照 roles/calico/templates/calico.yaml.j2文件注释和以下注意内容

+ 详细配置参数请参考[calico官方文档](https://projectcalico.docs.tigera.io/reference/node/configuration)
+ 配置ETCD_ENDPOINTS 、CA、证书等，所有{{ }}变量与ansible hosts文件中设置对应
+ 配置集群POD网络 CALICO_IPV4POOL_CIDR={{ CLUSTER_CIDR }}
+ 配置FELIX_DEFAULTENDPOINTTOHOSTACTION=ACCEPT 默认允许Pod到Node的网络流量，更多[felix配置选项](https://projectcalico.docs.tigera.io/reference/felix/configuration)

### 安装calico 网络

+ 安装前检查主机名不能有大写字母，只能由`小写字母` `-` `.` 组成 (name must consist of lower case alphanumeric characters, '-' or '.' (regex: [a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*))(calico-node v3.0.6以上已经解决主机大写字母问题)
+ **安装前必须确保各节点主机名不重复** ，calico node name 由节点主机名决定，如果重复，那么重复节点在etcd中只存储一份配置，BGP 邻居也不会建立。
+ 安装之前必须确保`kube_master`和`kube_node`节点已经成功部署
+ 轮询等待calico 网络插件安装完成，删除之前kube_node安装时默认cni网络配置

### [可选]配置calicoctl工具 [calicoctl.cfg.j2](roles/calico/templates/calicoctl.cfg.j2)

``` bash
apiVersion: projectcalico.org/v3
kind: CalicoAPIConfig
metadata:
spec:
  datastoreType: "etcdv3"
  etcdEndpoints: {{ ETCD_ENDPOINTS }}
  etcdKeyFile: /etc/calico/ssl/calico-key.pem
  etcdCertFile: /etc/calico/ssl/calico.pem
  etcdCACertFile: {{ ca_dir }}/ca.pem
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

**查看etcd中calico相关信息**

因为这里calico网络使用etcd存储数据，所以可以在etcd集群中查看数据

+ calico 3.x 版本默认使用 etcd v3存储，**登录集群的一个etcd 节点**，查看命令：

``` bash
# 查看所有calico相关数据
ETCDCTL_API=3 etcdctl --endpoints="http://127.0.0.1:2379" get --prefix /calico
# 查看 calico网络为各节点分配的网段
ETCDCTL_API=3 etcdctl --endpoints="http://127.0.0.1:2379" get --prefix /calico/ipam/v2/host
```


## 下一步：[设置 BGP Route Reflector](calico-bgp-rr.md)
