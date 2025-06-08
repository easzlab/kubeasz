## kubeasz 3.6.7

kubeasz 3.6.7 发布：支持k8s v1.33 版本，组件更新和bugfix。

### 版本更新

- k8s: v1.33.1
- etcd: v3.5.21
- containerd: 2.1.1
- runc: v1.2.6
- calico: v3.28.4
- cilium: 1.17.4
- coredns: 1.12.1
- cni: v1.7.1
- dnsNodeCache: 1.25.0
- harbor: v2.12.4
- local-path-provisioner: v0.0.31
- dashboard 7.12.0

### 更新

- 增加可选组件`kubeblocks`集成，增加多种数据库高可用方案
- 重写脚本ezdown中关于镜像下载保存部分，清理冗余，增加错误错误处理
- 修复添加/删除master节点时/etc/hosts问题 #1464
- 修复使用静态编译的containerd二进制，并设置日志为warn级别，避免当容器使用exec类健康检查时产生过多日志
- 修复./ezdown -D 偶发403报错 #1470
- 修复cilium 组件原cilium_connectivity_check脚本执行条件

### 文档更新

- 更新一致性认证文档 conformance.md
