## Horizontal Pod Autoscaling

自动水平伸缩，是指运行在k8s上的应用负载(POD)，可以根据资源使用率进行自动扩容、缩容；我们知道应用的资源使用率通常都有高峰和低谷，所以k8s的`HPA`特性应运而生；它也是最能体现区别于传统运维的优势之一，不仅能够弹性伸缩，而且完全自动化！

根据 CPU 使用率或自定义 metrics 自动扩展 Pod 数量（支持 replication controller、deployment）；k8s1.6版本之前是通过kubelet来获取监控指标，1.6版本之后是通过api server、heapster或者kube-aggregator来获取监控指标。

### Metrics支持

根据不同版本的API中，HPA autoscale时靠以下指标来判断资源使用率：
- autoscaling/v1: CPU
- autoscaling/v2alpha1
  - 内存
  - 自定义metrics
  - 多metrics组合: 根据每个metric的值计算出scale的值，并将最大的那个值作为扩容的最终结果

### 基础示例

本实验环境基于k8s 1.8 和 1.9，仅使用`autoscaling/v1` 版本API，**注意确保**`k8s` 集群插件`kubedns` 和 `heapster` 工作正常。

``` bash
# 创建deploy和service
$ kubectl run php-apache --image=pilchard/hpa-example --requests=cpu=200m --expose --port=80

# 创建autoscaler
$ kubectl autoscale deploy php-apache --cpu-percent=50 --min=1 --max=10

# 等待3~5分钟查看hpa状态
$ kubectl get hpa php-apache
NAME         REFERENCE               TARGETS    MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   0% / 50%   1         10        1          3m

# 增加负载
$ kubectl run --rm -it load-generator --image=busybox /bin/sh
Hit enter for command prompt
$ while true; do wget -q -O- http://php-apache; done;

# 等待约5分钟查看hpa显示负载增加，且副本数目增加为4
$ kubectl get hpa php-apache
NAME         REFERENCE               TARGETS      MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   430% / 50%   1         10        4          4m

# 注意k8s为了避免频繁增删pod，对副本的增加速度有限制
# 实验过程可以看到副本数目从1到4到8到10，大概都需要4~5分钟的缓冲期
$ kubectl get hpa php-apache
NAME         REFERENCE               TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   86% / 50%   1         10        8          9m
$ kubectl get hpa php-apache
NAME         REFERENCE               TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   52% / 50%   1         10        10         12m

# 清除负载，CTRL+C 结束上述循环程序，稍后副本数目变回1
$ kubectl get hpa php-apache
NAME         REFERENCE               TARGETS    MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   0% / 50%   1         10        1          17m
```

