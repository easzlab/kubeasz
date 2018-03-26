## 监控数据持久化
部署[heapster](heapster.md)时我们知道监控数据是持久化存储到`influxdb`中的，但是在influxdb.yaml文件中存储使用的是`emptyDir`类型，所以当influxdb POD被删除时，监控数据也就丢失了。本文档讲解一个使用nfs持久化保存监控数据的例子。

### 前提
环境准备一个nfs服务器，如果没有可以参考[nfs-server](nfs-server.md)创建。

### 创建 PV
PersistentVolume (PV) 和 PersistentVolumeClaim (PVC) 提供了方便的持久化卷；`PV` 是集群的存储资源，就像 `node`是集群的计算资源，PV 可以静态或动态创建，这里使用静态方式创建；`PVC` 就是用来申请`PV` 资源，它可以直接挂载在`POD` 里面使用。更多知识请访问[官网](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)。根据你监控日志的多少和需保存时间需求创建固定大小的`PV` 资源，例子：

``` bash
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-influxdb
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  volumeMode: Filesystem
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: slow
  nfs:
    # 根据实际共享目录修改
    path: /share
    # 根据实际 nfs服务器地址修改
    server: 192.168.1.208
```

### 修改influxdb 存储卷

使用PVC 替换 volumes `emptyDir:{}`，创建PVC 如下：

``` bash
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: influxdb-claim
  namespace: kube-system
spec:
  accessModes:
    - ReadWriteMany
  volumeMode: Filesystem
  resources:
    requests:
      storage: 3Gi
  storageClassName: slow
``` 
+ 注意`PV` 是不区分namespace，而`PVC` 是区分namespace的

### 安装持久化监控

如果之前已经安装本项目创建了`heapster`，请使用如下删除 `influxdb POD`:

``` bash
kubectl delete -f /etc/ansible/manifests/heapster/influxdb.yaml
```

然后使用如下命令新建持久化的 `influxdb POD` :

``` bash
kubectl create -f /etc/ansible/manifests/heapster/influxdb-with-pv/
```

### 验证监控数据的持久性

手动删除 `influxdb` POD，半小时后再次创建，登陆grafana 确认历史数据是否还在。

``` bash
# 删除 influxdb deploy
kubectl delete -f /etc/ansible/manifests/heapster/influxdb-with-pv/influxdb.yaml

# 等待半小时后重新创建
kubectl create -f /etc/ansible/manifests/heapster/influxdb-with-pv/influxdb.yaml
```
+ 如果同时删除了 influxdb-pvc，那么根据策略`Recycle`监控历史数据也就删除了。
