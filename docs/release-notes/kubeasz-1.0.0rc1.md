## kubeasz-1.0.0rc1 发布说明

- Note: kubeasz-1.x 第一个版本预发布，原 master 已并入 release-0 分支，停止主要更新，仅做 bug 修复版本；后续 master 分支将开始 kubeasz-1.x 版本发布。
- Action Required: 本次更新修改 ansible hosts 文件，请按照 example 目录中的对应例子修改`/etc/ansible/hosts`文件, 确保 ansible hosts 文件中主机组的顺序与例子一致。

CHANGELOG:
- 组件更新：
  - k8s: v1.13.4
  - cilium v1.4.1
- 集群安装：
  - **引入[easzctl](https://github.com/easzlab/kubeasz/blob/master/tools/easzctl)命令行工具**，后续它将作为推荐的集群常规管理工具，包括多集群管理(to do)
  - **新增 docker 运行安装 kubeasz**，请参考文档 https://github.com/easzlab/kubeasz/blob/master/docs/setup/docker_kubeasz.md
  - 优化 example hosts 配置，废弃 new-node/new-master/new-etcd 主机组，废弃变量K8S_VER，改为自动识别
  - 集成以下集群操作至 easzctl 命令行
    - [添加 master](https://github.com/easzlab/kubeasz/blob/master/docs/op/AddMaster.md)
    - [添加 node](https://github.com/easzlab/kubeasz/blob/master/docs/op/AddNode.md)
    - [添加 etcd](https://github.com/easzlab/kubeasz/blob/master/docs/op/op-etcd.md)
    - [删除 etcd](https://github.com/easzlab/kubeasz/blob/master/docs/op/op-etcd.md)
    - [删除节点](https://github.com/easzlab/kubeasz/blob/master/docs/op/clean_one_node.md)
    - [快速创建 aio 集群]()
  - 修改安装时生成随机 basic auth 密码
  - 修改优化部分安装脚本以兼容 docker 运行 kubeasz
  - update cilium v1.4.1，更新 cilium 文档(to do)
  - 增加启动 kubeasz 容器的脚本 tools/kubeasz-docker
- 其他：
  - 修复兼容 docker 18.09.x 版本安装
