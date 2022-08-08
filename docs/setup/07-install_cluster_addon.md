# 07-安装集群主要插件

目前挑选一些常用、必要的插件自动集成到安装脚本之中:  

## 集群默认安装

- [coredns](../guide/kubedns.md)
- [nodelocaldns](../guide/kubedns.md)
- [metrics-server](../guide/metrics-server.md)
- [dashboard](../guide/dashboard.md)

kubeasz 默认安装上述基础插件，并支持离线方式安装(./ezdown -D 命令会自动下载组件镜像，并推送到本地镜像仓库easzlab.io.local:5000)

## 集群可选安装

- [prometheus](../guide/prometheus.md)
- [network_check](network-plugin/network-check.md)
- [nfs_provisioner]()

kubeasz 默认不安装上述插件，可以在配置文件(clusters/xxx/config.yml)中开启，支持离线方式安装(./ezdown -X 会额外下载这些组件镜像，并推送到本地镜像仓库easzlab.io.local:5000)

## 安装脚本

详见`roles/cluster-addon/` 目录

- 1.根据hosts文件中配置的`CLUSTER_DNS_SVC_IP` `CLUSTER_DNS_DOMAIN`等参数生成kubedns.yaml和coredns.yaml文件
- 2.注册变量pod_info，pod_info用来判断现有集群是否已经运行各种插件
- 3.根据pod_info和`配置开关`逐个进行/跳过插件安装

## 下一步

- [创建ex_lb节点组](ex-lb.md), 向集群外提供高可用apiserver
- [创建集群持久化存储](08-cluster-storage.md)
