## dashboard

本文档基于 dashboard 1.8.3版本，k8s版本 1.9.x。旧版文档[dashboard-1.6.3](dashboard.1.6.3.md)

### 部署

如果之前已按照本项目部署dashboard1.6.3，先删除旧版本：`kubectl delete -f /etc/ansible/manifests/dashboard/1.6.3/`

1.8.3配置文件参考[官方文档](https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml)

+ 增加了通过`api-server`方式访问dashboard
+ 增加了`NodePort`方式暴露服务，这样集群外部可以使用 `https://NodeIP:NodePort` (注意是https不是http，区别于1.6.3版本) 直接访问 dashboard，生产环境建议关闭该访问途径。

安装部署

``` bash
# 部署dashboard 主yaml配置文件
$ kubectl create -f /etc/ansible/manifests/dashboard/kubernetes-dashboard.yaml
# 部署基本密码认证配置[可选]，密码文件位于 /etc/kubernetes/ssl/basic-auth.csv
$ kubectl create -f /etc/ansible/manifests/dashboard/ui-admin-rbac.yaml
$ kubectl create -f /etc/ansible/manifests/dashboard/ui-read-rbac.yaml
```

### 验证

``` bash
# 查看pod 运行状态
kubectl get pod -n kube-system | grep dashboard
kubernetes-dashboard-7c74685c48-9qdpn   1/1       Running   0          22s
# 查看dashboard service
kubectl get svc -n kube-system|grep dashboard
kubernetes-dashboard   NodePort    10.68.219.38   <none>        443:24108/TCP                   53s
# 查看集群服务
kubectl cluster-info|grep dashboard
kubernetes-dashboard is running at https://192.168.1.1:6443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy
# 查看pod 运行日志，关注有没有错误
kubectl logs kubernetes-dashboard-7c74685c48-9qdpn -n kube-system
```

### 访问

因为dashboard 作为k8s 原生UI，能够展示各种资源信息，甚至可以有修改、增加、删除权限，所以有必要对访问进行认证和控制，本项目部署的集群有以下安全设置：详见 [apiserver配置模板](../../roles/kube-master/templates/kube-apiserver.service.j2)

+ 启用 `TLS认证` `RBAC授权`等安全特性
+ 关闭 apiserver非安全端口8080的外部访问`--insecure-bind-address=127.0.0.1`
+ 关闭匿名认证`--anonymous-auth=false`
+ 补充启用基本密码认证 `--basic-auth-file=/etc/kubernetes/ssl/basic-auth.csv`，[密码文件模板](../../roles/kube-master/templates/basic-auth.csv.j2)中按照每行(密码,用户名,序号)的格式，可以定义多个用户

新版本dashboard登陆可以分为两步，类似流行的双因子登陆系统：

+ 第一步通过api-server本身安全认证流程，与之前1.6.3版本相同
+ 第二步通过dashboard自带的登陆流程，使用`Kubeconfig` `Token`等方式登陆

#### 1. 临时访问：使用 `https://NodeIP:NodePort` 方式直接访问 dashboard，生产环境建议关闭该途径
打开页面出现dashboard 新版本自带的登陆页面。Kubernetes仪表盘支持两种登录方式：Kubeconfig、令牌

- 令牌登录
选择“令牌(Token)”方式登陆，关于令牌的获取[参考](https://github.com/kubernetes/dashboard/wiki/Creating-sample-user)

``` bash
# 创建Service Account 和 ClusterRoleBinding
$ kubectl create -f /etc/ansible/manifests/dashboard/admin-user-sa-rbac.yaml
# 获取 Bearer Token，找到输出中 ‘token:’ 开头那一行
$ kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}')
```

- Kubeconfig登录
Kubeconfig文件默认位置：`/root/.kube/config`，该文件中默认没有token字段，使用Kubeconfig方式登录，还需要将token追加到该文件中，完整的文件格式如下：
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

#### 2. 用户+密码访问：安全性比证书方式差点，务必保管好密码文件`basic-auth.csv`

- 这里演示两种权限，使用admin 登陆dashboard拥有所有权限，使用readonly 登陆后仅查看权限，首先在 master节点文件 `/etc/kubernetes/ssl/basic-auth.csv` 确认用户名和密码，如果要增加或者修改用户，修改保存该文件后记得逐个重启你的master 节点
- 为了演示用户密码访问，如果你已经完成证书访问方式，你可以在浏览器删除证书，或者访问时候浏览器询问你证书时不选证书
- 2.1 设置用户admin 的RBAC 权限，如下运行配置文件 `kubectl create -f ui-admin-rbac.yaml`

``` bash
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ui-admin
rules:
- apiGroups:
  - ""
  resources:
  - services
  - services/proxy
  verbs:
  - '*'

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ui-admin-binding
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ui-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: admin
```  
- 2.2 设置用户readonly 的RBAC 权限，如下运行配置文件 `kubectl create -f ui-read-rbac.yaml`

``` bash
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ui-read
rules:
- apiGroups:
  - ""
  resources:
  - services
  - services/proxy
  verbs:
  - get
  - list
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ui-read-binding
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ui-read
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: readonly
```
- 2.3 访问 `https://x.x.x.x:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy` (该URL具体使用`kubectl cluster-info`查看) 使用 admin登陆拥有所有权限，比如删除某个部署；使用 readonly登陆只有查看权限，尝试删除某个部署会提示错误 `forbidden: User \"readonly\" cannot delete services/proxy in the namespace \"kube-system\"`

- dashboard自带的登陆流程同上

#### 3. 证书访问：最安全的方式，配置较复杂
- 使用集群CA 生成客户端证书，可以根据需要生成权限不同的证书，这里为了演示直接使用 kubectl使用的证书和key(在03.kubectl.yml阶段生成)，该证书拥有所有权限
- 指定格式导出该证书，进入`/etc/kubernetes/ssl`目录，使用命令`openssl pkcs12 -export -in admin.pem -inkey admin-key.pem -out kube-admin.p12` 提示输入证书密码和确认密码，可以用密码再增加一层保护，也可以直接回车跳过，完成后目录下多了 `kube-admin.p12`文件，将它分发给授权的用户
- 用户将 `kube-admin.p12` 双击导入证书即可，`IE` 和`Chrome` 中输入`https://x.x.x.x:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy`(该URL具体使用`kubectl cluster-info`查看) 即可访问。补充：最新firefox需要在浏览器中单独导入 [选项] - [隐私与安全] - [证书/查看证书] - [您的证书] 页面点击 [导入] 该证书
- dashboard自带的登陆流程同上

#### 4. 授予admin权限，跳过登录
**注意：** 首先需要确保你知道这样做的后果，授予admin权限后安全性较低，不建议在生产环境中使用。

- 创建admin角色
```
$ kubectl create -f /etc/ansible/manifests/dashboard/admin-user-sa-rbac.yaml
```

- 修改dashboard角色配置
编辑`/etc/ansible/manifests/dashboard/kubernetes-dashboard.yaml`文件

找到以下配置：
```
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kubernetes-dashboard-minimal
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kubernetes-dashboard-minimal
subjects:
- kind: ServiceAccount
  name: kubernetes-dashboard
  namespace: kube-system
```

修改为：
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubernetes-dashboard-admin
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: kubernetes-dashboard
  namespace: kube-system
```

- 最后再创建dashboard
`# kubectl create -f /etc/ansible/manifests/dashboard/kubernetes-dashboard.yaml`

访问dashboard：
`https://x.x.x.x:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy`(该URL具体使用`kubectl cluster-info`查看) ，直接点击跳过按钮即可


### 小结

+ dashboard 访问控制实现较复杂，文档中给出的例子也有助于你理解 RBAC的灵活控制能力，当然最好去[官方文档](https://kubernetes.io/docs/admin/authorization/rbac/)学习一下，这块篇幅不长
+ 由于还未部署 Heapster 插件，当前 dashboard 不能展示 Pod、Nodes 的 CPU、内存等 metric 图形，后续部署 heapster后自然能够看到
+ 本文中的权限设置仅供演示用，生产环境请在此基础上修改成适合你安全需求的方式

