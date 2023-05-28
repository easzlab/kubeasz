# 使用kubeapps管理集群应用

Kubeapps 是一个基于 Web 的应用程序，它可以在 Kubernetes 集群上进行一站式安装，并使用户能够部署、管理和升级应用程序。

<img alt="kubeapps_dashboard" width="400" height="300" src="https://github.com/vmware-tanzu/kubeapps/raw/main/site/content/docs/latest/img/dashboard-login.png">

项目地址：https://github.com/vmware-tanzu/kubeapps
部署项目地址：https://github.com/bitnami/charts/tree/main/bitnami/kubeapps

## 使用kubeasz部署

- 1.编辑集群配置文件：clusters/${集群名}/config.yml

``` bash
kubeapps_install: "yes"                    # 启用安装
kubeapps_install_namespace: "kubeapps"     # 设置安装命名空间
kubeapps_working_namespace: "default"      # 设置默认应用命名空间
kubeapps_storage_class: "local-path"       # 设置存储storageclass，默认使用local-path-provisioner
kubeapps_chart_ver: "12.4.3"
```

- 2.下载相关容器镜像

``` bash
# 下载kubeapps镜像
/etc/kubeasz/ezdown -X kubeapps

# 下载local-path-provisioner镜像
/etc/kubeasz/ezdown -X local-path-provisioner
```

- 3.安装cluster-addon

``` bash
$ dk ezctl setup ${集群名} 07

# 执行成功后验证
$ kubectl get pod --all-namespaces |grep kubeapps
```

## 验证使用kubeapps

阅读文档：https://github.com/vmware-tanzu/kubeapps/blob/main/site/content/docs/latest/tutorials/getting-started.md

正式使用建议配置OAuth2/OIDC用户认证，这里仅验证使用k8s ServiceAccount 方式登陆，项目已预装三个用户权限：

- 1.kubeapps-admin-token，全局cluster-admin权限，不建议使用
- 2.kubeapps-edit-token，某个命名空间下的应用可写权限
- 3.kubeapps-view-token，某个命名空间下的应用只读权限

``` bash
# 获取UI访问地址，默认使用NodePort
kubectl get svc -n kubeapps kubeapps
NAME       TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
kubeapps   NodePort   10.68.92.88   <none>        80:32490/TCP   117m

# 获取admin token
kubectl get secrets -n kube-system kubeapps-admin-token -o go-template='{{.data.token | base64decode}}'

# 获取某命名空间应用部署权限 token
kubectl get secrets -n default kubeapps-edit-token -o go-template='{{.data.token | base64decode}}'

# 获取某命名空间应用部署权限 token
kubectl get secrets -n default kubeapps-view-token -o go-template='{{.data.token | base64decode}}'
```

打开浏览器，访问http://${Node_IP}:32490，输入上面合适权限的token即可。
