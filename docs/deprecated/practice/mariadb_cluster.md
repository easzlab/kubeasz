# Mariadb 数据库集群

Mariadb 是从 MySQL 衍生出来的开源关系型数据库，目前兼容 mysql 5.7 版本；它也非常流行，拥有 Google Facebook 等重要企业用户。本文档介绍使用 helm charts 方式安装 mariadb cluster，仅供实践交流使用。

## 前提条件

- 已部署 k8s 集群，参考[这里](../setup/quickStart.md)
- 已部署 helm，参考[这里](../guide/helm.md)
- 集群提供持久性存储，参考[这里](../setup/08-cluster-storage.md)


## mariadb charts 配置修改

按照惯例，直接把 chart 下载到本地，然后把配置复制 values.yaml 出来进行修改，这样方便以后整体更新 chart，安装实际使用需要修改配置文件

``` bash
$ cd /etc/kubeasz/manifests/mariadb-cluster
# 编辑 my-values.yaml 修改以下部分

service:
  type: NodePort     # 方便集群外部访问
  port: 3306
  nodePort:
    master: 33306    # 设置主库的nodePort
    slave: 33307     # 设置从库的nodePort

rootUser:            # 设置 root 密码
  password: test.c0m
  forcePassword: true

db:                  # 设置初始测试数据库
  user: hello
  password: hello
  name: hello
  forcePassword: true

replication:         # 设置主从复制
  enabled: true
  user: replicator
  password: R4%forep11CAT0r
  forcePassword: true

master:
  affinity: {}
  antiAffinity: soft
  tolerations: []
  persistence:
    enabled: true    # 启用持久化存储
    mountPath: /bitnami/mariadb
    storageClass: "nfs-db"  # 设置使用 nfs-db 存储类
    annotations: {}
    accessModes:
    - ReadWriteOnce
    size: 5Gi        # 设置存储容量 

slave:
  replicas: 1
  affinity: {}
  antiAffinity: soft
  tolerations: []
  persistence:
    enabled: false   # 从库这里没有启用持久性存储
```

## 安装

使用 helm 安装

``` bash
$ cd /etc/kubeasz/manifests/mariadb-cluster
$ helm install --name mariadb --namespace default -f my-values.yaml ./mariadb
```

## 验证

``` bash
$ kubectl get pod,svc | grep mariadb
pod/mariadb-mariadb-master-0      1/1     Running   0          27m
pod/mariadb-mariadb-slave-0       1/1     Running   0          29m

service/mariadb                       NodePort    10.68.170.168   <none>        3306:33306/TCP       29m
service/mariadb-mariadb-slave         NodePort    10.68.151.95    <none>        3306:33307/TCP       29m
```

