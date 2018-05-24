# nfs 动态存储

## 前言
在kubernetes(k8s)中对于存储的资源抽象了两个概念，分别是PersistentVolume(PV)、PersistentVolumeClaim(PVC)。
- PV是集群中的资源
- PVC是对这些资源的请求。

如上面所说PV和PVC都只是抽象的概念，在k8s中是通过插件的方式提供具体的存储实现。目前包含有NFS、iSCSI和云提供商指定的存储系统，更多的存储实现[参考官方文档](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)。

这里PV又有两种提供方式: 静态或者动态。
这篇文档主要就是介绍**NFS动态存储**的实现。

## NFS服务器
首先我们需要一个NFS服务器，用于提供底层存储。通过文档[nfs-server](nfs-server.md)，我们可以创建一个NFS服务器。

## 创建nfs-client
切换到项目`manifests/nfs-provisioner`目录，
修改`nfs-client-provisioner.yaml`文件，将`nfs server`地址和`共享目录`调整为我们自己NFS服务器的配置, 在示例中分别就是`10.1.241.230`和`/home/share/k8s-pv`.

调整完后，执行以下命令进行创建：
`kubectl create -f nfs-client-provisioner.yaml`

若出现类似错误：`Unable to mount volumes for pod`
则可能是NFS服务器的配置不正确，需要调整为你自己的服务器及共享目录。


## 测试
查看运行情况，执行`kubectl get pod --namespace=kube-system`，当出现类似以下信息时，则表示nfs-client运行正常了:
```
nfs-client-provisioner-667bbdcc94-7vdl8       1/1       Running   0          35s
```

我们这里最后再测试一下动态存储。
为了便于学习理解，这里分别将StorageClass、PVC、Pod放在了3个文件中。

### 创建StorageClass
我们首先创建`StorageClass`，用于引用`nfs-client`中的提供者。  

`kubectl create -f nfs-dynamic-storageclass.yaml`

### 创建PVC
根据`StorageClass`，我们再创建PVC。在静态提供方式中，PVC是根据PV来进行绑定的，这里我们是根据`StorageClass`动态方式进行绑定。

`kubectl create -f test/test-claim.yaml`  

创建完PVC后，我们可以看看是否有绑定成功：
```
# kubectl get pvc
NAME         STATUS    VOLUME
test-claim   Bound     pvc-a877172b-5f49-11e8-b675-d8cb8ae6325a
```
当`STATUS`字段出现**Bound**时，就表明PVC已经绑定成功了。

### 创建测试Pod
最后我们再来创建Pod，引用我们刚刚创建的PVC。

`kubectl create -f test/test-pod.yaml`

这个Pod很简单，就是启动完成后，在挂载的目录中创建一个`SUCCESS`文件。
启动完成后，我们可以到NFS服务器去看下：
```
.
└── default-test-claim-pvc-a877172b-5f49-11e8-b675-d8cb8ae6325a
    └── SUCCESS
```
如上，可以发现挂载的时候，nfs-client根据PVC自动创建了一个目录，我们Pod中挂载的`/mnt`，实际引用的就是该目录，而我们在`/mnt`下创建的`SUCCESS`文件，也自动写入到了这里。

经过了简单的测试，我们可以了解动态存储的使用。
这里`StorageClass`并不是每次都要进行创建的，只需要创建一次就好，后面其他的PVC可以重复使用，而不是像静态PV一样，只能被一个PVC绑定。

# 后续
后面当我们需要为上层应用提供持久化存储时，只需要提供`StorageClass`即可。
很多应用都会根据`StorageClass`来创建他们的所需的PVC, 最后再把PVC挂载到他们的Deployment或StatefulSet中使用，比如：efk、jenkins等
