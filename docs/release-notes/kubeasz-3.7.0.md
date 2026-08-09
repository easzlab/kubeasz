## kubeasz 3.7.0

首个 kubeasz 3.7.x 版本发布，默认使用nftables 模式，包含大量组件更新。

### 版本/组件更新

- k8s: v1.36.3
- etcd: v3.7.1
- containerd: 2.3.3
- runc: v1.5.1
- cni: v1.9.1
- calico: v3.32.1
- flannel: v0.28.9
- coredns: 1.14.6
- nodelocaldns: 1.26.8
- metrics: v0.9.0
- pause: 3.10.2
- headlamp: 0.43.0
- openebs: 4.5.1
- higress: 2.2.3
- node-problem-detector: v1.35.3
- kube-prometheus-stack: 88.1.5

### 重要更新

- 支持 k8s v1.36集群默认使用 kube-proxy 的nftables 模式，需要较新的系统内核支持(kernal >= 5.13 and nft >= 1.0.0)；更新默认支持的操作系统列表，尽量选择较新系统内核的稳定linux发行版。
- 新增组件node-problem-detector 和若干自定义检查插件。
- 新增组件 headlamp 替换 dashboard，使用更加方便。
- 新增组件 higress ingress controller 替换 ingress-nginx。

### 其他更新

- 新增 easzlab/kubeasz 镜像签名信息和SBOM 信息，感谢kobihikri #1538
- 新增openeuler22/24系统软件离线包下载，感谢[whowhopipi](https://github.com/easzlab/dockerfile-kubeasz-sys-pkg/pull/12)
- 部分文档更新
