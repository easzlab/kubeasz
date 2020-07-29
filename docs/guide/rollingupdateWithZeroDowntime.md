## 1、前言
在当下微服务架构盛行的时代，用户希望应用程序时时刻刻都是可用，为了满足不断变化的新业务，需要不断升级更新应用程序，有时可能需要频繁的发布版本。实现"零停机"、“零感知”的持续集成(Continuous Integration)和持续交付/部署(Continuous Delivery)应用程序，一直都是软件升级换代不得不面对的一个难题和痛点，也是一种追求的理想方式，也是DevOps诞生的目的。
## 2、滚动发布
把一次完整的发布过程，合理地分成多个批次，每次发布一个批次，**成功后**，再发布下一个批次，最终完成所有批次的发布。在整个滚动过程期间，保证始终有可用的副本在运行，从而平滑的发布新版本，实现**零停机(without an outage)**、用户**零感知**，是一种非常主流的发布方式。由于其自动化程度比较高，通常需要复杂的发布工具支撑，而k8s可以完美的胜任这个任务。 
## 3、k8s滚动更新机制
**k8s创建副本应用程序的最佳方法就是部署(Deployment)，部署自动创建副本集(ReplicaSet)，副本集可以精确地控制每次替换的Pod数量，从而可以很好的实现滚动更新**。具体来说，k8s每次使用一个新的副本控制器(replication controller)来替换已存在的副本控制器，从而始终使用一个新的Pod模板来替换旧的pod模板。
>大致步骤如下：
>1. 创建一个新的replication controller。
>2. 增加或减少pod副本数量，直到满足当前批次期望的数量。
>3. 删除旧的replication controller。

## 4、演示
>使用kubectl更新一个已部署的应用程序，并模拟回滚。为了方便分析，将应用程序的pod副本数量设置为10。
``` bash
$ kubectl run busy --image=busybox:1.28.4 sleep 36000000 --replicas=10
```
### 4.1. 发布微服务
- 当前服务状态查看
``` bash
# 查看部署列表
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        10        10           10          5m

# 查看正在运行的pod
root@kube-aio:~# kubectl get pod | grep busy
busy-794c95f5d7-56b6w        1/1       Running   0          5m
busy-794c95f5d7-8ddjr        1/1       Running   0          5m
busy-794c95f5d7-8zm8r        1/1       Running   0          5m
busy-794c95f5d7-9hjhp        1/1       Running   0          5m
busy-794c95f5d7-df2r2        1/1       Running   0          5m
busy-794c95f5d7-fsn94        1/1       Running   0          5m
busy-794c95f5d7-k4w8r        1/1       Running   0          5m
busy-794c95f5d7-lsmgb        1/1       Running   0          5m
busy-794c95f5d7-rg8kw        1/1       Running   0          5m
busy-794c95f5d7-xpxxt        1/1       Running   0          5m

# 通过pod描述，查看应用程序的当前映像版本
root@kube-aio:~# kubectl describe pod busy-794c95f5d7-56b6w |grep Image
    Image:         busybox:1.28.4
    Image ID:      docker-pullable://busybox@sha256:141c253bc4c3fd0a201d32dc1f493bcf3fff003b6df416dea4f41046e0f37d47
```
- 升级镜像版本到1.29
  - 为了更清晰看到更新过程，可另开一个窗口使用`$ watch kubectl get deployment busy`实时查看变化
``` bash
$ kubectl set image deployments/busy busy=busybox:1.29
```
### 4.2. 验证发布
``` bash
# 检查rollout状态
root@kube-aio:~# kubectl rollout status deployments/busy
deployment "busy" successfully rolled out

# 检查pod详情
root@kube-aio:~# kubectl describe pod busy-665cdb7b-44jnt |grep Image
    Image:         busybox:1.29
    Image ID:      docker-pullable://busybox@sha256:cb63aa0641a885f54de20f61d152187419e8f6b159ed11a251a09d115fdff9bd
```
从上面可以看到，镜像已经升级到1.29版本
### 4.3. 回滚发布
``` bash
# 回滚发布
root@kube-aio:~# kubectl rollout undo deployments/busy
deployment.apps "busy" 

# 回滚完成
root@kube-aio:~# kubectl rollout status deployments/busy
deployment "busy" successfully rolled out

# 镜像又回退到1.28.4 版本
root@kube-aio:~# kubectl describe pod busy-794c95f5d7-4x9bn |grep Image
    Image:         busybox:1.28.4
    Image ID:      docker-pullable://busybox@sha256:141c253bc4c3fd0a201d32dc1f493bcf3fff003b6df416dea4f41046e0f37d47
```

到目前为止，整个滚动发布工作就圆满完成了！！！
**那么如果我们想回滚到指定版本呢？答案是k8s完美支持，并且还可以通过资源文件进行配置保留的历史版次量**。由于篇幅有限，感兴趣的朋友，可以自己下去实战，回滚命令如下：
```javascript
kubectl rollout undo deployment/busy --to-revision=<版次>
```
## 5、原理
k8s精确地控制着整个发布过程，分批次有序地进行着滚动更新，直到把所有旧的副本全部更新到新版本。实际上，k8s是通过两个参数来精确地控制着每次滚动的pod数量：

>* **`maxSurge` 滚动更新过程中运行操作期望副本数的最大pod数，可以为绝对数值(eg：5)，但不能为0；也可以为百分数(eg：10%)。**
>* **`maxUnavailable`  滚动更新过程中不可用的最大pod数，可以为绝对数值(eg：5)，但不能为0；也可以为百分数(eg：10%)。**

如果未指定这两个可选参数，则k8s会使用默认配置：  
``` bash
root@kube-aio:~# kubectl get deploy busy -o yaml
apiVersion: apps/v1 
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "3"
  creationTimestamp: 2018-08-19T02:42:56Z
  generation: 3
  labels:
    run: busy
  name: busy
  namespace: default
  resourceVersion: "199461"
  uid: 93fde307-a359-11e8-a93b-525400c61543
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      run: busy
  strategy:
    rollingUpdate:
      maxSurge: 1	# 滚动更新中最多超过预期值的 pod数
      maxUnavailable: 1	# 滚动更新中最多不可用的 pod数
    type: RollingUpdate
...
```
### 5.1. 浅析部署概况
``` bash
# 初始状态
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        10        10           10          1h

# 再做一遍回退
root@kube-aio:~# kubectl rollout undo deploy busy
deployment.apps "busy" 

# 更新过程1
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        11        2            9           1h

# 更新过程2
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        11        4            9           1h

# 更新过程3
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        11        6            9           1h

# 更新结束
root@kube-aio:~# kubectl get deploy busy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
busy      10        10        10           10          1h
```
>* `DESIRED`    最终期望处于READY状态的副本数   
>* `CURRENT`   当前的副本总数    
>* `UP-TO-DATE`   当前完成更新的副本数   
>* `AVAILABLE`   当前可用的副本数     

当前的副本总数：10(DESIRED) + 1(maxSurge) = 11，所以CURRENT为11。
当前可用的副本数：10(DESIRED) - 1(maxUnavailable) = 9，所以AVAILABLE为9。

### 5.2. 浅析部署详情

``` bash
root@kube-aio:~# kubectl describe deploy busy
Name:                   busy
Namespace:              default
CreationTimestamp:      Sun, 19 Aug 2018 12:27:19 +0800
Labels:                 run=busy
Annotations:            deployment.kubernetes.io/revision=2
Selector:               run=busy
Replicas:               10 desired | 10 updated | 10 total | 10 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  1 max unavailable, 1 max surge
Pod Template:
  Labels:  run=busy
  Containers:
   busy:
    Image:      busybox:1.29
    Port:       <none>
    Host Port:  <none>
    Args:
      sleep
      3600000
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   busy-84cb46955d (10/10 replicas created)
Events:
  Type    Reason             Age                 From                   Message
  ----    ------             ----                ----                   -------
  Normal  ScalingReplicaSet  1m                  deployment-controller  Scaled up replica set busy-9669c8599 to 10
  Normal  ScalingReplicaSet  46s                 deployment-controller  Scaled up replica set busy-84cb46955d to 1
  Normal  ScalingReplicaSet  46s                 deployment-controller  Scaled down replica set busy-9669c8599 to 9
  Normal  ScalingReplicaSet  46s                 deployment-controller  Scaled up replica set busy-84cb46955d to 2
  Normal  ScalingReplicaSet  43s                 deployment-controller  Scaled down replica set busy-9669c8599 to 8
  Normal  ScalingReplicaSet  43s                 deployment-controller  Scaled up replica set busy-84cb46955d to 3
  Normal  ScalingReplicaSet  43s                 deployment-controller  Scaled down replica set busy-9669c8599 to 7
  Normal  ScalingReplicaSet  43s                 deployment-controller  Scaled up replica set busy-84cb46955d to 4
  Normal  ScalingReplicaSet  40s                 deployment-controller  Scaled down replica set busy-9669c8599 to 6
  Normal  ScalingReplicaSet  28s (x12 over 40s)  deployment-controller  (combined from similar events): Scaled down replica set busy-9669c8599 to 0
```
整个滚动过程是通过控制两个副本集来完成的，新的副本集：busy-84cb46955d；旧的副本集：busy-9669c8599 。
理想状态下的滚动过程：
>1. 创建新副本集，并为其分配1个新版本的pod。
>2. 通知旧副本集，销毁1个旧版本的pod。
>3. 当旧副本销毁成功后，通知新副本集，再新增1个新版本的pod；当新副本创建成功后，通知旧副本再减少1个pod。
>只要销毁成功，新副本集就会创造新的pod，一直循环，直到旧的副本集pod数量为0。
### 5.4 总结
**`无论理想还是不理想，k8s最终都会使应用程序全部更新到期望状态，都会始终保持最大的副本总数和可用副本总数的不变性！！！`**

[阅读原文](http://www.cnblogs.com/justmine/p/8688828.html)

