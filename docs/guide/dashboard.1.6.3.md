## dashboard

本文档基于 dashboard 1.6.3版本，从 1.7.x 版本以后，dashboard 默认开启自带的登录验证界面，登录流程差异详见[新版本](dashboard.md)。

+ 注意：实际测试k8s版本<=1.9.1支持dashboard 1.6.3, 建议k8s 1.9 以后使用 dashboard 新版本。

### 部署

``` bash
# 部署dashboard 主yaml配置文件
$ kubectl create -f /etc/kubeasz/manifests/dashboard/1.6.3/kubernetes-dashboard.yaml
# 部署基本密码认证配置[可选]，密码文件位于 /etc/kubernetes/ssl/basic-auth.csv
$ kubectl create -f /etc/kubeasz/manifests/dashboard/1.6.3/ui-admin-rbac.yaml
$ kubectl create -f /etc/kubeasz/manifests/dashboard/1.6.3/ui-read-rbac.yaml
```

请在另外窗口打开 [kubernetes-dashboard.yaml](../../manifests/dashboard/1.6.3/kubernetes-dashboard.yaml)

+ 由于 kube-apiserver 启用了 RBAC授权，dashboard使用的 ServiceAccount `kubernetes-dashboard` 必须有相应的权限去访问apiserver(在新版本1.8.0中，该访问权限已按最小化方式授权)，在1.6.3 版本，先粗放一点，把`kubernetes-dashboard` 与 集群角色 `cluster-admin` 绑定，这样dashboard就拥有了所有访问apiserver的权限。
+ 开发测试环境为了方便配置dashboard-service时候，指定 `NodePort`方式暴露服务，这样集群外部可以使用 `http://NodeIP:NodePort` 方式直接访问 dashboard，生产环境建议关闭该访问途径。

### 验证

``` bash
# 查看pod 运行状态
kubectl get pod -n kube-system | grep dashboard
kubernetes-dashboard-86bd8778bf-w4974      1/1       Running   0          12h
# 查看dashboard service
kubectl get svc -n kube-system|grep dashboard
kubernetes-dashboard   NodePort    10.68.7.67      <none>        80:5452/TCP	12h
# 查看集群服务
kubectl cluster-info|grep dashboard
kubernetes-dashboard is running at https://192.168.1.10:6443/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy
# 查看pod 运行日志，关注有没有错误
kubectl logs kubernetes-dashboard-86bd8778bf-w4974 -n kube-system
```

### 访问

因为dashboard 作为k8s 原生UI，能够展示各种资源信息，甚至可以有修改、增加、删除权限，所以有必要对访问进行认证和控制，本项目预置部署的集群有以下安全设置：详见 [apiserver配置模板](../../roles/kube-master/templates/kube-apiserver.service.j2)

+ 启用 `TLS认证` `RBAC授权`等安全特性
+ 关闭 apiserver非安全端口8080的外部访问`--insecure-bind-address=127.0.0.1`
+ 关闭匿名认证`--anonymous-auth=false`
+ 补充启用基本密码认证 `--token-auth-file=/etc/kubernetes/ssl/basic-auth.csv`，[密码文件模板](../../roles/kube-master/templates/basic-auth.csv.j2)中按照每行(密码,用户名,序号)的格式，可以定义多个用户

#### 1. 临时访问：使用 `http://NodeIP:NodePort` 方式直接访问 dashboard，生产环境建议关闭该途径

#### 2. 用户+密码访问：安全性比证书方式差点，务必保管好密码文件`basic-auth.csv`

- 这里演示两种权限，使用admin 登录dashboard拥有所有权限，使用readonly 登录后仅查看权限，首先在 master节点文件 `/etc/kubernetes/ssl/basic-auth.csv` 确认用户名和密码，如果要增加或者修改用户，修改保存该文件后记得逐个重启你的master 节点
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
- 2.3 访问 `https://x.x.x.x:6443/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy` 使用 admin登录拥有所有权限，比如删除某个部署；使用 readonly登录只有查看权限，尝试删除某个部署会提示错误 `forbidden: User \"readonly\" cannot delete services/proxy in the namespace \"kube-system\"`

#### 3. 证书访问：最安全的方式，配置较复杂
- 使用集群CA 生成客户端证书，可以根据需要生成权限不同的证书，这里为了演示直接使用 kubectl使用的证书和key(在03.kubectl.yml阶段生成)，该证书拥有所有权限
- 指定格式导出该证书，进入`/etc/kubernetes/ssl`目录，使用命令`openssl pkcs12 -export -in admin.pem -inkey admin-key.pem -out kube-admin.p12` 提示输入证书密码和确认密码，可以用密码再增加一层保护，也可以直接回车跳过，完成后目录下多了 `kube-admin.p12`文件，将它分发给授权的用户
- 用户将 `kube-admin.p12` 双击导入证书即可，`IE` 和`Chrome` 中输入`https://x.x.x.x:6443/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy` 或者 `https://x.x.x.x:6443/ui` 即可访问。补充：最新firefox需要在浏览器中单独导入 [选项] - [隐私与安全] - [证书/查看证书] - [您的证书] 页面点击 [导入] 该证书

### 小结

+ dashboard 版本 1.6.3 访问控制实现较复杂，文档中给出的例子也有助于你理解 RBAC的灵活控制能力，当然最好去[官方文档](https://kubernetes.io/docs/admin/authorization/rbac/)学习一下，这块篇幅不长
+ 由于还未部署 Heapster 插件，当前 dashboard 不能展示 Pod、Nodes 的 CPU、内存等 metric 图形，后续部署 heapster后自然能够看到
+ 本文中的权限设置仅供演示用，生产环境请在此基础上修改成适合你安全需求的方式

