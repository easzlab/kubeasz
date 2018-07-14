# Metrics Server

从 v1.8 开始，资源使用情况的度量（如容器的 CPU 和内存使用）可以通过 Metrics API 获取；前提是集群中要部署 Metrics Server，它从Kubelet 公开的Summary API采集指标信息，关于更多的背景介绍请参考如下文档：  
- Metrics Server[设计提案](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/metrics-server.md)

大致是说它符合k8s的监控架构设计，受heapster项目启发，并且比heapster优势在于：访问不需要apiserver的代理机制，提供认证和授权等；很多集群内组件依赖它（HPA,scheduler,kubectl top），因此它应该在集群中默认运行；部分k8s集群的安装工具已经默认集成了Metrics Server的安装，以下概述下它的安装：

- 1.metric-server是扩展的apiserver，依赖于[kube-aggregator](https://github.com/kubernetes/kube-aggregator)，因此需要在apiserver中开启相关参数。
- 2.需要在集群中运行deployment处理请求

从kubeasz 0.1.0 开始，metrics-server已经默认集成在集群安装脚本中，请查看`roles/cluster-addon/defaults/main.yml`中的设置

## 安装

默认已集成在90.setup.yml中，如果分步请执行`ansible-play /etc/ansible/07.cluster-addon.yml`

- 1.设置apiserver相关[参数](../../roles/kube-master/templates/kube-apiserver.service.j2)
``` bash
... # 省略
  --requestheader-client-ca-file={{ ca_dir }}/ca.pem \
  --requestheader-allowed-names=aggregator \
  --requestheader-extra-headers-prefix=X-Remote-Extra- \
  --requestheader-group-headers=X-Remote-Group \
  --requestheader-username-headers=X-Remote-User \
  --proxy-client-cert-file={{ ca_dir }}/aggregator-proxy.pem \
  --proxy-client-key-file={{ ca_dir }}/aggregator-proxy-key.pem \
  --enable-aggregator-routing=true \
```
- 2.生成[aggregator proxy相关证书](../../roles/kube-master/tasks/main.yml)

参考1：https://kubernetes.io/docs/tasks/access-kubernetes-api/configure-aggregation-layer/  
参考2：https://kubernetes.io/docs/tasks/access-kubernetes-api/setup-extension-api-server/

## 验证

- 查看生成的新api：v1beta1.metrics.k8s.io
``` bash
$ kubectl get apiservice|grep metrics
v1beta1.metrics.k8s.io                 1d
```

- 查看kubectl top命令（无需额外安装heapster）
``` bash
$ kubectl top node
NAME           CPU(cores)   CPU%      MEMORY(bytes)   MEMORY%   
192.168.1.1   116m         2%        2342Mi          60%       
192.168.1.2   79m          1%        1824Mi          47%       
192.168.1.3   82m          2%        1897Mi          49%  
$ kubectl top pod --all-namespaces 	# 输出略
```

- 验证基于metrics-server实现的基础hpa自动缩放，请参考[hpa.md](hpa.md)

## 补充

目前dashboard插件如果想在界面上显示资源使用率，它还依赖于`heapster`；另外，测试发现k8s 1.8版本的`kubectl top`也依赖`heapster`，因此建议补充安装`heapster`，无需安装`influxdb`和`grafana`。

``` bash
$ kubectl apply -f /etc/ansible/manifests/heapster/heapster.yaml
```
