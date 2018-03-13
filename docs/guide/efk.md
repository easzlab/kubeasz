## EFK

`EFK` 插件是`k8s`项目的一个日志解决方案，它包括三个组件：[Elasticsearch](),[Fluentd](),[Kibana]()；Elasticsearch 是日志存储和日志搜索引擎，Fluentd 负载把`k8s`集群的日志发送给 Elasticsearch, Kibana 是可视化界面查看和检索存储在 Elasticsearch 的数据。

### 部署

下载官方最新[release](https://github.com/kubernetes/kubernetes/release)，进入目录: `kubernetes/cluster/addons/fluentd-elasticsearch`，参考官方配置的基础上使用本项目`manifests/efk/`部署，以下为几点主要的修改：

+ 官方提供的`kibana-deployment.yaml`中的参数`SERVER_BASEPATH`在k8s v1.8 版本以后部署需要按照本项目调整
+ 修改官方docker镜像，方便国内下载加速

请使用`kubectl create -f /etc/ansible/manifests/efk/`进行安装

注意：Fluentd 是以 DaemonSet 形式运行且只会调度到有`beta.kubernetes.io/fluentd-ds-ready=true`标签的节点，所以对需要收集日志的节点逐个打上标签：

``` bash
$ kubectl label nodes 192.168.1.2 beta.kubernetes.io/fluentd-ds-ready=true
node "192.168.1.2" labeled
```

### 验证

``` bash
kubectl get pods -n kube-system|grep -E 'elasticsearch|fluentd|kibana'
elasticsearch-logging-0                    1/1       Running   0          19h
elasticsearch-logging-1                    1/1       Running   0          19h
fluentd-es-v2.0.2-6c95c                    1/1       Running   0          17h
fluentd-es-v2.0.2-f2xh8                    1/1       Running   0          8h
fluentd-es-v2.0.2-pv5q5                    1/1       Running   0          8h
kibana-logging-d5cffd7c6-9lz2p             1/1       Running   0          1m
```
kibana Pod 第一次启动时会用较长时间(10-20分钟)来优化和 Cache 状态页面，可以查看 Pod 的日志观察进度，等待 `Ready` 状态

``` bash
$ kubectl logs -n kube-system kibana-logging-d5cffd7c6-9lz2p -f
...
{"type":"log","@timestamp":"2018-03-13T07:33:00Z","tags":["listening","info"],"pid":1,"message":"Server running at http://0:5601"}
{"type":"log","@timestamp":"2018-03-13T07:33:00Z","tags":["status","ui settings","info"],"pid":1,"state":"green","message":"Status changed from uninitialized to green - Ready","prevState":"uninitialized","prevMsg":"uninitialized"}
```

### 访问 Kibana

这里介绍 `kube-apiserver`方式访问，获取访问 URL

``` bash
$ kubectl cluster-info | grep Kibana
Kibana is running at https://192.168.1.10:8443/api/v1/namespaces/kube-system/services/kibana-logging/proxy
```
浏览器访问 URL：`https://192.168.1.10:8443/api/v1/namespaces/kube-system/services/kibana-logging/proxy`，然后使用`basic auth`或者`证书` 的方式认证后即可，关于认证可以参考[dashboard文档](dashboard.md)

首次登陆需要在`Management` - `Index Patterns` 创建 `index pattern`，可以使用默认的 logstash-* pattern，点击 Create; 创建Index后，稍等几分钟就可以在 Discover 菜单看到 ElasticSearch logging 中汇聚的日志；


[前一篇](ingress.md) -- [目录](index.md) -- [后一篇](harbor.md)
