## kubeasz-2.2.3 发布说明

CHANGELOG:
- 组件更新
  - k8s: v1.20.1, v1.19.6, v1.18.14, v1.17.16
  - containerd v1.4.3
  - docker: 19.03.14
  - calico v3.15.3
  - dashboard: v2.1.0
- 集群安装
  - 更新支持containerd 1.4.3
  - 修改etcd启动参数auto-compaction-mode=periodic #951 by lushenle
  - 修改docker默认开启live-restore功能
- 工具脚本
  - easzup: 移除下载containerd代码，已合并在镜像easzlab/kubeasz-ext-bin:0.8.1 中
  - start-aio: 增加懒人一键下载并启动aio集群脚本 ./start-aio ${kubeasz_version}
- 文档
  - 少量文档更新
