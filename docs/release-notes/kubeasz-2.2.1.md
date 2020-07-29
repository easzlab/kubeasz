## kubeasz-2.2.1 发布说明

CHANGELOG:
- 组件更新
  - k8s: v1.18.3, v1.17.6, v1.16.10, v1.15.12
  - docker: 19.03.8
  - coredns: v1.6.7
  - calico: v3.8.8
  - flannel: v0.12.0-amd64
  - pause: 3.2
  - dashboard: v2.0.1
  - easzlab/kubeasz-ext-bin:0.5.2
- 集群安装
  - 更新etcd 安装参数 #823 
  - 修改kubelet.service部分启动前置条件
  - 预处理增加设置内核参数net.core.somaxconn = 32768
- 工具脚本
  - easzup: 调整下载/安装docker等
  - docker-tag: 修复原功能，增加支持harbor镜像查询 #814 #824
- 文档
  - 更新首页连接文档：Allinone安装文档、离线安装文档等
  - 更新helm3文档，举例redis集群安装
  - 增加文档 kubesphere 安装 #804
- 其他
  - 删除'azure.cn'的docker镜像加速
  - fix: 节点角色删除时/opt/kube/bin被误删 #837
  - fix：pause-amd64:3.2 repos
  - fix: 部分文档错误
