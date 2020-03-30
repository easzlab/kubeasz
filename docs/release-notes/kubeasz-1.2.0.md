## kubeasz-1.2.0 发布说明

IMPORTANT: 本次更新增加容器运行时`containerd`支持，需要在`ansible hosts`中增加全局变量`CONTAINER_RUNTIME`（可选 docker/containerd ），参考 example/ 中的例子。

NOTE: kubeasz 项目正式从 github.com/gjmzj/kubeasz 迁移至组织 github.com/easzlab/kubeasz

CHANGELOG:
- 组件更新：
  - k8s: v1.14.2
  - traefik v1.7.11
  - efk: es/kibana 6.6.1
- 集群安装：
  - 增加 containerd 支持及[简单介绍](https://github.com/easzlab/kubeasz/blob/master/docs/guide/containerd.md)
  - 增加 EFK 日志清理工具及[说明](https://github.com/easzlab/kubeasz/blob/master/docs/guide/efk.md#%E7%AC%AC%E5%9B%9B%E9%83%A8%E5%88%86%E6%97%A5%E5%BF%97%E8%87%AA%E5%8A%A8%E6%B8%85%E7%90%86)
  - 增加 Amazon Linux 支持 by lusyoe
  - 更新 containerd/docker 仓库国内镜像设置
  - 增加 containerd 与 harbor 集成
  - 更新集群清理、离线镜像推送等脚本以支持 containerd 集成
- easzctl 命令行 
  - 修复`easzctl basic-auth`命令执行问题 #544
- 文档：
  - 更新 efk 文档
  - 更新集群节点规划文档
  - 更新公有云部署文档
  - 增加 AWS 高可用部署文档 by lusyoe
  - 更新腾讯云部署文档 by waitingsong
  - 更新安装文档 istio v1.1.7 by waitingsong
  - 更新 harbor 文档
- 其他：
  - 更新项目 logo
  - fix: 清理node时罕见错误删除hosts中其他node信息 #541
  - fix: 在没有创建集群context下运行`easzctl add-node`成功时返回失败提示 
  - 更新项目迁移部分 URL 连接内容
