---
title: "Istio 1.0.3 安装　"
date: 2018-11-12T13:44:34+08:00
draft: false
---

#### Service Mesh(服务网格)  

--- 
Kubernetes 已经给我们带来了诸多的好处。但是仍有些需求比如 A/B 测试、金丝雀发布、限流、访问控制,端到端认证等需要运维人员进一步去解决。

Istio 是完全开源的服务网格,提供了一套完整的解决方案，可以透明地分层到现有的分布式应用程序上。对开发人员几乎无感的同时获得超能力。

如果想要现有的服务支持 Istio，只需要在当前的环境中部署一个特殊的 sidecar 代理，即可。

##### 前提   

---- 

- 安装 Kubernetes 集群 1.9+ 
- [安装 Helm](./helm.md) 

##### 准备

---- 

进入 [Istio release](https://github.com/istio/istio/releases) 页面下载最新版安装包(1.0.3)并解压到当前目录,

```
curl -L https://git.io/getLatestIstio | sh -

ll istio-1.0.3/
total 28
drwxr-xr-x  2 root root    22 10月 26 07:36 bin
drwxr-xr-x  6 root root    79 10月 26 07:36 install
-rw-r--r--  1 root root   648 10月 26 07:36 istio.VERSION
-rw-r--r--  1 root root 11343 10月 26 07:36 LICENSE
-rw-r--r--  1 root root  5817 10月 26 07:36 README.md
drwxr-xr-x 12 root root   212 10月 26 07:36 samples
drwxr-xr-x  8 root root  4096 10月 26 07:36 tools
```
- install  Kubernetes 安装所需的 .yaml 文件
- samples  Task中的示例应用
- bin/istioctl 客户端工具
- istio.VERSION 配置文件

#### 安装 

---

##### 安装　istio
注意事项

Istio 默认使用‘负载均衡器’服务对象类型。对于裸机安装没有负载均衡器的情况下，安装需指定‘NodePort’类型。

```
helm install --name istio install/kubernetes/helm/istio --namespace istio-system --set gateways.istio-ingressgateway.type=NodePort --set gateways.istio-egressgateway.type=NodePort
```

##### 验证
```
kubectl get pod -n istio-system
NAME                                     READY   STATUS    RESTARTS   AGE
istio-citadel-6955bc9cb7-qh846           1/1     Running   0          3d22h
istio-egressgateway-7dc5cbbc56-k4cgh     1/1     Running   0          3d22h
istio-galley-545b6b8f5b-k7ssx            1/1     Running   0          3d22h
istio-ingressgateway-7958d776b5-ptdsc    1/1     Running   0          3d22h
istio-pilot-56bfdbffff-mtcn6             2/2     Running   0          3d22h
istio-policy-5c689f446f-6bzlq            2/2     Running   0          3d15h
istio-policy-5c689f446f-dvmfq            2/2     Running   0          3d22h
istio-policy-5c689f446f-f2kl8            2/2     Running   0          3d3h
istio-policy-5c689f446f-nfv2l            2/2     Running   0          3d1h
istio-policy-5c689f446f-qdtql            2/2     Running   0          3d2h
istio-sidecar-injector-99b476b7b-dt24k   1/1     Running   0          3d22h
istio-telemetry-55d68b5dfb-52ftl         2/2     Running   0          3d22h
istio-telemetry-55d68b5dfb-dvdvz         2/2     Running   0          3d22h
istio-telemetry-55d68b5dfb-ln2sr         2/2     Running   0          3d
istio-telemetry-55d68b5dfb-m2mb8         2/2     Running   0          3d
istio-telemetry-55d68b5dfb-sjgq8         2/2     Running   0          3d
prometheus-65d6f6b6c-dsv26               1/1     Running   0          3d22h

```
```
kubectl get svc -n istio-system
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                                                                   AGE
istio-citadel            ClusterIP   10.68.7.100     <none>        8060/TCP,9093/TCP                                                                                                         3d22h
istio-egressgateway      NodePort    10.68.67.237    <none>        80:30060/TCP,443:38194/TCP                                                                                                3d22h
istio-galley             ClusterIP   10.68.12.54     <none>        443/TCP,9093/TCP                                                                                                          3d22h
istio-ingressgateway     NodePort    10.68.87.79     <none>        80:31380/TCP,443:31390/TCP,31400:31400/TCP,15011:31812/TCP,8060:30957/TCP,853:23011/TCP,15030:22292/TCP,15031:23663/TCP   3d22h
istio-pilot              ClusterIP   10.68.84.101    <none>        15010/TCP,15011/TCP,8080/TCP,9093/TCP                                                                                     3d22h
istio-policy             ClusterIP   10.68.94.206    <none>        9091/TCP,15004/TCP,9093/TCP                                                                                               3d22h
istio-sidecar-injector   ClusterIP   10.68.191.221   <none>        443/TCP                                                                                                                   3d22h
istio-telemetry          ClusterIP   10.68.199.8     <none>        9091/TCP,15004/TCP,9093/TCP,42422/TCP                                                                                     3d22h
prometheus               ClusterIP   10.68.91.13     <none>        9090/TCP    
```

##### Sidecar 的自动注入

注意事项

需要在kube-apiserver 启动 admission-control 参数中加入 MutatingAdmissionWebhook 和 ValidatingAdmissionWebhook并确保正确的顺序,如果是多master安装，确保每个kube-apiserver都要进行修改。

```
/bin/kube-apiserver --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook
```

重启 kube-apiserver 即可

##### 部署应用验证

istio 的samples目录中有很多示例。我们现在使用samples/sleep/sleep.yaml 来验证刚刚开启的Sidecar自动注入功能。

进入目录 istio-1.0.3/ 部署一个新的应用

```
cd istio-1.0.3/
kubectl apply -f samples/sleep/sleep.yaml

kubectl get pod 
NAME                            READY   STATUS    RESTARTS   AGE
sleep-7549f66447-wv8cl          1/1     Running   0          1m
```

一切都是熟悉的味道。下面给 default 命名空间设置标签：istio-injection=enabled，这样就会在pod 创建时触发 Sidecar 的注入过程。从此default 名称空间拥有了超能力.

```
kubectl label namespace default istio-injection=enabled
kubectl get namespace -L istio-injection
NAME           STATUS   AGE     ISTIO-INJECTION
default        Active   1h    enabled
istio-system   Active   3d22h   
kube-public    Active   4d2h    
kube-system    Active   4d2h
```
接下来删除上面创建的pod，观察下有什么变化。

```
kubectl delete pod sleep-7549f66447-wv8cl
pod "sleep-7549f66447-wv8cl" deleted

kubectl get pod 
NAME                            READY   STATUS    RESTARTS   AGE
sleep-7549f66447-x4td6          2/2     Running   0          37s
```
刚刚的pod里面现在已经拥有两个容器，进入pod一探究竟。
```
 kubectl describe pod sleep-7549f66447-x4td6

 ....

  Containers:
   sleep:
   
     .... 
   
   istio-proxy:
 
     ....
    
```
多出了一个 istio-proxy 容器及其对应的存储卷


#### 卸载istio 

---

```
helm delete --purge istio

```

