## Ingress简介

ingress就是从kubernetes集群外访问集群的入口，将用户的URL请求转发到不同的service上。ingress相当于nginx反向代理服务器，它包括的规则定义就是URL的路由信息；它的实现需要部署`Ingress controller`(比如 [traefik](https://github.com/containous/traefik) [ingress-nginx](https://github.com/kubernetes/ingress-nginx) 等)，`Ingress controller`通过apiserver监听ingress和service的变化，并根据规则配置负载均衡并提供访问入口，达到服务发现的作用。

+ 未配置ingress：

集群外部 -> NodePort -> K8S Service

+ 配置ingress:

集群外部 -> Ingress -> K8S Service

+ 注意：ingress 本身也需要部署`Ingress controller`时暴露`NodePort`让外部访问

### 部署 Traefik

Traefik 提供了一个简单好用 `Ingress controller`，下文基于它讲解一个简单的 ingress部署和测试例子。请查看yaml配置 [traefik-ingress.yaml](../../manifests/ingress/traefik-ingress.yaml)，参考[traefik 官方k8s例子](https://github.com/containous/traefik/tree/master/examples/k8s)

#### 安装 traefik ingress-controller

``` bash
kubectl create -f /etc/ansible/manifests/ingress/traefik-ingress.yaml
```
+ 注意需要配置 `RBAC`授权
+ 注意trafik `Service`中 `80`端口为 traefik ingress-controller的服务端口，`8080`端口为 traefik 的管理WEB界面；为后续配置方便指定`80` 端口暴露`NodePort`端口为 `23456`(对应于在hosts配置中`NODE_PORT_RANGE`范围内可用端口)

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
kubectl run test-hello --image=nginx --expose --port=80
##
# kubectl get deploy test-hello
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
test-hello   1         1         1            1           56s
# kubectl get svc test-hello
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
test-hello   ClusterIP   10.68.124.115   <none>        80/TCP    1m
```
+ 然后为这个应用创建 ingress，`kubectl create -f /etc/ansible/manifests/ingress/test-hello.ing.yaml`

``` bash
# test-hello.ing.yaml内容
apiVersion: extensions/v1beta1
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
+ 集群内部尝试访问: `curl -H Host:hello.test.com 10.68.69.170(traefik-ingress-service的服务地址)` 能够看到欢迎页面 `Welcome to nginx!`；在集群外部尝试访问(假定集群一个NodeIP为 192.168.1.1): `curl -H Host:hello.test.com 192.168.1.1:23456`，也能够看到欢迎页面 `Welcome to nginx!`，说明ingress测试成功

+ 最后我们可以为traefik WEB管理页面也创建一个ingress, `kubectl create -f /etc/ansible/manifests/ingress/traefik-ui.ing.yaml`

``` bash
# traefik-ui.ing.yaml内容
---
apiVersion: extensions/v1beta1
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
这样在集群外部可以使用 `curl -H Host:traefik-ui.test.com 192.168.1.1:23456` 尝试访问WEB管理页面，返回 `<a href="/dashboard/">Found</a>.`说明 traefik-ui的ingress配置生效了。

### [可选] 部署`ingress-service`的代理

在客户端主机上可以通过修改本机 `hosts` 文件，如上例子，增加两条记录：

``` text
192.168.1.1	hello.test.com
192.168.1.1	traefik-ui.test.com
```
打开浏览器输入域名 `http://hello.test.com:23456` 和 `http://traefik-ui.test.com:23456` 就可以访问k8s的应用服务了。

当然如果你的环境中有类似 nginx/haproxy 等代理，可以做代理转发以去掉 `23456`这个端口，这里以 haproxy演示下。

如果你的集群根据本项目部署了高可用方案，那么可以利用`LB` 节点haproxy 来做，当然如果生产环境K8S应用已经部署非常多，建议还是使用独立的 `nginx/haproxy`集群

在 LB 主备节点，修改 `/etc/haproxy/haproxy.cfg`类似如下：

``` bash
global
        log /dev/log    local0
        log /dev/log    local1 notice
        chroot /var/lib/haproxy
        stats socket /run/haproxy/admin.sock mode 660 level admin
        stats timeout 30s
        user haproxy
        group haproxy
        daemon
        nbproc 1

defaults
        log     global
        timeout connect 5000
        timeout client  50000
        timeout server  50000

listen kube-master
        bind 0.0.0.0:8443
        mode tcp
        option tcplog
        balance source
        # 根据实际kube-master 节点数量增减如下endpoints
        server s1 192.168.1.1:6443  check inter 10000 fall 2 rise 2 weight 1
        server s2 192.168.1.2:6443  check inter 10000 fall 2 rise 2 weight 1

listen kube-node
	# 先确认 LB节点80端口可用
        bind 0.0.0.0:80		
        mode tcp
        option tcplog
        balance source
        # 根据实际kube-node 节点数量增减如下endpoints
        server s1 192.168.1.1:23456  check inter 10000 fall 2 rise 2 weight 1
        server s2 192.168.1.2:23456  check inter 10000 fall 2 rise 2 weight 1
        server s3 192.168.1.3:23456  check inter 10000 fall 2 rise 2 weight 1
```
修改保存后，重启haproxy服务；

这样我们就可以访问集群`master-VIP`的`80`端口，由haproxy代理转发到实际的node节点和nodePort端口上了。这时可以修改客户端本机 `hosts`文件如下：(假定 master-VIP=192.168.1.10)

``` text
192.168.1.10     hello.test.com
192.168.1.10    traefik-ui.test.com
```
打开浏览器输入域名 `http://hello.test.com` 和 `http://traefik-ui.test.com`可以正常访问。


