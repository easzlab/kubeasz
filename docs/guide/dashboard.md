## dashboard

本文档基于 dashboard 7.12.0 版本，k8s 1.32 版本，dashboard 7.0.0 以后引入大量不兼容变化。

### 部署

参考 https://github.com/kubernetes/dashboard

新版dashboard使用helm chart部署如下：（以单机集群为例，其他情况请修改集群名称'default'为实际的名称）

``` bash
# 1. 修改 clusters/default/config.yml 文件，设置 dashboard_install: "yes"

# 2. 下载dashboard 需要的镜像
./ezdown -X dashboard

# 3. 执行安装，配置文件位于 clusters/default/yml/dashboard/ 目录
./ezctl setup default 07
```

+ 增加`NodePort`方式暴露服务，这样集群外部可以使用 `https://NodeIP:NodePort` (注意是https不是http) 直接访问 dashboard。

### 验证

``` bash
# 查看pod 运行状态
kubectl get pod -n kube-system |grep kubernetes-dashboard
kubernetes-dashboard-api-6d77cb7964-4tklq               1/1     Running   0          17h
kubernetes-dashboard-auth-5fbd64f659-f9dst              1/1     Running   0          17h
kubernetes-dashboard-kong-6dcdbf5dfd-829h4              1/1     Running   0          17h
kubernetes-dashboard-metrics-scraper-7757c48476-4lcrq   1/1     Running   0          17h
kubernetes-dashboard-web-5f9f47979-7khrk                1/1     Running   0          17h

# 查看service
kubectl get svc -n kube-system |grep kong
kubernetes-dashboard-kong-proxy        NodePort    10.68.148.170   <none>   443:31544/TCP  17h
```

### 登陆

因为dashboard 作为k8s 原生UI，能够展示各种资源信息，甚至可以有修改、增加、删除权限，所以有必要对访问进行认证和控制，为演示方便这里使用 `https://NodeIP:NodePort` 方式访问 dashboard，目前支持登录方式：令牌(Token)

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

### 参考

- [旧版文档 dashboard 1.6.3](dashboard.1.6.3.md)
- [旧版文档 dashboard 2.x](dashboard.2.x.md)
