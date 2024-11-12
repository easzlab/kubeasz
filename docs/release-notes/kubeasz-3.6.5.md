## kubeasz 3.6.5

kubeasz 3.6.5 发布：支持k8s v1.31 版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.31.2
- etcd: v3.5.16
- containerd: 1.7.23
- runc: v1.1.15
- calico: v3.28.2
- coredns: 1.11.3
- dnsnodecache: 1.23.1
- cilium: 1.16.3
- flannel: v0.26.0
- cni: v1.6.0
- harbor: v2.11.1
- metrics-server: v0.7.2
- pause: 3.10

### 更新

- 修正centos9 下prepare脚本运行的问题 #1397 By GitHubAwan
- style: trim trailing whitespace & add logger source line number #1413 By kelein
- 操作系统：增加测试支持 Ubuntu 2404
  - 修复在ubuntu 2404上使用网络插件calico ipSet兼容性问题（calico v3.28.2）

### 其他

- 修复calico hostname 设置
- 更新部分文档
- 
