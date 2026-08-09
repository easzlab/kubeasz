## headlamp

Kubernetes Dashboard 已归档，不再维护；本项目推荐使用 Headlamp 作为 Kubernetes Web UI。Headlamp 支持桌面客户端和集群内两种运行方式，本文介绍在集群内通过 Helm 安装的方法。

### 集成部署

假设已经使用kubeasz 部署k8s集群完成；headlamp 部署如下：（以单机集群为例，其他情况请修改集群名称'default'为实际的名称）

``` bash
# 1. 修改 clusters/default/config.yml 文件，设置 headlamp_install: "yes"

# 2. 下载需要的镜像
./ezdown -X headlamp

# 3. 执行安装，配置文件位于 clusters/default/yml/headlamp/ 目录
dk ezctl setup default 07
```

### 验证

``` bash
# 查看 Pod 运行状态
kubectl get pod -n kube-system -l app.kubernetes.io/name=headlamp

# 查看 Service
kubectl get svc -n kube-system headlamp
```

+ 增加`NodePort`方式暴露服务，这样集群外部可以使用 `http://NodeIP:NodePort` 直接访问 headlamp。

### 登陆

Headlamp 使用 Kubernetes RBAC 控制资源访问权限。桌面客户端会直接使用本机 kubeconfig；集群内访问可使用 ServiceAccount Token 或 OIDC 登录。

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

将输出的 Token 粘贴到 Headlamp 登录页。以上两个内置的token 仅用于演示；生产环境应根据实际角色创建最小权限的 RBAC 规则，避免长期使用高权限 Token。

### 参考

- [Headlamp 集群内安装文档](https://headlamp.dev/docs/latest/installation/in-cluster/)
- [Headlamp 认证文档](https://headlamp.dev/docs/latest/installation/)
- [Kubernetes Dashboard 归档仓库](https://github.com/kubernetes-retired/dashboard)
- [旧版 Dashboard 文档 1.6.3](dashboard.1.6.3.md)
- [旧版 Dashboard 文档 2.x](dashboard.2.x.md)
