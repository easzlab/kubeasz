## kubeasz-3.1.1 (Autumnal Equinox)

昼夜均，寒暑平，中秋祭月。kubeasz 3.1.1 小版本更新。

### 正式通过k8s一致性认证

kubeasz 用户可以确认集群各项功能符合预期，既符合k8s设计标准；

- v1.20 [已完成](https://github.com/cncf/k8s-conformance/pull/1326)
- v1.21 [已完成](https://github.com/cncf/k8s-conformance/pull/1398)
- v1.22 [进行中]

### 组件更新

- k8s: v1.22.2, v1.21.5, v.1.20.11, v1.19.15
- etcd: v3.5.0
- docker: 20.10.8
- calico: v3.19.2
- coredns: 1.8.4
- pause: 3.5
- dashboard: v2.3.1
- metrics-server: v0.5.0

### 其他

- 更新：kuboard 文档 #1014 #1023
- 更新：判断服务状态直接使用systemctl is-active #1019
- 修复：etcd dir bug #1036
- 更新：traefik为Daemonset部署,增加健康检测功能以及Node节点亲和性调度 #1028
- 更新：dashboard 部署文件和文档
- 更新：metrics-server 部署文件和文档
- 修复：coredns 1.8.4 rbac settings
- 修复：docker/containerd是否需要安装的判断条件
- 修复：暂时绕过centos7.9开启KUBE_RESERVED的问题
- 调整：docker/containerd运行时安装互不影响
