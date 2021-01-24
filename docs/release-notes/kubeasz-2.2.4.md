## kubeasz-2.2.4 发布说明

CHANGELOG:
- 组件更新
  - k8s: v1.20.2
  - kube-ovn: 1.5.3
- 集群安装
  - fix: 删除etcd节点hosts文件不更新
  - fix: ubuntu 20.04安装集群dns问题 #970
  - fix: kube-proxy的metrics绑定非本地环回，支持prometheus 指标拉取 #971
  - fix: 清理脚本容器目录无法删除问题
  - minor fix: calico离线镜像下载
  - fix: calico 网络 backend 设置为 vxlan none 时，calico 部署失败 #959
  - fix: 单机安装报错"/etc/ansible/bin"不是目录 #957
  - 更新traefik v2安装方式 #955
- 文档
  - KubeSphere guide updated #968
