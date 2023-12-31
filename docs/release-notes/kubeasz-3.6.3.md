## kubeasz 3.6.3

kubeasz 3.6.3 发布：支持k8s v1.29版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.29.0
- etcd: v3.5.10
- containerd: 1.6.26
- runc: v1.1.10
- calico: v3.26.4
- cilium: 1.14.5

### 修改kubeasz支持k8s版本对应规则 

原有模式每个k8s大版本都有推荐对应的kubeasz版本，这样做会导致kubeasz版本碎片化，追踪问题很麻烦，而且也影响普通用户安装体验。从kubeasz 3.6.2版本开始，默认最新版本kubeasz兼容支持安装最新的三个k8s大版本。具体安装说明如下：

(如果/etc/kubeasz/bin 目录下已经有kube* 文件，需要先删除 rm -f /etc/kubeasz/bin/kube*)

- 安装 k8s v1.29: 使用 kubeasz 3.6.3，执行./ezdown -D 默认下载即可
- 安装 k8s v1.28: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.28.5 下载
- 安装 k8s v1.27: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.27.9 下载
- 安装 k8s v1.26: 使用 kubeasz 3.6.2，执行./ezdown -D -k v1.26.12 下载

### 重要更新

- deprecated role: os-harden，因为扩大支持更多linux发行版，系统加固方式无法在各种系统上充分测试，感谢 #1338 issue 反馈问题 
- adjust docker setup scripts
- update harbor v2.8.4 and fix harbor setup
- fix nodelocaldns yaml

### 其他

- docs update: add argocd guide 
- docs: fix the quickStart.md url in network-plugin
