## 部署集群 DNS

DNS 是 k8s 集群首先需要部署的，集群中的其他 pods 使用它提供域名解析服务；主要可以解析 `集群服务名 SVC` 和 `Pod hostname`；目前 k8s v1.9+ 版本可以有两个选择：`kube-dns` 和 `coredns`，可以选择其中一个部署安装。

### 部署 dns

配置文件参考 `https://github.com/kubernetes/kubernetes` 项目目录 `kubernetes/cluster/addons/dns`

+ 安装 

``` bash
# 安装 kube-dns
$ kubectl create -f /etc/ansible/manifests/kubedns

# 或者选择安装 coredns
$ kubectl create -f /etc/ansible/manifests/coredns
```

+ 集群 pod默认继承 node的dns 解析，修改 kubelet服务启动参数 --resolv-conf=""，可以更改这个特性，详见 kubelet 启动参数

### 验证 dns服务

新建一个测试nginx服务

`kubectl run nginx --image=nginx --expose --port=80`

确认nginx服务

``` bash
kubectl get pod|grep nginx
nginx-7cbc4b4d9c-fl46v   1/1       Running   0          1m
kubectl get svc|grep nginx
nginx        ClusterIP   10.68.33.167   <none>        80/TCP    1m
```

测试pod busybox

``` bash
kubectl run busybox --rm -it --image=busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # cat /etc/resolv.conf
nameserver 10.68.0.2
search default.svc.cluster.local. svc.cluster.local. cluster.local.
options ndots:5
# 测试集群内部服务解析
/ # nslookup nginx
Server:    10.68.0.2
Address 1: 10.68.0.2 kube-dns.kube-system.svc.cluster.local

Name:      nginx
Address 1: 10.68.33.167 nginx.default.svc.cluster.local
/ # nslookup kubernetes
Server:    10.68.0.2
Address 1: 10.68.0.2 kube-dns.kube-system.svc.cluster.local

Name:      kubernetes
Address 1: 10.68.0.1 kubernetes.default.svc.cluster.local
# 测试外部域名的解析，默认集成node的dns解析
/ # nslookup www.baidu.com
Server:    10.68.0.2
Address 1: 10.68.0.2 kube-dns.kube-system.svc.cluster.local

Name:      www.baidu.com
Address 1: 180.97.33.108
Address 2: 180.97.33.107
/ #
```

[前一篇](index.md) -- [目录](index.md) -- [后一篇](dashboard.md)
