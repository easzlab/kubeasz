## kubeasz 3.6.2

kubeasz 3.6.2 发布：支持k8s v1.28版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.28.1
- etcd: v3.5.9
- containerd: 1.6.23
- runc: v1.1.9
- cni: v1.3.0
- coredns: 1.11.1
- cilium: 1.13.6
- flannel: v0.22.2

### 修改kubeasz支持k8s版本对应规则 

原有模式每个k8s大版本都有推荐对应的kubeasz版本，这样做会导致kubeasz版本碎片化，追踪问题很麻烦，而且也影响普通用户安装体验。从kubeasz 3.6.2版本开始，默认最新版本kubeasz兼容支持安装最新的三个k8s大版本。具体安装说明如下：

(如果/etc/kubeasz/bin 目录下已经有kube* 文件，需要先删除 rm -f /etc/kubeasz/bin/kube*)

- 安装 k8s v1.28: 使用 kubeasz 3.6.2，执行./ezdown -D 默认下载即可
- 安装 k8s v1.27: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.27.5 下载
- 安装 k8s v1.26: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.26.8 下载
- 安装 k8s v1.25: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.25.13 下载
- 安装 k8s v1.24: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.24.17 下载


### 重要更新

- 增加支持containerd 可配置trusted insecure registries 
- 修复calico rr 模式的节点设置 #1308
- 修复自定义节点名称设置 /etc/hosts方案
- fix: kubelet failed when enabling kubeReserved or systemReserved

### 其他

- 修复：disable selinux on deploy host
- helm部署redis-ha添加国内可访问镜像 by heyanyanchina123
- 修复多集群管理时, 若当前ezctl配置不是升级集群,会导致升级失败 by learn0208
- add ipvs配置打开strictARP #1298
- revert for supporting k8s version <= 1.26
- add kubetail, by WeiLai
- update manifests:es-cluster/mysql-cluster
