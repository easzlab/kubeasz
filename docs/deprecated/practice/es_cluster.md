# Elasticsearch 部署实践

`Elasticsearch`是目前全文搜索引擎的首选，它可以快速地储存、搜索和分析海量数据；也可以看成是真正分布式的高效数据库集群；`Elastic`的底层是开源库`Lucene`；封装并提供了`REST API`的操作接口。

## 单节点 docker 测试安装 
 
``` bash
cat > es-start.sh << EOF
#!/bin/bash

sysctl -w vm.max_map_count=262144

docker run --detach \
   --name es01 \
   -p 9200:9200 -p 9300:9300 \
   -e "discovery.type=single-node" \
   -e "bootstrap.memory_lock=true" --ulimit memlock=-1:-1 \
   --ulimit nofile=65536:65536 \
   --volume /srv/elasticsearch/data:/usr/share/elasticsearch/data \
   --volume /srv/elasticsearch/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
   jmgao1983/elasticsearch:6.4.0
EOF
```

执行`sh es-start.sh`后，就在本地运行了。

- 验证 docker 镜像运行情况  

``` bash
root@docker-ts:~# docker ps -a
CONTAINER ID        IMAGE                           COMMAND                  CREATED             STATUS              PORTS                                            NAMES
171f3fecb596        jmgao1983/elasticsearch:6.4.0   "/usr/local/bin/do..."   2 hours ago         Up 2 hours          0.0.0.0:9200->9200/tcp, 0.0.0.0:9300->9300/tcp   es01
```

- 验证 es 健康检查  

``` bash
root@docker-ts:~# curl http://127.0.0.1:9200/_cat/health
epoch      timestamp cluster       status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1535523956 06:25:56  docker-es     green           1         1      0   0    0    0        0             0                  -                100.0%
```

## 在 k8s 上部署 Elasticsearch 集群

在生产环境下，Elasticsearch 集群由不同的角色节点组成：

- master 节点：参与主节点选举，不存储数据；建议3个以上，维护整个集群的稳定可靠状态
- data 节点：不参与选主，负责存储数据；主要消耗磁盘，内存
- client 节点：不参与选主，不存储数据；负责处理用户请求，实现请求转发，负载均衡等功能

这里使用`helm chart`来部署 (https://github.com/helm/charts/tree/master/incubator/elasticsearch)

- 1.安装 helm: 以本项目[安全安装helm](../guide/helm.md)为例
- 2.准备 PV: 以本项目[K8S 集群存储](../setup/08-cluster-storage.md)创建`nfs`动态 PV 为例
- 3.安装 elasticsearch chart  

``` bash
$ cd /etc/kubeasz/manifests/es-cluster
# 如果你的helm安装没有启用tls证书，请忽略以下--tls参数
$ helm install --tls --name es-cluster --namespace elastic -f es-values.yaml elasticsearch
```

- 4.验证 es 集群  

``` bash
# 验证k8s上 es集群状态
$ kubectl get pod,svc -n elastic 
NAME                                                   READY   STATUS    RESTARTS   AGE
pod/es-cluster-elasticsearch-client-778df74c8f-7fj4k   1/1     Running   0          2m17s
pod/es-cluster-elasticsearch-client-778df74c8f-skh8l   1/1     Running   0          2m3s
pod/es-cluster-elasticsearch-data-0                    1/1     Running   0          25m
pod/es-cluster-elasticsearch-data-1                    1/1     Running   0          11m
pod/es-cluster-elasticsearch-master-0                  1/1     Running   0          25m
pod/es-cluster-elasticsearch-master-1                  1/1     Running   0          12m
pod/es-cluster-elasticsearch-master-2                  1/1     Running   0          10m

NAME                                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
service/es-cluster-elasticsearch-client      NodePort    10.68.157.105   <none>        9200:29200/TCP,9300:29300/TCP   25m
service/es-cluster-elasticsearch-discovery   ClusterIP   None            <none>        9300/TCP                        25m

# 验证 es集群本身状态
$ curl $NODE_IP:29200/_cat/health
1539335131 09:05:31 es-on-k8s green 7 2 0 0 0 0 0 0 - 100.0%

$ curl $NODE_IP:29200/_cat/indices?v
health status index uuid pri rep docs.count docs.deleted store.size pri.store.size
root@k8s401:/etc/kubeasz# curl 10.100.97.41:29200/_cat/nodes?
172.31.2.4 27 80 5 0.09 0.11 0.21 mi - es-cluster-elasticsearch-master-0
172.31.1.7 30 97 3 0.39 0.29 0.27 i  - es-cluster-elasticsearch-client-778df74c8f-skh8l
172.31.3.7 20 97 3 0.11 0.17 0.18 i  - es-cluster-elasticsearch-client-778df74c8f-7fj4k
172.31.1.5  8 97 5 0.39 0.29 0.27 di - es-cluster-elasticsearch-data-0
172.31.2.5  8 80 3 0.09 0.11 0.21 di - es-cluster-elasticsearch-data-1
172.31.1.6 18 97 4 0.39 0.29 0.27 mi - es-cluster-elasticsearch-master-2
172.31.3.6 20 97 4 0.11 0.17 0.18 mi * es-cluster-elasticsearch-master-1
```

### es 性能压测

如上已使用 chart 在 k8s上部署了 **7** 节点的 elasticsearch 集群；各位应该十分好奇性能怎么样；官方提供了压测工具[esrally](https://github.com/elastic/rally)可以方便的进行性能压测，这里省略安装和测试过程；压测机上执行：  
`esrally --track=http_logs --target-hosts="$NODE_IP:29200" --pipeline=benchmark-only --report-file=report.md`  
压测过程需要1-2个小时，部分压测结果如下：  

``` bash
------------------------------------------------------
    _______             __   _____
   / ____(_)___  ____ _/ /  / ___/_________  ________
  / /_  / / __ \/ __ `/ /   \__ \/ ___/ __ \/ ___/ _ \
 / __/ / / / / / /_/ / /   ___/ / /__/ /_/ / /  /  __/
/_/   /_/_/ /_/\__,_/_/   /____/\___/\____/_/   \___/
------------------------------------------------------

|   Lap |                               Metric |         Task |       Value |    Unit |
|------:|-------------------------------------:|-------------:|------------:|--------:|
...
|   All |                       Min Throughput | index-append |     16903.2 |  docs/s |
|   All |                    Median Throughput | index-append |     17624.4 |  docs/s |
|   All |                       Max Throughput | index-append |     19382.8 |  docs/s |
|   All |              50th percentile latency | index-append |     1865.74 |      ms |
|   All |              90th percentile latency | index-append |     3708.04 |      ms |
|   All |              99th percentile latency | index-append |     6379.49 |      ms |
|   All |            99.9th percentile latency | index-append |     8389.74 |      ms |
|   All |           99.99th percentile latency | index-append |     9612.84 |      ms |
|   All |             100th percentile latency | index-append |     9861.02 |      ms |
|   All |         50th percentile service time | index-append |     1865.74 |      ms |
|   All |         90th percentile service time | index-append |     3708.04 |      ms |
|   All |         99th percentile service time | index-append |     6379.49 |      ms |
|   All |       99.9th percentile service time | index-append |     8389.74 |      ms |
|   All |      99.99th percentile service time | index-append |     9612.84 |      ms |
|   All |        100th percentile service time | index-append |     9861.02 |      ms |
|   All |                           error rate | index-append |           0 |       % |
|   All |                       Min Throughput |      default |        0.66 |   ops/s |
|   All |                    Median Throughput |      default |        0.66 |   ops/s |
|   All |                       Max Throughput |      default |        0.66 |   ops/s |
|   All |              50th percentile latency |      default |      770131 |      ms |
|   All |              90th percentile latency |      default |      825511 |      ms |
|   All |              99th percentile latency |      default |      838030 |      ms |
|   All |             100th percentile latency |      default |      839382 |      ms |
|   All |         50th percentile service time |      default |      1539.4 |      ms |
|   All |         90th percentile service time |      default |     1635.39 |      ms |
|   All |         99th percentile service time |      default |     1728.02 |      ms |
|   All |        100th percentile service time |      default |      1736.2 |      ms |
|   All |                           error rate |      default |           0 |       % |
...
```  

从测试结果看：集群的吞吐可以（k8s es-client pod还可以扩展）；延迟略高一些（因为使用了nfs共享存储）；整体效果不错。

### 中文分词安装

安装 ik 插件即可，可以自定义已安装ik插件的es docker镜像：创建如下 Dockerfile  

``` bash
FROM jmgao1983/elasticsearch:6.4.0

RUN /usr/share/elasticsearch/bin/elasticsearch-plugin install \
  --batch https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v6.4.0/elasticsearch-analysis-ik-6.4.0.zip \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

### 参考阅读

1. [Elasticsearch 入门教程](http://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)
2. [Elasticsearch 压测方案之 esrally 简介](https://segmentfault.com/a/1190000011174694)
