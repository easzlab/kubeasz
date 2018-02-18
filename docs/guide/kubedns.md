## 部署 kubedns

kubedns 是 k8s 集群首先需要部署的，集群中的其他 pods 使用它提供域名解析服务；主要可以解析 `集群服务名` 和 `Pod hostname`；

配置文件参考 `https://github.com/kubernetes/kubernetes` 项目目录 `kubernetes/cluster/addons/dns` 

更新 `kube-dns to 1.14.8`，如果集群中已经运行kubedns插件，请使用`RollingUpdate`如下：

``` bash
kubectl set image -n kube-system deploy/kube-dns kubedns=mirrorgooglecontainers/k8s-dns-kube-dns-amd64:1.14.8
kubectl set image -n kube-system deploy/kube-dns dnsmasq=mirrorgooglecontainers/k8s-dns-dnsmasq-nanny-amd64:1.14.8
kubectl set image -n kube-system deploy/kube-dns sidecar=mirrorgooglecontainers/k8s-dns-sidecar-amd64:1.14.8
```

### 安装

**kubectl create -f /etc/ansible/manifests/kubedns/[kubedns.yaml](../../manifests/kubedns/kubedns.yaml)**

+ 注意deploy中使用的 serviceAccount `kube-dns`，该预定义的 ClusterRoleBinding system:kube-dns 将 kube-system 命名空间的 kube-dns ServiceAccount 与 system:kube-dns ClusterRole 绑定， 因此POD 具有访问 kube-apiserver DNS 相关 API 的权限；
+ 集群 pod默认继承 node的dns 解析，修改 kubelet服务启动参数 --resolv-conf=""，可以更改这个特性，详见 kubelet 启动参数

### 验证 kubedns

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
