## kubeasz-0.3.0 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.11.2/v1.10.6/v1.9.10/v1.8.15
  - calico: v3.1.3
  - kube-router: v0.2.0-beta.9 
- 功能更新：
  - **增加集群备份与恢复** 功能与[说明](https://github.com/easzlab/kubeasz/blob/master/docs/op/cluster_restore.md)
  - **增加cilium网络插件** ，文档待更新
  - **增加cluster-storage角色** 与[文档说明](https://github.com/easzlab/kubeasz/blob/master/docs/setup/08-cluster-storage.md)
  - 增加阿里云NAS存储支持
  - 增加集群个性化[配置说明](https://github.com/easzlab/kubeasz/blob/master/docs/setup/config_guide.md)与生成脚本`tools/init_vars.yml` 
  - 支持deploy节点与ansible执行节点分离，为一份代码创建多个集群准备
- 其他：
  - 更新 jenkins and plugins (#258)
  - 重写 nfs动态存储脚本与文档
  - 优化 cluster-addon 安装脚本
  - 增加 docker 配置文件
  - 更新 offline images 0.3
  - 增加 batch/v2alpha支持
  - 移动 DNS yaml文件至 /opt/kube/kube-system
  - fix 多主集群下change_k8s_network时vip丢失问题
  - fix 禁止节点使用系统swap
  - fix 解压后的harbor安装文件没有执行权限问题
  - fix Ubuntu 18.04无法安装haproxy、keepalived问题
