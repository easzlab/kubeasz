# 07-安装集群主要插件

目前挑选一些常用、必要的插件自动集成到安装脚本之中:  
- [自动脚本](../../roles/cluster-addon/tasks/main.yml)
- 配置开关
  - 参照[配置指南](config_guide.md)，生成后在`roles/cluster-addon/defaults/main.yml`配置

## 脚本介绍

- 1.根据hosts文件中配置的`CLUSTER_DNS_SVC_IP` `CLUSTER_DNS_DOMAIN`等参数生成kubedns.yaml和coredns.yaml文件
- 2.注册变量pod_info，pod_info用来判断现有集群是否已经运行各种插件
- 3.根据pod_info和`配置开关`逐个进行/跳过插件安装

## 下一步

- [创建ex-lb节点组](ex-lb.md), 向集群外提供高可用apiserver
- [创建集群持久化存储](08-cluster-storage.md)
