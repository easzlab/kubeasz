## kubeasz 3.6.6

kubeasz 3.6.6 发布：支持k8s v1.32 版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.32.3
- etcd: v3.5.20
- containerd: 2.0.4
- runc: v1.2.6
- calico: v3.28.3
- coredns: 1.11.4
- cni: v1.6.2
- harbor: v2.12.2

### 更新

- 更新国内docker镜像仓库加速设置，解决ezdown脚本无法下载镜像问题；同步更新containerd 镜像仓库加速设置
- 主要组件大版本更新：containerd 从 1.7.x 更新大版本 2.0.x，更新主要配置文件；runc 从 1.1.x 更新大版本 1.2.x
- 安装逻辑更新：新增节点不再重复执行网络插件安装，避免部分网络插件自动重启业务pod，by gogeof
- 安装逻辑更新：每次执行脚本 containerd 都会被重新安装，不管原先是否已经运行
- 优化更新 ezctl 脚本从 ezdown 加载变量方式，by RadPaperDinosaur


### 其他

- 修复 CLUSTER_DNS_SVC_IP & CLUSTER_KUBERNETES_SVC_IP 地址生成规则，by yunpiao
- 更新conformance文档
- 
