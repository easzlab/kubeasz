## heapster

+ 本文档基于heapster 1.5.1和k8s 1.9.x，旧版文档请看[heapster 1.4.3](heapster.1.4.3.md)

`Heapster` 监控整个集群资源的过程：首先kubelet内置的cAdvisor收集本node节点的容器资源占用情况，然后heapster从kubelet提供的api采集节点和容器的资源占用，最后heapster 持久化数据存储到`influxdb`中（也可以是其他的存储后端,Google Cloud Monitoring等）。

`Grafana` 则通过配置数据源指向上述 `influxdb`，从而界面化显示监控信息。

### 部署

访问 [heapster release](https://github.com/kubernetes/heapster)页面下载最新 release 1.5.1，参考目录`heapster-1.5.1/deploy/kube-config/influxdb`，请在参考官方yaml文件的基础上使用本项目提供的yaml文件

1. [grafana](../../manifests/heapster/grafana.yaml)
1. [heapster](../../manifests/heapster/heapster.yaml)
1. [influxdb](../../manifests/heapster/influxdb.yaml)

安装比较简单 `kubectl create -f /etc/ansible/manifests/heapster/`，主要讲一下注意事项

#### grafana.yaml配置

+ 参数`- name: GF_SERVER_ROOT_URL`的设置要根据后续访问grafana的方式确定，如果使用 NodePort方式访问，必须设置成:`value: /`；如果使用apiserver proxy方式，必须设置成`value: /api/v1/namespaces/kube-system/services/monitoring-grafana/proxy/`
+ `kubernetes.io/cluster-service: 'true'` 和 `type: NodePort` 根据上述的访问方式设置，建议使用apiserver 方式，可以增加安全控制

#### heapster.yaml配置

+ 需要配置 RBAC 把 ServiceAccount `heapster` 与集群预定义的集群角色 `system:heapster` 绑定，这样heapster pod才有相应权限去访问 apiserver

#### influxdb.yaml配置

+ influxdb 官方建议使用命令行或 HTTP API 接口来查询数据库，从 v1.1.0 版本开始默认关闭 admin UI, 从 v1.3.3 版本开始已经移除 admin UI 插件，如果你因特殊原因需要访问admin UI，请使用 v1.1.1 版本并使用configMap 配置开启它。参考[heapster 1.4.3](heapster.1.4.3.md)，具体配置yaml文件参考[influxdb v1.1.1](../../manifests/heapster/influxdb-v1.1.1/influxdb.yaml)
+ 如果你需要监控持久化存储支持，请参考[influxdb-with-pv.md](influxdb-with-pv.md)

### 验证

``` bash
$ kubectl get pods -n kube-system | grep -E 'heapster|monitoring'
heapster-3273315324-tmxbg               1/1       Running   0          11m
monitoring-grafana-2255110352-94lpn     1/1       Running   0          11m
monitoring-influxdb-884893134-3vb6n     1/1       Running   0          11m
```
扩展检查Pods日志：
``` bash
$ kubectl logs heapster-3273315324-tmxbg -n kube-system
$ kubectl logs monitoring-grafana-2255110352-94lpn -n kube-system
$ kubectl logs monitoring-influxdb-884893134-3vb6n -n kube-system
```
部署完heapster，使用上一步介绍方法查看kubernets dashboard 界面，就可以看到各 Nodes、Pods 的 CPU、内存、负载等利用率曲线图，如果 dashboard上还无法看到利用率图，使用以下命令重启 dashboard pod：
+ 首先删除 `kubectl scale deploy kubernetes-dashboard --replicas=0 -n kube-system`
+ 然后新建 `kubectl scale deploy kubernetes-dashboard --replicas=1 -n kube-system`

### 访问 grafana

#### 1.通过apiserver 访问（建议的方式）

``` bash
kubectl cluster-info | grep grafana
monitoring-grafana is running at https://x.x.x.x:6443/api/v1/namespaces/kube-system/services/monitoring-grafana/proxy
```
请参考上一步 [访问dashboard](dashboard.md)同样的方式，使用证书或者密码认证（参照hosts文件配置，默认：用户admin 密码test1234），访问`https://x.x.x.x:6443/api/v1/namespaces/kube-system/services/monitoring-grafana/proxy`即可，如图可以点击[Home]选择查看 `Cluster` `Pods`的监控图形

![grafana](../../pics/grafana.png)

#### 2.通过NodePort 访问

+ 修改 `Service` 允许 type: NodePort
+ 修改 `Deployment`中参数`- name: GF_SERVER_ROOT_URL`为 `value: /`
+ 如果之前grafana已经运行，使用 `kubectl replace --force -f /etc/ansible/manifests/heapster/grafana.yaml` 重启 grafana插件

``` bash
kubectl get svc -n kube-system|grep grafana
monitoring-grafana        NodePort    10.68.135.50    <none>        80:5855/TCP		11m
```
然后用浏览器访问 http://NodeIP:5855 



[前一篇](dashboard.md) -- [目录](index.md) -- [后一篇](ingress.md)
