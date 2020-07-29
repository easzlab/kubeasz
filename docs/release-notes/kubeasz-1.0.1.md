## kubeasz-1.0.1 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.13.5 v1.12.7 v1.11.9
  - cni v0.7.5
  - coredns 1.4.0
- 集群安装：
  - 优化journald日志服务配置，避免与syslog采集重复，节省节点资源
  - 修复CVE-2019-3874 (work around)，[详情](https://mp.weixin.qq.com/s/CnzK8722pJUWRAitWBRPcw)
  - 修复首个etcd成员故障时apiserver也故障的bug，详见 kubernetes issue #72102
  - 修复add-master时偶然出现的兼容性问题 issue #490
  - 修复apiserver启用basic_auth认证时用户rbac设置问题
  - 修复docker安装时变量 DOCKER_VER 需要显式转换成 float
  - 调整ca证书有效期等配置
  - 增加kubectl使用可选参数配置kubeconfig
- easzctl 命令行 
  - 新增[升级集群](https://github.com/easzlab/kubeasz/blob/master/docs/op/upgrade.md)
  - 修复easzctl setup启动前检查项
  - 增加/删除etcd节点成功后配置重启apiserver
- 其他：
  - 更新 DOCKER_VER 使用新版本格式 #488
  - fix download url for harbor v1.5.x #492
