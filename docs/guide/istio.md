---
title: "Istio 1.1.7 安装　"
date: 2019-05-19T19:44:00+08:00
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

进入 [Istio release](https://github.com/istio/istio/releases) 页面下载最新版安装包并解压到当前目录,

```sh
curl -L https://git.io/getLatestIstio | sh -


ll istio-1.1.7/
total 40
drwxr-xr-x  2 root root  4096 May 15 08:59 bin
drwxr-xr-x  6 root root  4096 May 15 08:59 install
-rw-r--r--  1 root root   602 May 15 08:59 istio.VERSION
-rw-r--r--  1 root root 11343 May 15 08:59 LICENSE
-rw-r--r--  1 root root  5921 May 15 08:59 README.md
drwxr-xr-x 15 root root  4096 May 15 08:59 samples
drwxr-xr-x  7 root root  4096 May 15 08:59 tools
```
- install  Kubernetes 安装所需的 .yaml 文件
- samples  Task中的示例应用
- bin/istioctl 客户端工具
- istio.VERSION 配置文件

#### 安装 
---- 

注意事项

- Node 节点内存不能低于 4G，否则相关容器可能启动失败  
- Istio 默认使用‘负载均衡器’服务对象类型。对于裸机安装没有负载均衡器的情况下，安装需指定‘NodePort’类型。


##### 方案1：使用 Helm template 进行安装

```bash
cd /usr/local/src/istio-1.1.7

kubectl create namespace istio-system

# 安装 istio-init chart，来启动 Istio CRD 的安装过程
helm template install/kubernetes/helm/istio-init --name istio-init --namespace istio-system --set gateways.istio-ingressgateway.type=NodePort --set gateways.istio-egressgateway.type=NodePort | kubectl apply -f -

# 稍等一会儿执行
# 输出 23 或者 28 （若开启了 cert-manager）
kubectl get crds | grep 'istio.io\|certmanager.k8s.io' | wc -l

# 部署与你选择的配置文件相对应的 Istio 的核心组件
# 不同配置说明 https://istio.io/zh/docs/setup/kubernetes/additional-setup/config-profiles/

# 选择 default 配置
helm template install/kubernetes/helm/istio --name istio --namespace istio-system \
  --set gateways.istio-ingressgateway.type=NodePort \
  --set gateways.istio-egressgateway.type=NodePort | kubectl apply -f -

# 或选择 demo 配置
helm template install/kubernetes/helm/istio --name istio --namespace istio-system \
  --set gateways.istio-ingressgateway.type=NodePort \
  --set gateways.istio-egressgateway.type=NodePort \
  --values install/kubernetes/helm/istio/values-istio-demo.yaml | kubectl apply -f -
```

##### 方案2：在 Helm 和 Tiller 的环境中使用 helm install 命令进行安装

见[官方文档](https://istio.io/zh/docs/setup/kubernetes/install/helm/#%E6%96%B9%E6%A1%88-2-%E5%9C%A8-helm-%E5%92%8C-tiller-%E7%9A%84%E7%8E%AF%E5%A2%83%E4%B8%AD%E4%BD%BF%E7%94%A8-helm-install-%E5%91%BD%E4%BB%A4%E8%BF%9B%E8%A1%8C%E5%AE%89%E8%A3%85)


##### 验证
```bash
kubectl get pod -n istio-system

# default 配置时
NAME                                     READY   STATUS    RESTARTS   AGE
istio-citadel-899dfb67c-5hlsc             1/1     Running     0          49s
istio-cleanup-secrets-1.1.7-nkdxt         0/1     Completed   0          50s
istio-galley-555dd7c7d7-rpfln             1/1     Running     0          49s
istio-ingressgateway-5b547dfb7b-ctm5l     1/1     Running     0          49s
istio-init-crd-10-l9xcj                   0/1     Completed   0          66s
istio-init-crd-11-nqvml                   0/1     Completed   0          66s
istio-pilot-9f5c75ddf-n5s6p               2/2     Running     0          49s
istio-policy-bd45d757d-6qcdg              2/2     Running     1          49s
istio-security-post-install-1.1.7-nbwwv   0/1     Completed   0          50s
istio-sidecar-injector-998dd6cbb-n2hdm    1/1     Running     0          49s
istio-telemetry-656df5b64-k8vkf           2/2     Running     1          49s
prometheus-7f87866f5f-t97wc               1/1     Running     0          49s

# demo 配置时
grafana-749c78bcc5-fbzmn                  1/1     Running     0          101s
istio-citadel-899dfb67c-8shx2             1/1     Running     0          100s
istio-cleanup-secrets-1.1.7-jbhsl         0/1     Completed   0          102s
istio-egressgateway-748d5fd794-x5bjt      1/1     Running     0          101s
istio-galley-555dd7c7d7-86r2b             1/1     Running     0          101s
istio-grafana-post-install-1.1.7-kq7b4    0/1     Completed   0          103s
istio-ingressgateway-55dd86767f-jd9m4     1/1     Running     0          101s
istio-init-crd-10-l9xcj                   0/1     Completed   0          16m
istio-init-crd-11-nqvml                   0/1     Completed   0          16m
istio-pilot-6964dd4957-7bzdq              2/2     Running     0          101s
istio-policy-689687bd77-ncw2n             2/2     Running     1          101s
istio-security-post-install-1.1.7-t2kwh   0/1     Completed   0          102s
istio-sidecar-injector-998dd6cbb-7mwkh    1/1     Running     0          100s
istio-telemetry-8564679887-59c8z          2/2     Running     1          101s
istio-tracing-595796cf54-jn49s            1/1     Running     0          100s
kiali-5df77dc9b6-psjs4                    1/1     Running     0          101s
prometheus-7f87866f5f-hrbgt               1/1     Running     0          100s

```

```bash
kubectl get svc -n istio-system

# default 配置时
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                        AGE
istio-citadel            ClusterIP   10.68.236.249   <none>        8060/TCP,15014/TCP             75s
istio-galley             ClusterIP   10.68.105.102   <none>        443/TCP,15014/TCP,9901/TCP     75s
istio-ingressgateway     NodePort    10.68.181.46    <none>        15020:32761/TCP,80:31380/TCP,443:31390/TCP,31400:31400/TCP,15029:33185/TCP,15030:20745/TCP,15031:36208/TCP,15032:34095/TCP,15443:36244/TCP   75s
istio-pilot              ClusterIP   10.68.252.143   <none>        15010/TCP,15011/TCP,8080/TCP,15014/TCP   75s
istio-policy             ClusterIP   10.68.40.51     <none>        9091/TCP,15004/TCP,15014/TCP   75s
istio-sidecar-injector   ClusterIP   10.68.55.134    <none>        443/TCP                        74s
istio-telemetry          ClusterIP   10.68.16.11     <none>        9091/TCP,15004/TCP,15014/TCP,42422/TCP       75s
prometheus               ClusterIP   10.68.65.238    <none>        9090/TCP                       75s

# demo 配置时
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                        AGE
grafana                  ClusterIP   10.68.65.248    <none>        3000/TCP                       2m27s
istio-citadel            ClusterIP   10.68.72.100    <none>        8060/TCP,15014/TCP             2m26s
istio-egressgateway      NodePort    10.68.21.24     <none>        80:26775/TCP,443:28249/TCP,15443:38494/TCP  2m27s
istio-galley             ClusterIP   10.68.73.9      <none>        443/TCP,15014/TCP,9901/TCP     2m27s
istio-ingressgateway     NodePort    10.68.122.190   <none>        15020:39248/TCP,80:31380/TCP,443:31390/TCP,31400:31400/TCP,15029:33522/TCP,15030:26010/TCP,15031:27064/TCP,15032:32158/TCP,15443:30848/TCP   2m27s
istio-pilot              ClusterIP   10.68.116.5     <none>        15010/TCP,15011/TCP,8080/TCP,15014/TCP  2m26s
istio-policy             ClusterIP   10.68.239.246   <none>        9091/TCP,15004/TCP,15014/TCP   2m27s
istio-sidecar-injector   ClusterIP   10.68.93.151    <none>        443/TCP                        2m26s
istio-telemetry          ClusterIP   10.68.117.254   <none>        9091/TCP,15004/TCP,15014/TCP,42422/TCP  2m26s
jaeger-agent             ClusterIP   None            <none>        5775/UDP,6831/UDP,6832/UDP     2m25s
jaeger-collector         ClusterIP   10.68.103.8     <none>        14267/TCP,14268/TCP            2m26s
jaeger-query             ClusterIP   10.68.73.252    <none>        16686/TCP                      2m26s
kiali                    ClusterIP   10.68.214.228   <none>        20001/TCP                      2m27s
prometheus               ClusterIP   10.68.203.209   <none>        9090/TCP                       2m26s
tracing                  ClusterIP   10.68.113.236   <none>        80/TCP                         2m25s
zipkin                   ClusterIP   10.68.96.189    <none>        9411/TCP                       2m25s
```

##### Sidecar 的自动注入

注意事项

需要在kube-apiserver 启动 admission-control 参数中加入 MutatingAdmissionWebhook 和 ValidatingAdmissionWebhook并确保正确的顺序,如果是多master安装，确保每个kube-apiserver都要进行修改。


##### 部署应用验证

istio 的samples目录中有很多示例。我们现在使用samples/sleep/sleep.yaml 来验证刚刚开启的Sidecar自动注入功能。

进入目录 istio-1.1.7/ 部署一个新的应用

```bash
cd istio-1.1.7/
kubectl apply -f samples/sleep/sleep.yaml

kubectl get pod 
NAME                            READY   STATUS    RESTARTS   AGE
sleep-7549f66447-wv8cl          1/1     Running   0          1m
```

一切都是熟悉的味道。下面给 default 命名空间设置标签：istio-injection=enabled，这样就会在pod 创建时触发 Sidecar 的注入过程。从此default 名称空间拥有了超能力.

```bash
kubectl label namespace default istio-injection=enabled
kubectl get namespace -L istio-injection
NAME           STATUS   AGE     ISTIO-INJECTION
default        Active   1h    enabled
istio-system   Active   3d22h   
kube-public    Active   4d2h    
kube-system    Active   4d2h
```
接下来删除上面创建的pod，观察下有什么变化。

```bash
kubectl delete pod sleep-7549f66447-wv8cl
pod "sleep-7549f66447-wv8cl" deleted

kubectl get pod 
NAME                            READY   STATUS    RESTARTS   AGE
sleep-7549f66447-x4td6          2/2     Running   0          37s
```

刚刚的pod里面现在已经拥有两个容器，进入pod一探究竟。
```bash
 kubectl describe pod sleep-7549f66447-x4td6

 ....
  Containers:
   sleep:
    Container ID:   docker://
    Image:         pstauffer/curl
    .... 
   
   istio-proxy:
    Container ID:   docker://
    Image:         docker.io/istio/proxyv2:1.1.7
    ....
    
```
多出了一个 `istio-proxy` 容器及其对应的存储卷


#### 卸载istio 

---

```bash
# 采用 default 配置安装
helm template install/kubernetes/helm/istio --name istio --namespace istio-system | kubectl delete -f -
# 采用 demo 配置安装
helm template install/kubernetes/helm/istio --name istio --namespace istio-system \
  --values install/kubernetes/helm/istio/values-istio-demo.yaml | kubectl delete -f -

kubectl delete namespace istio-system
```


#### 资源
- [官方安装文档](https://istio.io/zh/docs/setup/kubernetes/install/helm/)
