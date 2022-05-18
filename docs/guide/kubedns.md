# 集群 DNS

DNS 是 k8s 集群首要部署的组件，它为集群中的其他 pods 提供域名解析服务；主要可以解析 `集群服务名 SVC` 和 `Pod hostname`；目前建议部署 `coredns`。

NodeLocal DNSCache在集群的上运行一个dnsCache daemonset来提高clusterDNS性能和可靠性。在K8S集群上的一些测试表明：相比于纯coredns方案，nodelocaldns + coredns方案能够大幅降低DNS查询timeout的频次，提升服务稳定性。参考官方文档：https://kubernetes.io/docs/tasks/administer-cluster/nodelocaldns/

### 部署 dns

配置文件参考 `https://github.com/kubernetes/kubernetes` 项目目录 `kubernetes/cluster/addons/dns`

+ 安装

目前 kubeasz 已经自动集成安装 coredns 和 nodelocaldns 组件，配置模板位于`roles/cluster-addon/templates/`目录。

``` bash
# 默认已经集成安装，假设集群名为xxxx
ezctl setup xxxx all

# 如果需要分步安装
ezctl setup xxxx 07

# 如果需要手动安装
kubectl apply -f /etc/kubeasz/clusters/xxxx/yml/coredns.yaml
kubectl apply -f /etc/kubeasz/clusters/xxxx/yml/nodelocaldns.yaml
```

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

测试pod alpine

``` bash
kubectl run test --rm -it --image=alpine /bin/sh
If you don't see a command prompt, try pressing enter.

/ # cat /etc/resolv.conf
nameserver 10.68.0.2
search default.svc.cluster.local. svc.cluster.local. cluster.local.
options ndots:5

# 测试集群内部服务解析
/ # nslookup nginx.default.svc.cluster.local
Server:    10.68.0.2
Address 1: 10.68.0.2 kube-dns.kube-system.svc.cluster.local

Name:      nginx
Address 1: 10.68.33.167 nginx.default.svc.cluster.local

/ # nslookup kubernetes.default.svc.cluster.local
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

- Note1: 如果你使用`calico`网络组件，安装完集群后，直接安装dns组件，可能会出现如下BUG，分析是因为calico分配pod地址时候会从网段的第一个地址（网络地址）开始，详见提交的 [ISSUE #1710](https://github.com/projectcalico/calico/issues/1710)，临时解决办法为手动删除POD，重新创建后获取后面的IP地址

```
# BUG出现现象
$ kubectl get pod --all-namespaces -o wide
NAMESPACE     NAME                                       READY     STATUS             RESTARTS   AGE       IP              NODE
default       busy-5cc98488d4-s894w                      1/1       Running            0          28m       172.20.24.193   192.168.97.24
kube-system   calico-kube-controllers-6597d9c664-nq9hn   1/1       Running            0          1h        192.168.97.24   192.168.97.24
kube-system   calico-node-f8gnf                          2/2       Running            0          1h        192.168.97.24   192.168.97.24
kube-system   kube-dns-69bf9d5cc9-c68mw                  0/3       CrashLoopBackOff   27         31m       172.20.24.192   192.168.97.24

# 解决办法，删除pod，自动重建
$ kubectl delete pod -n kube-system kube-dns-69bf9d5cc9-c68mw
```

- Note2: 使用``` kubectl run test -it --rm --image=busybox /bin/sh``` 进行解析测试可能会失败, busybox内的nslookup程序有bug, 详见 https://github.com/kubernetes/dns/issues/109
