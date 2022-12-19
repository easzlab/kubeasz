## Ingress简介

本文档已过期，安装最新版本，请参考相关官方文档。

ingress就是从外部访问k8s集群的入口，将用户的URL请求转发到不同的service上。ingress相当于nginx反向代理服务器，它包括的规则定义就是URL的路由信息；它的实现需要部署`Ingress controller`(比如 [traefik](https://github.com/containous/traefik) [ingress-nginx](https://github.com/kubernetes/ingress-nginx) 等)，`Ingress controller`通过apiserver监听ingress和service的变化，并根据规则配置负载均衡并提供访问入口，达到服务发现的作用。

- 未配置ingress：

集群外部 -> NodePort -> K8S Service

- 配置ingress:

集群外部 -> Ingress -> K8S Service

- **注意：ingress 本身也需要部署`Ingress controller`时使用以下几种方式让外部访问**
  - 使用`NodePort`方式
  - 使用`hostPort`方式
  - 使用LoadBalancer地址方式

- 以下讲解基于`Traefik`，如果想要了解`ingress-nginx`的原理与实践，推荐阅读博客[烂泥行天下](https://www.ilanni.com/?p=14501)的相关文章

### 部署 Traefik

Traefik 提供了一个简单好用 `Ingress controller`，下文侧重讲解 ingress部署和测试例子。请查看yaml配置 [traefik-ingress.yaml](../../manifests/ingress/traefik/traefik-ingress.yaml)，参考[traefik 官方k8s例子](https://github.com/containous/traefik/tree/master/examples/k8s)

#### 安装 traefik ingress-controller

``` bash
kubectl create -f /etc/kubeasz/manifests/ingress/traefik/traefik-ingress.yaml
```
+ 注意需要配置 `RBAC`授权
+ 注意`trafik pod`中 `80`端口为 traefik ingress-controller的服务端口，`8080`端口为 traefik 的管理WEB界面；为后续配置方便指定`80` 端口暴露`NodePort`端口为 `23456`(对应于在hosts配置中`NODE_PORT_RANGE`范围内可用端口)

#### 验证 traefik ingress-controller

``` bash
# kubectl get deploy -n kube-system traefik-ingress-controller
NAME                         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
traefik-ingress-controller   1         1         1            1           4m

# kubectl get svc -n kube-system traefik-ingress-service
NAME                      TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)                       AGE
traefik-ingress-service   NodePort   10.68.69.170   <none>        80:23456/TCP,8080:34815/TCP   4m
```
+ 可以看到`traefik-ingress-service` 服务端口`80`暴露的nodePort确实为`23456`

#### 测试 ingress

+ 首先创建测试用K8S应用，并且该应用服务不用nodePort暴露，而是用ingress方式让外部访问

``` bash
kubectl run test-hello --image=nginx:alpine --expose --port=80
##
# kubectl get deploy test-hello
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
test-hello   1         1         1            1           56s
# kubectl get svc test-hello
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
test-hello   ClusterIP   10.68.124.115   <none>        80/TCP    1m
```
+ 然后为这个应用创建 ingress，`kubectl create -f /etc/kubeasz/manifests/ingress/test-hello.ing.yaml`

``` bash
# test-hello.ing.yaml内容
apiVersion: networking.k8s.io/v1beta1 
kind: Ingress
metadata:
  name: test-hello
spec:
  rules:
  - host: hello.test.com
    http:
      paths:
      - path: /
        backend:
          serviceName: test-hello
          servicePort: 80
```
+ 集群内部尝试访问: `curl -H Host:hello.test.com 10.68.69.170(traefik-ingress-service的服务地址)` 能够看到欢迎页面 `Welcome to nginx!`；
+ 在集群外部尝试访问(假定集群一个NodeIP为 192.168.1.1): `curl -H Host:hello.test.com 192.168.1.1:23456`，也能够看到欢迎页面 `Welcome to nginx!`，说明ingress测试成功

#### 为 traefik WEB 管理页面创建 ingress 规则 

`kubectl create -f /etc/kubeasz/manifests/ingress/traefik/traefik-ui.ing.yaml`

``` bash
# traefik-ui.ing.yaml内容
---
apiVersion: networking.k8s.io/v1beta1 
kind: Ingress
metadata:
  name: traefik-web-ui
  namespace: kube-system
spec:
  rules:
  - host: traefik-ui.test.com
    http:
      paths:
      - path: /
        backend:
          serviceName: traefik-ingress-service
          servicePort: 8080
```

+ 在集群外部可以使用 `curl -H Host:traefik-ui.test.com 192.168.1.1:23456` 尝试访问WEB管理页面，返回 `<a href="/dashboard/">Found</a>.`说明 traefik-ui的ingress配置生效了。

+ 在客户端主机也可以通过修改本机 `hosts` 文件，如上例子，增加两条记录：

``` text
192.168.1.1	hello.test.com
192.168.1.1	traefik-ui.test.com
```
打开浏览器输入域名 `http://hello.test.com:23456` 和 `http://traefik-ui.test.com:23456` 就可以访问k8s的应用服务了。

### 可选1: 使用`LoadBalancer`服务类型来暴露ingress，自有环境（非公有云）可以参考[metallb文档](metallb.md)

``` bash
# 修改traefik-ingress 使用 LoadBalancer服务
$ sed -i 's/NodePort$/LoadBalancer/g' /etc/kubeasz/manifests/ingress/traefik/traefik-ingress.yaml
# 创建traefik-ingress
$ kubectl apply -f /etc/kubeasz/manifests/ingress/traefik/traefik-ingress.yaml
# 验证
$ kubectl get svc --all-namespaces |grep traefik
kube-system   traefik-ingress-service   LoadBalancer   10.68.163.243   192.168.1.241   80:23456/TCP,8080:37088/TCP   1m
```
这时可以修改客户端本机 `hosts`文件：(如上例192.168.1.241)

``` text
192.168.1.241     hello.test.com
192.168.1.241     traefik-ui.test.com
```
打开浏览器输入域名 `http://hello.test.com` 和 `http://traefik-ui.test.com`可以正常访问。

### 可选2: 部署`ingress-service`的负载均衡

- 利用 nginx/haproxy 等集群，可以做代理转发以去掉 `23456`这个端口。如果你的集群根据本项目部署了高可用方案，那么可以利用`LB` 节点haproxy 来做，当然如果生产环境K8S应用已经部署非常多，建议还是使用独立的 `nginx/haproxy`集群。

具体参考[配置转发 ingress nodePort](../op/loadballance_ingress_nodeport.md)，如上配置访问集群`MASTER_IP`的`80`端口时，由haproxy代理转发到实际的node节点暴露的nodePort端口上了。这时可以修改客户端本机 `hosts`文件如下：(假定 MASTER_IP=192.168.1.10)

``` text
192.168.1.10     hello.test.com
192.168.1.10    traefik-ui.test.com
```
打开浏览器输入域名 `http://hello.test.com` 和 `http://traefik-ui.test.com`可以正常访问。

## 下一步[配置https ingress](ingress-tls.md)
