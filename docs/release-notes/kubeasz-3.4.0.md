## kubeasz 3.4.0 (White Dew)

蒹葭苍苍，白露为霜。kubeasz 3.4.0 发布更新支持k8s 1.25 版本。

### 组件更新

- k8s: v1.25.1
- containerd: 1.6.8
- calico: v3.23.3
- cilium: 1.12.2
- flannel: v0.19.2
- kube-prometheus-stack: 39.11.0
- nodelocaldns: 1.22.8
- dashboard: v2.6.1
- pause: 3.8

### 其他

- fix: custom PATH settings
- 修复calico ipip隧道模式说明错误，并完善可选参数说明以及使用场景 (#1168 by Hello-Linux)
- fix: checking bash shell (#1171 by EamonZhang)
- fix: create etcd certs (#1172 by EamonZhang )
- fix: ca-config.json format (#1174 by libinglong)
