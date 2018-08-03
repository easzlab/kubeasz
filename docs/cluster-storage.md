# K8S 集群存储 

## 前言
在kubernetes(k8s)中对于存储的资源抽象了两个概念，分别是PersistentVolume(PV)、PersistentVolumeClaim(PVC)。
- PV是集群中的资源
- PVC是对这些资源的请求。

如上面所说PV和PVC都只是抽象的概念，在k8s中是通过插件的方式提供具体的存储实现。目前包含有NFS、iSCSI和云提供商指定的存储系统，更多的存储实现[参考官方文档](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)。

这里PV又有两种提供方式: 静态或者动态。
本篇以介绍 **NFS存储** 为例，讲解k8s 众多存储方案中的一个实现。

## 静态 PV
首先我们需要一个NFS服务器，用于提供底层存储。通过文档[nfs-server](nfs-server.md)，我们可以创建一个NFS服务器。

- 创建静态 pv，指定容量，访问模式，回收策略，存储类等；参考[这里](https://github.com/feiskyer/kubernetes-handbook/blob/master/zh/concepts/persistent-volume.md)

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

## 创建动态PV

在一个工作k8s 集群中，`PVC`请求会很多，如果每次都需要管理员手动去创建对应的 `PV`资源，那就很不方便；因此 K8S还提供了多种 `provisioner`来动态创建 `PV`，不仅节省了管理员的时间，还可以根据`StorageClasses`封装不同类型的存储供 PVC 选用。

项目中的 `role: cluster-storage`目前支持自建nfs 和aliyun_nas 的动态`provisioner`

- 1.更新项目源码，生成自定义配置文件（该配置文件被.gitignore忽略）

``` bash
$ ansible-playbook /etc/ansible/tools/init_vars.yml
```
- 2.编辑自定义配置文件：上述命令执行后生成的roles/cluster-storage/vars/main.yml

``` bash
# 比如创建nfs provisioner
storage_nfs_enabled: "yes"
nfs_server: "192.168.1.8"
nfs_server_path: "/data/nfs"
nfs_storage_class: "class-nfs-01"
nfs_provisioner_name: "nfs-provisioner-01"
```
- 3.创建 nfs provisioner

``` bash
$ ansible-playbook /etc/ansible/roles/cluster-storage/cluster-storage.yml
# 执行成功后验证
$ kubectl get pod --all-namespaces |grep nfs-prov
kube-system   nfs-provisioner-01-6b7fbbf9d4-bh8lh        1/1       Running   0          1d
```
**注意** k8s集群可以使用多个nfs provisioner，重复上述步骤2 修改使用不同的`nfs server` `nfs_storage_class` `nfs_provisioner_name`后执行步骤3创建即可。

## 验证使用动态 PV

切换到项目`manifests/storage`目录，编辑`test.yaml`文件，根据前文配置情况修改`storageClassName`即可；然后执行以下命令进行创建：

``` bash
$ kubectl apply -f test.yaml

# 验证测试pod
$ kubectl get pod --all-namespaces |grep test
default       test                                       1/1       Running   0          1m

# 验证自动创建的pv 资源，
$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                STORAGECLASS           REASON    AGE
pvc-8f1b4ced-92d2-11e8-a41f-5254008ec7c0   1Mi        RWX            Delete           Bound     default/test-claim   nfs-dynamic-class-01             3m

# 验证PVC已经绑定成功：STATUS字段为 Bound
$ kubectl get pvc
NAME         STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS           AGE
test-claim   Bound     pvc-8f1b4ced-92d2-11e8-a41f-5254008ec7c0   1Mi        RWX            nfs-dynamic-class-01   3m
```

另外，Pod启动完成后，在挂载的目录中创建一个`SUCCESS`文件。我们可以到NFS服务器去看下：

```
.
└── default-test-claim-pvc-a877172b-5f49-11e8-b675-d8cb8ae6325a
    └── SUCCESS
```
如上，可以发现挂载的时候，nfs-client根据PVC自动创建了一个目录，我们Pod中挂载的`/mnt`，实际引用的就是该目录，而我们在`/mnt`下创建的`SUCCESS`文件，也自动写入到了这里。

# 后续
后面当我们需要为上层应用提供持久化存储时，只需要提供`StorageClass`即可。很多应用都会根据`StorageClass`来创建他们的所需的PVC, 最后再把PVC挂载到他们的Deployment或StatefulSet中使用，比如：efk、jenkins等
