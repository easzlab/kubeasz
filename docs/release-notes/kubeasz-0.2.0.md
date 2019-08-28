## kubeasz-0.2.0 发布说明

CHANGELOG:
- 组件更新：
  - 增加新网络插件 kube-router，可在ansible hosts配置`CLUSTER_NETWORK="kube-router"`
- 功能更新：
  - 增加IPVS/LVS服务代理模式，比默认的kube-proxy服务代理更高效；在选择kube-router网络插件时配置`SERVICE_PROXY="IPVS"`
  - 增加部署metrics-server，以替代heapster 提供metrics API
  - 增加自动集成安装kube-dns/dashboard等组件，可在`roles/cluster-addon/defaults/main.yml`配置
- 脚本更新：
  - 增加删除单个节点脚本，docs/op/del_one_node.md
  - 增加等待网络插件正常运行
  - Bug fix: 更新99.clean.yml清理脚本，解决集群重装后cni地址分配问题 kubernetes #57280
  - Bug fix: 从0.1.0版本升级时，kube-apiserver服务启动失败问题
- 其他：
  - 修改部分镜像拉取策略统一为：`imagePullPolicy: IfNotPresent`
  - 新增metrics-server、cluster-addon文档
  - 更新kube-router相关文档 
  - 更新集群升级说明文档 docs/op/upgrade.md
