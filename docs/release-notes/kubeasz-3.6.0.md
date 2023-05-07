## kubeasz 3.6.0 (Beginning of Summer)

微雨过，小荷翻。榴花开欲然。kubeasz 3.6.0 发布：支持k8s v1.27版本，支持更多操作系统安装，以及组件更新和一些bugfix。

### 版本更新

- k8s: v1.27.1
- cilium: v1.13.2
- flannel: v0.21.4
- harbor: v2.6.4
- metrics-server: v0.6.3
- k8s-dns-node-cache: 1.22.20
- kube-prometheus-stack: 45.23.0

### 调整项目分支更新规则

k8s大版本对应kubeasz特定的大版本号，详见README.md 中版本对照表，当前积极更新的分支如下：

- master：默认保持与最新分支同步，当前与v3.6同步
- v3.6：对应k8s v1.27 版本，持续保持更新
- v3.5：对应k8s v1.26 版本，主要使用cherry-pick方式合并后续版本中的重要commit
- v3.4：对应k8s v1.25 版本，主要使用cherry-pick方式合并后续版本中的重要commit
- v3.3：对应k8s v1.24 版本，主要使用cherry-pick方式合并后续版本中的重要commit

### 支持更多操作系统安装

本次增加测试支持大部分使用systemd的linux发行版，如果安装有问题先请查看(docs/setup/multi_os.md)；如果某个能够支持安装的系统没有在列表中，欢迎提PR 告知。

- **Alibaba Linux** 2.1903, 3.2104([notes](docs/setup/multi_os.md#Alibaba))
- **Alma Linux** 8, 9
- **Anolis OS** 8.x RHCK, 8.x ANCK([notes](docs/setup/multi_os.md#Anolis))
- **CentOS/RHEL** 7, 8, 9
- **Debian** 10, 11([notes](docs/setup/multi_os.md#Debian))
- **Fedora** 34, 35, 36, 37
- **openSUSE** Leap 15.x([notes](docs/setup/multi_os.md#openSUSE))
- **Rocky Linux** 8, 9
- **Ubuntu** 16.04, 18.04, 20.04, 22.04

### 重要更新

- 重写`ezdown`脚本支持下载多系统软件包部分
- 重写`role:prepare`支持离线安装多系统软件包部分
- 简化harbor安装后集成使用，目前在containerd容器运行时中额外配置允许insecure仓库方式
- 修复pod挂载 hostpath volume，删除pod会卡住问题 (#1259) by itswl
- 增加设置limits for pids #1265 by AsonZhang

### 其他

- 增加项目`ISSUE`模版 
- 修复chronyd 服务可能出现 enable失败问题 (#1254) by Roach57
- 增加ezctl setup脚本执行时打印版本信息
