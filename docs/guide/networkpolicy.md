## Network Policy

`Network Policy`提供了基于策略的网络控制，用于隔离应用并减少攻击面。它使用标签选择器模拟传统的分段网络，并通过策略控制它们之间的流量以及来自外部的流量；目前基于`linux iptables`实现，使用类似`nf_conntrack`检查记录网络流量`session`从而决定流量是否阻断；因此它是`状态检测防火墙`。

- 网络插件要支持 Network Policy，如 Calico、Romana、Weave Net

### 简单示例

实验环境：k8s v1.9, calico 2.6.5

首先部署测试用nginx服务

``` bash
$ kubectl run nginx --image=nginx --replicas=3 --port=80 --expose
# 验证测试nginx服务
$ kubectl get pod -o wide 
NAME                     READY     STATUS    RESTARTS   AGE       IP               NODE
nginx-7587c6fdb6-p2fpz   1/1       Running   0          55m       172.20.125.2     10.0.96.7
nginx-7587c6fdb6-pbw7c   1/1       Running   0          55m       172.20.124.2     10.0.96.6
nginx-7587c6fdb6-v48db   1/1       Running   0          55m       172.20.121.195   10.0.96.4
$ kubectl get svc nginx
NAME      TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)   AGE
nginx     ClusterIP   10.68.7.183   <none>        80/TCP    1h
```
默认情况下，其他pod可以访问nginx服务

``` bash
$ kubectl run busy1 --rm -it --image=busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # wget --spider --timeout=1 nginx
Connecting to nginx (10.68.7.183:80)
```
创建`DefaultDeny Network Policy`后，其他Pod（包括namespace外部）不能访问nginx

``` bash
$ cat > default-deny.yaml << EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
  - Ingress
EOF
$ kubectl create -f default-deny.yaml
networkpolicy "default-deny" created
$ kubectl run busy1 --rm -it --image=busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # wget --spider --timeout=1 nginx
Connecting to nginx (10.68.7.183:80)
wget: download timed out
```
创建一个允许带有access=true的Pod访问nginx的网络策略

``` bash
$ cat > nginx-policy.yaml << EOF
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: access-nginx
spec:
  podSelector:
    matchLabels:
      run: nginx
  ingress:
  - from:
    - podSelector:
        matchLabels:
          access: "true"
EOF
$ kubectl create -f nginx-policy.yaml
networkpolicy "access-nginx" created

# 不带access=true标签的Pod还是无法访问nginx服务
$ kubectl run busy1 --rm -it --image=busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # wget --spider --timeout=1 nginx
Connecting to nginx (10.68.7.183:80)
wget: download timed out

# 而带有access=true标签的Pod可以访问nginx服务
$ kubectl run busy2 --rm -it --labels="access=true" --image=busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # wget --spider --timeout=1 nginx
Connecting to nginx (10.68.7.183:80)
```

### 示例策略解读

``` bash
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-network-policy
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: db
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - ipBlock:
        cidr: 172.17.0.0/16
        except:
        - 172.17.1.0/24
    - namespaceSelector:
        matchLabels:
          project: myproject
    - podSelector:
        matchLabels:
          role: frontend
    ports:
    - protocol: TCP
      port: 6379
  egress:
  - to:
    - ipBlock:
        cidr: 10.0.0.0/24
    ports:
    - protocol: TCP
      port: 5978
```
- 策略作用的对象Pods：default命名空间下带有`role=db`标签的Pod
  - 内向流量策略
    - 允许属于`172.17.0.0/16`网段但不属于`172.17.1.0/24`的源地址访问该对象Pods的TCP 6379端口
    - 允许带有project=myprojects标签的namespace中所有Pod访问该对象Pods的TCP 6379端口
    - 允许default命名空间下带有role=frontend标签的Pod访问该对象Pods的TCP 6379端口
    - 拒绝其他所有主动访问该对象Pods的网络流量
  - 外向流量策略
    - 允许该对象Pods主动访问目的地址属于`10.0.0.0/24`网段且目的端口为TCP 5978的流量
    - 拒绝该对象Pods其他所有主动外向网络流量

### 使用场景

参考阅读[ahmetb/kubernetes-network-policy-recipes](https://github.com/ahmetb/kubernetes-network-policy-recipes) 该项目举例一些使用NetworkPolicy的场景，并有形象的配图

#### 拒绝其他namespaces访问服务

![deny_from_other_namespaces](https://github.com/ahmetb/kubernetes-network-policy-recipes/blob/master/img/4.gif)

+ 场景1：你的k8s集群应用按照namespaces区分生产、测试环境，你要确保生产环境不会受到测试环境错误访问影响
+ 场景2：你的k8s集群有多租户应用采用namespaces区分的，你要确保多租户之间的应用隔离

在你需要隔离的命名空间创建如下策略:

``` bash
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  namespace: your-ns
  name: deny-other-namespaces
spec:
  podSelector:
    matchLabels:
  ingress:
  - from:
    - podSelector: {}
```

#### 允许外部访问服务

+ 场景：暴露特定Pod的特定端口给外部访问

![allow_from_external](https://github.com/ahmetb/kubernetes-network-policy-recipes/blob/master/img/8.gif)

``` bash
# 创建示例应用待暴露服务
$ kubectl run web --image=nginx --labels=app=web --port 80 --expose

# 创建网络策略
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: web-allow-external
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - from: []
    ports:
    - protocol: TCP
      port: 80
```
