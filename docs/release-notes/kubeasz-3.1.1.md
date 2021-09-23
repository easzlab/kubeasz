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

### 其他

- fix:增加/删除节点时ansible hosts文件更新错误 
