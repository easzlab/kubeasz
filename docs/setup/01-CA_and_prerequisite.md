# 01-创建证书和环境准备

本步骤主要完成: 

- (optional) role:os-harden，可选系统加固，符合linux安全基线，详见[upstream](https://github.com/dev-sec/ansible-collection-hardening/tree/master/roles/os_hardening)
- (optional) role:chrony，[可选集群节点时间同步](../guide/chrony.md)
- role:deploy，创建CA证书、集群组件访问apiserver所需的各种kubeconfig
- role:prepare，系统基础环境配置、分发CA证书、kubectl客户端安装

## deploy 角色

主要任务讲解：roles/deploy/tasks/main.yml

### 创建 CA 证书

kubernetes 系统各组件需要使用 TLS 证书对通信进行加密，使用 CloudFlare 的 PKI 工具集生成自签名的 CA 证书，用来签名后续创建的其它 TLS 证书。[参考阅读](https://coreos.com/os/docs/latest/generate-self-signed-certificates.html)

根据认证对象可以将证书分成三类：服务器证书`server cert`，客户端证书`client cert`，对等证书`peer cert`(既是`server cert`又是`client cert`)，在kubernetes 集群中需要的证书种类如下：

+ `etcd` 节点需要标识自己服务的`server cert`，也需要`client cert`与`etcd`集群其他节点交互，当然可以分别指定2个证书，为方便这里使用一个对等证书
+ `master` 节点需要标识 apiserver服务的`server cert`，也需要`client cert`连接`etcd`集群，这里也使用一个对等证书
+ `kubectl` `calico` `kube-proxy` 只需要`client cert`，因此证书请求中 `hosts` 字段可以为空
+ `kubelet` 需要标识自己服务的`server cert`，也需要`client cert`请求`apiserver`，也使用一个对等证书

整个集群要使用统一的CA 证书，只需要在ansible控制端创建，然后分发给其他节点；为了保证安装的幂等性，如果已经存在CA 证书，就跳过创建CA 步骤

#### 创建 CA 配置文件 [ca-config.json.j2](../../roles/deploy/templates/ca-config.json.j2)
``` bash
{
  "signing": {
    "default": {
      "expiry": "{{ CERT_EXPIRY }}"
    },
    "profiles": {
      "kubernetes": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "{{ CERT_EXPIRY }}"
      },
      "kcfg": {
        "usages": [
            "signing",
            "key encipherment",
            "client auth"
        ],
        "expiry": "{{ CUSTOM_EXPIRY }}"
      }
    }
  }
}
```
+ `signing`：表示该证书可用于签名其它证书；生成的 ca.pem 证书中 `CA=TRUE`；
+ `server auth`：表示可以用该 CA 对 server 提供的证书进行验证；
+ `client auth`：表示可以用该 CA 对 client 提供的证书进行验证；
+ `profile kubernetes` 包含了`server auth`和`client auth`，所以可以签发三种不同类型证书；expiry 证书有效期，默认50年
+ `profile kcfg` 在后面客户端kubeconfig证书管理中用到

#### 创建 CA 证书签名请求 [ca-csr.json.j2](../../roles/deploy/templates/ca-csr.json.j2)
``` bash
{
  "CN": "kubernetes-ca",
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
- `ca expiry` 指定ca证书的有效期，默认100年

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

#### 生成 admin 用户证书

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

### 创建kube-controller-manager 和 kube-scheduler 组件的kubeconfig 文件

过程与创建kube-proxy.kubeconfig 类似，略。

## prepare 角色

请在另外窗口打开[roles/prepare/tasks/main.yml](../../roles/prepare/tasks/main.yml) 文件，比较简单直观

1. 设置基础操作系统软件和系统参数，请阅读脚本中的注释内容
1. 创建一些基础文件目录、环境变量以及添加本地镜像仓库`easzlab.io.local`的域名解析
1. 分发kubeconfig等配置文件


[后一篇](02-install_etcd.md)
