# easzctl 命令行介绍

## 为什么使用 easzctl

作为 kubeasz 项目的推荐命令行脚本，easzctl 十分轻量、简单；（后续会不断完善补充）

- 命令集 1：集群层面操作
  - 切换/创建集群 context
  - 删除当前集群
  - 显示所有集群
  - 创建集群
  - 创建单机集群（类似 minikube）
- 命令集 2：集群内部操作
  - 增加工作节点
  - 增加主节点
  - 增加 etcd 节点
  - 删除 etcd 节点
  - 删除任意节点
  - 升级集群
- 命令集3：额外操作
  - 开启/关闭基础认证

集群 context 由 ansible hosts 配置、roles 配置等组成，用以区分不同的 k8s 集群，从而实现多集群的创建和管理；当然 easzctl 命令行不是必须的，你仍旧可以使用之前熟悉的方式安装/管理集群。

典型 easzctl 创建管理的集群拓扑如下：

```
+----------------+               +-----------------+
|easzctl 1.1.1.1 |               |cluster-aio:     |
+--+---+---+-----+               |deploy 4.4.4.4   |
   |   |   |                     |master 4.4.4.4   |
   |   |   +-------------------->+etcd   4.4.4.4   |
   |   |                         |node   4.4.4.4   |
   |   +--------------+          +-----------------+
   |                  |
   v                  v
+--+------------+ +---+----------------------------+
| cluster-1:    | | cluster-2:                     |
| deploy 2.2.2.1| | deploy 3.3.3.1                 |
| master 2.2.2.1| | master 3.3.3.1/3.3.3.2         |
| etcd   2.2.2.2| | etcd   3.3.3.1/3.3.3.2/3.3.3.3 |
| node   2.2.2.3| | node   3.3.3.4/3.3.3.5/3.3.3.6 |
+---------------+ +--------------------------------+
```

## 使用 easzctl 举例

- 随时运行 `easzctl help` 获取命令行提示信息

- 1.创建 context：准备集群名称（例如：test-cluster1），运行 `easzctl checkout test-cluster1`
  - 如果 context: test-cluster1 不存在，那么会根据 default 配置创建它；如果存在则切换当前 context 为 test-cluster1

- 2.准备 context 以后，根据你的需要配置 ansible hosts 文件和其他配置，然后运行 `easzctl setup`

- 3.安装成功后，运行 `easzctl list` 显示当前所有集群信息

- 4.重复步骤 1/2 可以创建多个集群

- 5.切换到某个集群 `easzctl checkout xxxx`，然后执行增加/删除节点操作

That's it! 赶紧动手测试吧，欢迎通过 Issues 和 PRs 反馈您的意见和建议！
