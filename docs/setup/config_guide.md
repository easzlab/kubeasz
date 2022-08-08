# 个性化集群参数配置

`kubeasz`创建集群主要在以下两个地方进行配置：(假设集群名xxxx)

- clusters/xxxx/hosts 文件（模板在example/hosts.multi-node）：集群主要节点定义和主要参数配置、全局变量
- clusters/xxxx/config.yml（模板在examples/config.yml）：其他参数配置或者部分组件附加参数

## clusters/xxxx/hosts (ansible hosts)

如[集群规划与安装概览](00-planning_and_overall_intro.md)中介绍，主要包括集群节点定义和集群范围的主要参数配置

- 尽量保持配置简单灵活
- 尽量保持配置项稳定

常用设置项：

- 修改容器运行时: CONTAINER_RUNTIME="containerd"
- 修改集群网络插件：CLUSTER_NETWORK="calico"
- 修改容器网络地址：CLUSTER_CIDR="192.168.0.0/16"
- 修改NodePort范围：NODE_PORT_RANGE="30000-32767"

## clusters/xxxx/config.yml

主要包括集群某个具体组件的个性化配置，具体组件的配置项可能会不断增加；可以在不做任何配置更改情况下使用默认值创建集群

根据实际需要配置 k8s 集群，常用举例

- 配置使用离线安装系统包：INSTALL_SOURCE: "offline" （需要ezdown -P 下载离线系统软件）
- 配置CA证书以及其签发证书的有效期
- 配置 apiserver 支持公网域名：MASTER_CERT_HOSTS
- 配置 cluster-addon 组件安装
- ...
