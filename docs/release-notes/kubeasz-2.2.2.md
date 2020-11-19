## kubeasz-2.2.2 发布说明

CHANGELOG:
- 组件更新
  - k8s: v1.19.4, v1.18.12, v1.17.14
  - docker: 19.03.13
  - etcd: v3.4.13
  - coredns: v1.7.1
  - cni-plugins: v0.8.7
  - flannel: v0.13.0-amd64
  - dashboard: v2.0.4
- 集群安装
  - 替换apiserver参数--basic-auth-file为--token-auth-file
  - kubelet启动参数修改for debian 10 #912
  - 修复debian 10 默认iptables问题 #909
  - roles/calico/defaults/main.yaml 增加 CALICO_NETWORKING_BACKEND 变量 #895
- 工具脚本
  - easzup: 调整部分下载脚本, 增加下载containerd #918
- 文档
  - Update kuboard.md #861
- 其他
  - 调整etcd备份文件 #902 #932 #933
