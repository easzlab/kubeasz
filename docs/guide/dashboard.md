## dashboard

本文档基于 dashboard 2.2 版本，k8s 1.22 版本，因 dashboard 1.7 以后默认开启了自带的登录验证机制，因此不同版本登录有差异：

- 旧版（<= 1.6）建议通过apiserver访问，直接通过apiserver 认证授权机制去控制 dashboard权限，详见[旧版文档](dashboard.1.6.3.md)
- 新版（>= 1.7）可以使用自带的登录界面，使用不同Service Account Tokens 去控制访问 dashboard的权限

### 部署

参考 https://github.com/kubernetes/dashboard

+ 增加了通过`api-server`方式访问dashboard
+ 增加了`NodePort`方式暴露服务，这样集群外部可以使用 `https://NodeIP:NodePort` (注意是https不是http，区别于1.6.3版本) 直接访问 dashboard。

安装部署

``` bash
# ezctl 集成部署组件，xxxx 代表集群部署名
# dashboard 部署文件位于 /etc/kubeasz/clusters/xxxx/yml/dashboard/ 目录
./ezctl setup xxxx 07
```

### 验证部署

``` bash
# 查看pod 运行状态
kubectl get pod -n kube-system | grep dashboard
dashboard-metrics-scraper-856586f554-l6bf4   1/1     Running   0          35m
kubernetes-dashboard-698d4c759b-67gzg        1/1     Running   0          35m

# 查看dashboard service
kubectl get svc -n kube-system|grep dashboard
kubernetes-dashboard   NodePort    10.68.219.38   <none>        443:24108/TCP                   53s

# 查看pod 运行日志
kubectl logs -n kube-system kubernetes-dashboard-698d4c759b-67gzg
```

### 登陆

因为dashboard 作为k8s 原生UI，能够展示各种资源信息，甚至可以有修改、增加、删除权限，所以有必要对访问进行认证和控制，为演示方便这里使用 `https://NodeIP:NodePort` 方式访问 dashboard，支持两种登录方式：Kubeconfig、令牌(Token)

**注意：** 使用chrome浏览器访问 `https://NodeIP:NodePort` 可能提示安全风险无法访问，可以换firefox浏览器设置安全例外，继续访问。 

- Token令牌方式登录（admin）

选择 Token 方式登录，复制下面输出的admin token 字段到输入框

``` bash
# 获取 Bearer Token，找到输出中 ‘token:’ 开头的后面部分
$ kubectl describe -n kube-system secrets admin-user 
```

- Token令牌方式登录（只读）

选择 Token 方式登录，复制下面输出的read token 字段到输入框

``` bash
# 获取 Bearer Token，找到输出中 ‘token:’ 开头的后面部分
$ kubectl describe -n kube-system secrets dashboard-read-user 
```

- Kubeconfig登录（admin）
Admin kubeconfig文件默认位置：`/root/.kube/config`，该文件中默认没有token字段，使用Kubeconfig方式登录，还需要将token追加到该文件中，完整的文件格式如下：
```
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdxxxxxxxxxxxxxx
    server: https://192.168.1.2:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: admin
  name: kubernetes
current-context: kubernetes
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRxxxxxxxxxxx
    client-key-data: LS0tLS1CRUdJTxxxxxxxxxxxxxx
    token: eyJhbGcixxxxxxxxxxxxxxxx
```

- Kubeconfig登录（只读）
首先[创建只读权限 kubeconfig文件](../op/kcfg-adm.md)，然后类似追加只读token到该文件，略。

### 参考

- 1.[Dashboard docs](https://github.com/kubernetes/dashboard/blob/master/docs/README.md)
- 2.[a-read-only-kubernetes-dashboard](https://blog.cowger.us/2018/07/03/a-read-only-kubernetes-dashboard.html)
