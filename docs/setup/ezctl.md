# ezctl 命令行介绍

## 为什么使用 ezctl

kubeasz 项目使用ezctl 方便地创建和管理多个k8s 集群，ezctl 使用shell 脚本封装ansible-playbook 执行命令，它十分轻量、简单和易于扩展。

### 使用帮助

随时运行 ezctl 获取命令行提示信息，如下

```
Usage: ezctl COMMAND [args]
-------------------------------------------------------------------------------------
Cluster setups:
    list		             to list all of the managed clusters
    checkout    <cluster>            to switch default kubeconfig of the cluster
    new         <cluster>            to start a new k8s deploy with name 'cluster'
    setup       <cluster>  <step>    to setup a cluster, also supporting a step-by-step way
    start       <cluster>            to start all of the k8s services stopped by 'ezctl stop'
    stop        <cluster>            to stop all of the k8s services temporarily
    upgrade     <cluster>            to upgrade the k8s cluster
    destroy     <cluster>            to destroy the k8s cluster
    backup      <cluster>            to backup the cluster state (etcd snapshot)
    restore     <cluster>            to restore the cluster state from backups
    start-aio		             to quickly setup an all-in-one cluster with 'default' settings

Cluster ops:
    add-etcd    <cluster>  <ip>      to add a etcd-node to the etcd cluster
    add-master  <cluster>  <ip>      to add a master node to the k8s cluster
    add-node    <cluster>  <ip>      to add a work node to the k8s cluster
    del-etcd    <cluster>  <ip>      to delete a etcd-node from the etcd cluster
    del-master  <cluster>  <ip>      to delete a master node from the k8s cluster
    del-node    <cluster>  <ip>      to delete a work node from the k8s cluster

Extra operation:
    kcfg-adm    <cluster>  <args>    to manage client kubeconfig of the k8s cluster

Use "ezctl help <command>" for more information about a given command.
```

- 命令集 1：集群安装相关操作
  - 显示当前所有管理的集群
  - 切换默认集群
  - 创建新集群配置
  - 安装新集群
  - 启动临时停止的集群
  - 临时停止某个集群（包括集群内运行的pod）
  - 升级集群k8s组件版本
  - 删除集群
  - 备份集群（仅etcd数据，不包括pv数据和业务应用数据）
  - 从备份中恢复集群
  - 创建单机集群（类似 minikube）
- 命令集 2：集群节点操作
  - 增加 etcd 节点
  - 增加主节点
  - 增加工作节点
  - 删除 etcd 节点
  - 删除主节点
  - 删除工作节点
- 命令集3：额外操作
  - 管理客户端kubeconfig

#### 举例创建、安装新集群流程

- 1.首先创建集群配置实例 

``` bash
~# ezctl new k8s-01
2021-01-19 10:48:23 DEBUG generate custom cluster files in /etc/kubeasz/clusters/k8s-01
2021-01-19 10:48:23 DEBUG set version of common plugins
2021-01-19 10:48:23 DEBUG cluster k8s-01: files successfully created.
2021-01-19 10:48:23 INFO next steps 1: to config '/etc/kubeasz/clusters/k8s-01/hosts'
2021-01-19 10:48:23 INFO next steps 2: to config '/etc/kubeasz/clusters/k8s-01/config.yml'
```
然后根据提示配置'/etc/kubeasz/clusters/k8s-01/hosts' 和 '/etc/kubeasz/clusters/k8s-01/config.yml'；为方便测试我们在hosts里面设置单节点集群（etcd/kube_master/kube_node配置同一个节点，注意节点需先设置ssh免密码登陆）, config.yml 使用默认配置即可。

- 2.然后开始安装集群

``` bash
# 一键安装
ezctl setup k8s-01 all

# 或者分步安装，具体使用 ezctl help setup 查看分步安装帮助信息
# ezctl setup k8s-01 01
# ezctl setup k8s-01 02
# ezctl setup k8s-01 03
# ezctl setup k8s-01 04
... 
```

- 3.重复步骤1，2可以创建、管理多个k8s集群（建议ezctl使用独立的部署节点）

ezctl 创建管理的多集群拓扑如下

```
+----------------+               +-----------------+
|ezctl 1.1.1.1   |               |cluster-aio:     |
+--+---+---+-----+               |                 |
   |   |   |                     |master 4.4.4.4   |
   |   |   +-------------------->+etcd   4.4.4.4   |
   |   |                         |node   4.4.4.4   |
   |   +--------------+          +-----------------+
   |                  |
   v                  v
+--+------------+ +---+----------------------------+
| cluster-1:    | | cluster-2:                     |
|               | |                                |
| master 2.2.2.1| | master 3.3.3.1/3.3.3.2         |
| etcd   2.2.2.2| | etcd   3.3.3.1/3.3.3.2/3.3.3.3 |
| node   2.2.2.3| | node   3.3.3.4/3.3.3.5/3.3.3.6 |
+---------------+ +--------------------------------+
```

That's it! 赶紧动手测试吧，欢迎通过 Issues 和 PRs 反馈您的意见和建议！
