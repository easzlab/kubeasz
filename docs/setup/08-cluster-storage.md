# K8S 集群存储 

## 前言
在kubernetes(k8s)中对于存储的资源抽象了两个概念，分别是PersistentVolume(PV)、PersistentVolumeClaim(PVC)。
- PV是集群中的资源
- PVC是对这些资源的请求。

如上面所说PV和PVC都只是抽象的概念，在k8s中是通过插件的方式提供具体的存储实现。目前包含有NFS、iSCSI和云提供商指定的存储系统，更多的存储实现[参考官方文档](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)。

以下介绍两种`provisioner`, 可以提供静态或者动态的PV

- [nfs-provisioner](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner): NFS存储目录供应者
- [local-path-provisioner](https://github.com/rancher/local-path-provisioner): 本地存储目录供应者

## NFS存储目录供应者

首先我们需要一个NFS服务器，用于提供底层存储。通过文档[nfs-server](../guide/nfs-server.md)，我们可以创建一个NFS服务器。

### 静态 PV
- 创建静态 pv，指定容量，访问模式，回收策略，存储类等

``` bash
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-es-0
spec:
  capacity:
    storage: 4Gi
  accessModes:
    - ReadWriteMany
  volumeMode: Filesystem
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: "es-storage-class"
  nfs:
    # 根据实际共享目录修改
    path: /share/es0
    # 根据实际 nfs服务器地址修改
    server: 192.168.1.208
```
- 创建 pvc即可绑定使用上述 pv了，具体请看后文 test pod例子

### 创建动态PV

在一个工作k8s 集群中，`PVC`请求会很多，如果每次都需要管理员手动去创建对应的 `PV`资源，那就很不方便；因此 K8S还提供了多种 `provisioner`来动态创建 `PV`，不仅节省了管理员的时间，还可以根据`StorageClasses`封装不同类型的存储供 PVC 选用。

项目中以nfs-client-provisioner为例 https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner

- 1.编辑集群配置文件：clusters/${集群名}/config.yml

``` bash
... 省略
# 在role:cluster-addon 中启用nfs-provisioner 安装
nfs_provisioner_install: "yes"			# 修改为yes
nfs_provisioner_namespace: "kube-system"
nfs_provisioner_ver: "v4.0.1"
nfs_storage_class: "managed-nfs-storage"	
nfs_server: "192.168.31.244"			# 修改为实际nfs server地址
nfs_path: "/data/nfs"				# 修改为实际的nfs共享目录

```

- 2.创建 nfs provisioner

``` bash
$ dk ezctl setup ${集群名} 07 

# 执行成功后验证
$ kubectl get pod --all-namespaces |grep nfs-client
kube-system   nfs-client-provisioner-84ff87c669-ksw95      1/1     Running     0          21m
```

- 3.验证使用动态 PV

在目录clusters/${集群名}/yml/nfs-provisioner/ 有个测试例子

``` bash
$ kubectl apply -f /etc/kubeasz/clusters/hello/yml/nfs-provisioner/test-pod.yaml

# 验证测试pod
kubectl get pod
NAME       READY   STATUS      RESTARTS   AGE
test-pod   0/1     Completed   0          6h36m

# 验证自动创建的pv 资源，
kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                STORAGECLASS          REASON   AGE
pvc-44d34a50-e00b-4f6c-8005-40f5cc54af18   2Mi        RWX            Delete           Bound    default/test-claim   managed-nfs-storage            6h36m

# 验证PVC已经绑定成功：STATUS字段为 Bound
kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
test-claim   Bound    pvc-44d34a50-e00b-4f6c-8005-40f5cc54af18   2Mi        RWX            managed-nfs-storage   6h37m
```

另外，Pod启动完成后，在挂载的目录中创建一个`SUCCESS`文件。我们可以到NFS服务器去看下：

```
.
└── default-test-claim-pvc-44d34a50-e00b-4f6c-8005-40f5cc54af18
    └── SUCCESS
```
如上，可以发现挂载的时候，nfs-client根据PVC自动创建了一个目录，我们Pod中挂载的`/mnt`，实际引用的就是该目录，而我们在`/mnt`下创建的`SUCCESS`文件，也自动写入到了这里。

后面当我们需要为上层应用提供持久化存储时，只需要提供`StorageClass`即可。很多应用都会根据`StorageClass`来创建他们的所需的PVC, 最后再把PVC挂载到他们的Deployment或StatefulSet中使用，比如：efk、jenkins等

## 本地存储目录供应者

当应用对于磁盘I/O性能要求高，比较适合本地文件目录存储，特别地可以本地挂载SSD磁盘（注意本地磁盘需要配置raid冗余策略）。Local Path Provisioner 可以方便地在k8s集群中使用本地文件目录存储。

在kubeasz项目中集成安装

- 1.编辑集群配置文件：clusters/${集群名}/config.yml

``` bash
... 省略
local_path_provisioner_install: "yes" # 修改为yes
# 设置默认本地存储路径
local_path_provisioner_dir: "/opt/local-path-provisioner"
```

- 2.创建 local path provisioner

``` bash
$ dk ezctl setup ${集群名} 07

# 执行成功后验证
$ kubectl get pod --all-namespaces |grep nfs-client-provisioner
```

- 3.验证使用（略）
