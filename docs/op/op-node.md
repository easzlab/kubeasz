# 管理 node 节点

## 1.增加 kube-node 节点

新增`kube-node`节点大致流程为：tools/02.addnode.yml
- [可选]新节点安装 chrony 时间同步
- 新节点预处理 prepare
- 新节点安装 docker 服务
- 新节点安装 kube-node 服务
- 新节点安装网络插件相关

### 操作步骤

首先配置 ssh 免密码登陆新增节点，然后执行 (假设待增加节点为 192.168.1.11)：

``` bash
$ easzctl add-node 192.168.1.11
```

### 验证

``` bash
# 验证新节点状态
$ kubectl get node

# 验证新节点的网络插件calico 或flannel 的Pod 状态
$ kubectl get pod -n kube-system

# 验证新建pod能否调度到新节点，略
```

## 2.删除 kube-node 节点

删除 node 节点流程：tools/12.delnode.yml
- 检测是否可以删除
- 迁移节点上的 pod
- 删除 node 相关服务及文件
- 从集群删除 node

### 操作步骤

``` bash
$ easzctl del-node 192.168.1.11 # 假设待删除节点为 192.168.1.11
```

### 验证

略
