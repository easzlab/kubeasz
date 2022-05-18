## kubeasz 3.2.0 (Slight Cold)

小寒料峭, 盼雪迎春。kubeasz 3.2.0 发布更新支持k8s 1.23 版本。

### 主要更新

- 修改默认容器运行时为containerd，如果需要仍旧使用docker，请对应修改clusters/${集群名}/hosts 配置项`CONTAINER_RUNTIME`
- 修改默认network-plugin:calico
- 修复CNCF Conformance Test，选择ipvs模式时sessionAffinity的问题
- 调整containerd配置文件与版本格式一致
- 调整kube-scheduler启动配置文件
- 调整kube-proxy启用配置文件
 
### 组件更新

- k8s: v1.23.1
- etcd: v3.5.1
- containerd: 1.5.8
- calico: v3.19.3
- flannel: v0.15.1
- coredns: 1.8.6
- cni-plugins: v1.0.1
- pause: 3.6
- dashboard: v2.4.0
- metrics-server: v0.5.2
- k8s-dns-node-cache: 1.21.1
- nfs-provisioner: v4.0.2

### 其他

- fix: avoid cleaning iptables rules on docker setup
- fix: controller-manager health check issue #1084
- feat: Add docker proxy config at ezdown
- fix: kubectl drain 参数版本差异导致失败的问题
- fix: prepare阶段的一些小问题
- fix: nf_conntrack模块安装判断等
