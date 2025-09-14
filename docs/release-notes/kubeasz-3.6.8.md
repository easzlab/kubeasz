## kubeasz 3.6.8

kubeasz 3.6.8 发布：支持k8s v1.34 版本，组件更新和bugfix。

### 版本更新

- k8s: v1.34.1
- etcd: v3.6.4
- containerd: 2.1.4
- runc: v1.3.1
- coredns: 1.12.4
- cni: v1.8.0
- dnsNodeCache: 1.26.4
- metrics: v0.8.0
- flannel: v0.27.3
- kubeblocks: 1.0.0
- kube-prometheus-stack: 75.7.0

### 重要更新

- 调整系统内核设置 commit f9bdbeb4e3bd6b98a03a900d3e50ef29da6a590f, #1478
- 新增支持 openEuler 22.03 LTS, 24.03 LTS
- 优化节点只需运行一次 prepare task
- 增加可选开启集群审计功能
- 修复 calico mtu 设置 #1444
- 修复 calico vxlan overlay 设置 #1492 
- 更新 containerd 配置容器镜像仓库方式

### 文档更新

- 实验性混合架构部署文档 https://github.com/easzlab/kubeasz/blob/master/docs/setup/mix_arch.md
- updat kernel_upgrade.md for centos7 by Zlanghu #1483

感谢新增贡献者：

vistamin #1444
newfzk #1477
learn0208 #1478
Zlanghu #1483
newfzk #1492
TOT-JIN #1495
