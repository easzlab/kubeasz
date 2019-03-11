# 个性化集群参数配置

对于刚接触项目者，如"快速指南"说明，只需要：

- **1** 个配置：`/etc/ansible/hosts`
- **1** 键安装：`ansible-playbook /etc/ansilbe/90.setup.yml`

具体来讲 `kubeasz`创建集群主要在以下两个地方进行配置：

- ansible hosts 文件（模板在examples目录）：集群主要节点定义和主要参数配置、全局变量
- roles/xxx/defaults/main.yml 文件：其他参数配置或者部分组件附加参数

## ansible hosts

项目在[快速指南](quickStart.md)或者[集群规划与安装概览](00-planning_and_overall_intro.md)已经介绍过，主要包括集群节点定义和集群范围的主要参数配置；目前提供四种集群部署模板。

- 尽量保持配置简单灵活
- 尽量保持配置项稳定

## roles/xxx/defaults/main.yml

主要包括集群某个具体组件的个性化配置，具体组件的配置项可能会不断增加；

- 可以在不做任何配置更改情况下使用默认值创建集群
- 可以根据实际需要配置 k8s 集群，常用举例
  - 配置 kube-proxy 使用 ipvs：修改 roles/kube-node/defaults/main.yml 变量 PROXY_MODE: "ipvs"
  - 配置 lb 节点负载均衡算法：修改 roles/lb/defaults/main.yml 变量 BALANCE_ALG: "roundrobin"
  - 配置 docker 国内镜像加速站点：修改 roles/docker/defaults/main.yml 相关变量
  - 配置 apiserver 支持公网域名：修改 roles/kube-master/defaults/main.yml 相关变量
  - 配置 flannel 使用镜像版本：修改 roles/flannel/defaults/main.yml 相关变量
  - 配置选择不同 addon 组件：修改roles/cluster-addon/defaults/main.yml
