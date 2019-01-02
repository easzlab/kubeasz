# 01-创建证书和环境配置

本步骤[01.prepare.yml](../../01.prepare.yml)主要完成:

- chrony role: 集群节点时间同步[可选]
- deploy role: 创建CA证书、kubeconfig、kube-proxy.kubeconfig
- prepare role: 分发CA证书、kubectl客户端安装、环境配置
- lb role: 安装负载均衡[可选]

## deploy 角色

请在另外窗口打开[roles/deploy/tasks/main.yml](../../roles/deploy/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建 CA 证书和秘钥 
``` bash
roles/deploy/
├── tasks
│   └── main.yml
└── templates
    ├── admin-csr.json.j2	# kubectl客户端使用的证书请求模板
    ├── ca-config.json.j2	# ca 配置文件模板
    ├── ca-csr.json.j2		# ca 证书签名请求模板
    ├── kubedns.yaml.j2
    └── kube-proxy-csr.json.j2	# kube-proxy使用的证书请求模板
```
kubernetes 系统各组件需要使用 TLS 证书对通信进行加密，使用 CloudFlare 的 PKI 工具集生成自签名的 CA 证书，用来签名后续创建的其它 TLS 证书。[参考阅读](https://coreos.com/os/docs/latest/generate-self-signed-certificates.html)

根据认证对象可以将证书分成三类：服务器证书`server cert`，客户端证书`client cert`，对等证书`peer cert`(表示既是`server cert`又是`client cert`)，在kubernetes 集群中需要的证书种类如下：

+ `etcd` 节点需要标识自己服务的`server cert`，也需要`client cert`与`etcd`集群其他节点交互，当然可以分别指定2个证书，也可以使用一个对等证书
+ `master` 节点需要标识 apiserver服务的`server cert`，也需要`client cert`连接`etcd`集群，这里也使用一个对等证书
+ `kubectl` `calico` `kube-proxy` 只需要`client cert`，因此证书请求中 `hosts` 字段可以为空
+ `kubelet` 证书比较特殊，不是手动生成，它由node节点`TLS BootStrap` 向`apiserver`请求，由`master`节点的`controller-manager` 自动签发，包含一个`client cert` 和一个`server cert`

整个集群要使用统一的CA 证书，只需要在 deploy 节点创建，然后分发给其他节点；为了保证安装的幂等性，如果已经存在CA 证书，就跳过创建CA 步骤

#### 创建 CA 配置文件 [ca-config.json.j2](../../roles/deploy/templates/ca-config.json.j2)
``` bash
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "kubernetes": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "87600h"
      }
    }
  }
}
```
+ `signing`：表示该证书可用于签名其它证书；生成的 ca.pem 证书中 `CA=TRUE`；
+ `server auth`：表示可以用该 CA 对 server 提供的证书进行验证；
+ `client auth`：表示可以用该 CA 对 client 提供的证书进行验证；
+ `profile kubernetes` 包含了`server auth`和`client auth`，所以可以签发三种不同类型证书；

#### 创建 CA 证书签名请求 [ca-csr.json.j2](../../roles/deploy/templates/ca-csr.json.j2)
``` bash
{
  "CN": "kubernetes",
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
  ],
  "ca": {
    "expiry": "876000h"
  }
}
```

#### 生成CA 证书和私钥
``` bash
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

### 生成 kubeconfig 配置文件

kubectl使用~/.kube/config 配置文件与kube-apiserver进行交互，且拥有管理 K8S集群的完全权限，

准备kubectl使用的admin 证书签名请求 [admin-csr.json.j2](../../roles/deploy/templates/admin-csr.json.j2)

``` bash
{
  "CN": "admin",
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
      "O": "system:masters",
      "OU": "System"
    }
  ]
}

```
+ kubectl 使用客户端证书可以不指定hosts 字段
+ 证书请求中 `O` 指定该证书的 Group 为 `system:masters`，而 `RBAC` 预定义的 `ClusterRoleBinding` 将 Group `system:masters` 与 ClusterRole `cluster-admin` 绑定，这就赋予了kubectl**所有集群权限**

``` bash
$ kubectl describe clusterrolebinding cluster-admin
Name:         cluster-admin
Labels:       kubernetes.io/bootstrapping=rbac-defaults
Annotations:  rbac.authorization.kubernetes.io/autoupdate=true
Role:
  Kind:  ClusterRole
  Name:  cluster-admin
Subjects:
  Kind   Name            Namespace
  ----   ----            ---------
  Group  system:masters  
```

#### 生成 cluster-admin 用户证书

```
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes admin-csr.json | cfssljson -bare admin
```

#### 生成 ~/.kube/config 配置文件

使用`kubectl config` 生成kubeconfig 自动保存到 ~/.kube/config，生成后 `cat ~/.kube/config`可以验证配置文件包含 kube-apiserver 地址、证书、用户名等信息。

```
kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=127.0.0.1:8443
kubectl config set-credentials admin --client-certificate=admin.pem --embed-certs=true --client-key=admin-key.pem
kubectl config set-context kubernetes --cluster=kubernetes --user=admin
kubectl config use-context kubernetes
```

### 生成 kube-proxy.kubeconfig 配置文件

创建 kube-proxy 证书请求

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
+ kube-proxy 使用客户端证书可以不指定hosts 字段
+ CN 指定该证书的 User 为 system:kube-proxy，预定义的 ClusterRoleBinding system:node-proxier 将User system:kube-proxy 与 Role system:node-proxier 绑定，授予了调用 kube-apiserver Proxy 相关 API 的权限；

``` bash
$ kubectl describe clusterrolebinding system:node-proxier
Name:         system:node-proxier
Labels:       kubernetes.io/bootstrapping=rbac-defaults
Annotations:  rbac.authorization.kubernetes.io/autoupdate=true
Role:
  Kind:  ClusterRole
  Name:  system:node-proxier
Subjects:
  Kind  Name               Namespace
  ----  ----               ---------
  User  system:kube-proxy  
```

#### 生成 system:kube-proxy 用户证书

```
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes kube-proxy-csr.json | cfssljson -bare kube-proxy
```

#### 生成 kube-proxy.kubeconfig

使用`kubectl config` 生成kubeconfig 自动保存到 kube-proxy.kubeconfig

```
kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=127.0.0.1:8443 --kubeconfig=kube-proxy.kubeconfig
kubectl config set-credentials kube-proxy --client-certificate=kube-proxy.pem --embed-certs=true --client-key=kube-proxy-key.pem --kubeconfig=kube-proxy.kubeconfig
kubectl config set-context default --cluster=kubernetes --user=kube-proxy --kubeconfig=kube-proxy.kubeconfig
kubectl config use-context default --kubeconfig=kube-proxy.kubeconfig
```

## prepare 角色

``` bash
roles/prepare/
├── files
│   ├── 95-k8s-sysctl.conf
└── tasks
    └── main.yml
```
请在另外窗口打开[roles/prepare/tasks/main.yml](../../roles/prepare/tasks/main.yml) 文件，比较简单直观

1. 首先创建一些基础文件目录
1. 修改环境变量，把{{ bin_dir }} 添加到$PATH，需要重新登陆 shell生效
1. 把证书工具 CFSSL 和 kubectl 下发到指定节点，并下发kubeconfig配置文件
1. 把CA 证书相关下发到指定节点的 {{ ca_dir }} 目录
1. 最后设置基础操作系统软件和系统参数，请阅读脚本中的注释内容

## LB 角色-负载均衡部署
``` bash
roles/lb
├── tasks
│   └── main.yml
└── templates
    ├── haproxy.cfg.j2
    ├── haproxy.service.j2
    ├── keepalived-backup.conf.j2
    └── keepalived-master.conf.j2
```

Haproxy支持四层和七层负载，稳定性好，根据官方文档，HAProxy可以跑满10Gbps-New benchmark of HAProxy at 10 Gbps using Myricom's 10GbE NICs (Myri-10G PCI-Express)；另外，openstack高可用也有用haproxy的。

keepalived观其名可知，保持存活，它是基于VRRP协议保证所谓的高可用或热备的，这里用来预防haproxy的单点故障。

keepalived与haproxy配合，实现master的高可用过程如下：

+ 1.keepalived利用vrrp协议生成一个虚拟地址(VIP)，正常情况下VIP存活在keepalive的主节点，当主节点故障时，VIP能够漂移到keepalived的备节点，保障VIP地址可用性。
+ 2.在keepalived的主备节点都配置相同haproxy负载配置，并且监听客户端请求在VIP的地址上，保障随时都有一个haproxy负载均衡在正常工作。并且keepalived启用对haproxy进程的存活检测，一旦主节点haproxy进程故障，VIP也能切换到备节点，从而让备节点的haproxy进行负载工作。
+ 3.在haproxy的配置中配置多个后端真实kube-apiserver的endpoints，并启用存活监测后端kube-apiserver，如果一个kube-apiserver故障，haproxy会将其剔除负载池。

请在另外窗口打开[roles/lb/tasks/main.yml](../../roles/lb/tasks/main.yml) 文件，对照看以下讲解内容。

#### 安装haproxy

+ 使用apt源安装

#### 配置haproxy [haproxy.cfg.j2](../../roles/lb/templates/haproxy.cfg.j2)
``` bash
global
        log /dev/log    local0
        log /dev/log    local1 notice
        chroot /var/lib/haproxy
        stats socket /run/haproxy/admin.sock mode 660 level admin
        stats timeout 30s
        user haproxy
        group haproxy
        daemon
        nbproc 1

defaults
        log     global
        timeout connect 5000
        timeout client  50000
        timeout server  50000

listen kube-master
        bind 0.0.0.0:{{ KUBE_APISERVER.split(':')[2] }}
        mode tcp
        option tcplog
        balance source
        server s1 {{ master1 }}  check inter 10000 fall 2 rise 2 weight 1
        server s2 {{ master2 }}  check inter 10000 fall 2 rise 2 weight 1
```
如果用apt安装的话，可以在/usr/share/doc/haproxy目录下找到配置指南configuration.txt.gz，全局和默认配置这里不展开，关注`listen` 代理设置模块，各项配置说明：
+ 名称 kube-master
+ bind 监听客户端请求的地址/端口，保证监听master的VIP地址和端口
+ mode 选择四层负载模式 (当然你也可以选择七层负载，请查阅指南，适当调整)
+ balance 选择负载算法 (负载算法也有很多供选择)
+ server 配置master节点真实的endpoits，必须与 [hosts文件](../../example/hosts.m-masters.example)对应设置

#### 安装keepalived

+ 使用apt源安装

#### 配置keepalived主节点 [keepalived-master.conf.j2](../../roles/lb/templates/keepalived-master.conf.j2)
``` bash
global_defs {
    router_id lb-master
}

vrrp_script check-haproxy {
    script "killall -0 haproxy"
    interval 5
    weight -30
}

vrrp_instance VI-kube-master {
    state MASTER
    priority 120
    dont_track_primary
    interface {{ LB_IF }}
    virtual_router_id {{ ROUTER_ID }}
    advert_int 3
    track_script {
        check-haproxy
    }
    virtual_ipaddress {
        {{ MASTER_IP }}
    }
}
```
+ vrrp_script 定义了监测haproxy进程的脚本，利用shell 脚本`killall -0 haproxy` 进行检测进程是否存活，如果进程不存在，根据`weight -30`设置将主节点优先级降低30，这样原先备节点将变成主节点。
+ vrrp_instance 定义了vrrp组，包括优先级、使用端口、router_id、心跳频率、检测脚本、虚拟地址VIP等
+ 特别注意 `virtual_router_id` 标识了一个 VRRP组，在同网段下必须唯一，否则出现 `Keepalived_vrrp: bogus VRRP packet received on eth0 !!!`类似报错

#### 配置keepalived备节点 [keepalived-backup.conf.j2](../../roles/lb/templates/keepalived-backup.conf.j2)
``` bash
global_defs {
    router_id lb-backup
}

vrrp_instance VI-kube-master {
    state BACKUP
    priority 110
    dont_track_primary
    interface {{ LB_IF }}
    virtual_router_id {{ ROUTER_ID }}
    advert_int 3
    virtual_ipaddress {
        {{ MASTER_IP }}
    }
}
```
+ 备节点的配置类似主节点，除了优先级和检测脚本，其他如 `virtual_router_id` `advert_int` `virtual_ipaddress`必须与主节点一致

### 启动 keepalived 和 haproxy 后验证

+ lb 节点验证

``` bash
systemctl status haproxy 	# 检查进程状态
journalctl -u haproxy		# 检查进程日志是否有报错信息
systemctl status keepalived 	# 检查进程状态
journalctl -u keepalived	# 检查进程日志是否有报错信息
netstat -antlp|grep 8443	# 检查tcp端口是否监听
```
+ 在 keepalived 主节点

``` bash
ip a				# 检查 master的 VIP地址是否存在
```
### keepalived 主备切换演练

1. 尝试关闭 keepalived主节点上的 haproxy进程，然后在keepalived 备节点上查看 master的 VIP地址是否能够漂移过来，并依次检查上一步中的验证项。
1. 尝试直接关闭 keepalived 主节点系统，检查各验证项。

[后一篇](02-install_etcd.md)
