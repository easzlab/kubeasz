## kubeasz-2.0.3 发布说明

**WARNNING:** 从 kubeasz 2.0.1 版本开始，项目仅支持 kubernetes 社区最近的4个大版本，当前为 v1.12/v1.13/v1.14/v1.15，更老的版本不保证兼容性。

CHANGELOG:
- 组件更新
  - k8s: v1.15.2 v1.14.5 v1.13.9
  - docker: 18.09.8
  - kube-ovn: 0.6.0 #644
- 集群安装
  - 修复增加/删除 etcd 节点的脚本（当待删除节点不可达时） 
  - 修复删除 master/node 节点的脚本（当待删除节点不可达时）
  - 修复 etcd 备份脚本
  - 设置 kube-proxy 默认使用 ipvs 模式 
  - 增加部分内核优化参数
  - 增加新节点后推送常用集群插件镜像 #650
  - 增加 Docker 安装后的内部（非安全）仓库 #651
  - 增加 flannel vxlan 可选开启 DirectRouting 特性 #652
  - 禁用内核参数 net.ipv4.tcp_tw_reuse
  - 使用 netaddr 模块进行 ip 地址计算 #658
- 工具脚本
  - 新增 tools/imgutils 方便拉取 gcr.io 等仓库镜像；方便批量存储/导入离线镜像等
- 文档
  - 更新 istio.md #641
  - 更新[修改 apiserver 证书文档](https://github.com/easzlab/kubeasz/blob/master/docs/op/ch_apiserver_cert.md)
- 其他
  - fix: easzctl 删除节点时的正则匹配
  - fix: kube-ovn 启动参数设置 #658
