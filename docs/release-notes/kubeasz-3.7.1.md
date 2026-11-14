## kubeasz 3.7.1

kubeasz 3.7.1 发布：支持k8s v1.36 版本，组件更新和bugfix。

### 版本/组件更新

- k8s: v1.36.6
- etcd: v3.7.3
- containerd: 2.3.7
- runc: v1.5.3
- cni: v1.10.0
- calico: v3.32.4
- flannel: v0.28.10
- coredns: 1.14.7
- metrics: v0.9.1
- harbor: v2.14.6

### 重要更新

- 支持 k8s v1.36集群默认使用 kube-proxy 的nftables 模式，需要较新的系统内核支持(kernal >= 5.13 and nft >= 1.0.0)；当没满足条件时也支持ipvs模式。
- 默认使用cgroup v2，当系统不支持时回退使用cgroup v1。
- 更新ansible镜像 ansibe-core:2.16.19

### 其他更新

- 增加友情链接：[GUI K8S 管理平台](https://github.com/liuhaixiao837-collab/k8smanager)
- 修复在windows git clone 仓库问题 #1543
- 修复harbor 安装找不到docker问题 
- 增加harbor 单独安装脚本 8a22e691f25b59926ca416eefa06756455b0904b
- 增加脚本从admin kubeconfig自动生成debug kubeconfig，适合AI排查问题 3c9ebc585ab4c754af06231355186b79e624b799
