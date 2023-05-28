## kubeasz 3.6.1

kubeasz 3.6.1 发布：支持k8s v1.27版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.27.2
- calico: v3.24.6
- kube-ovn: v1.11.5
- kube-router: v1.5.4

### 增加应用部署插件 kubeapps

Kubeapps 是一个基于 Web 的应用程序，它可以在 Kubernetes 集群上进行一站式安装，并使用户能够部署、管理和升级应用
程序。https://github.com/easzlab/kubeasz/blob/master/docs/guide/kubeapps.md

### 重要更新

- 重写`ezdown`脚本支持下载额外的应用容器镜像
- 增加`local-path-provisioner`本地文件目录提供者
- 设置允许kubelet并行拉取容器镜像

### 其他

- 增加kubectl-node-shell 脚本
- 修复ansible connect local 是 python 解析器不确定问题
- 修复typo #1273
- 部分文档更新
