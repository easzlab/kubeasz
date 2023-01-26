# 管理 node 节点

目录
- 1.增加 kube_node 节点
- 2.增加非标准ssh端口节点
- 3.删除 kube_node 节点

## 1.增加 kube_node 节点

新增`kube_node`节点大致流程为：(参考ezctl 里面add-node函数 和 playbooks/22.addnode.yml)
- [可选]新节点安装 chrony 时间同步
- 新节点预处理 prepare
- 新节点安装 container runtime
- 新节点安装 kube_node 服务
- 新节点安装网络插件相关

### 操作步骤

执行如下 (假设待增加节点为 192.168.1.11，k8s集群名为 test-k8s)：

``` bash
# ssh 免密码登录
$ ssh-copy-id 192.168.1.11

# 部分操作系统需要配置python软链接
$ ssh 192.168.1.11 ln -s /usr/bin/python3 /usr/bin/python

# 新增节点
$ ezctl add-node test-k8s 192.168.1.11

# 同理，重复上面步骤再新增节点并自定义nodename
$ ezctl add-node test-k8s 192.168.1.12 k8s_nodename=worker-03
```

### 验证

``` bash
# 验证新节点状态
$ kubectl get node

# 验证新节点的网络插件calico 或flannel 的Pod 状态
$ kubectl get pod -n kube-system

# 验证新建pod能否调度到新节点，略
```

## 2.增加非标准ssh端口节点

假设待添加节点192.168.2.1，ssh 端口 10022；然后执行 

``` bash
$ ssh-copy-id -p 10022 192.168.2.1
$ ssh -p10022 192.168.2.1 ln -s /usr/bin/python3 /usr/bin/python
$ ezctl add-node test-k8s 192.168.2.1 ansible_ssh_port=10022
```

- 注意：如果在添加节点时需要设置其他个性化变量，可以同理在后面不断添加


## 3.删除 kube_node 节点

删除 node 节点流程：(参考ezctl 里面del-node函数 和 playbooks/32.delnode.yml)
- 检测是否可以删除
- 迁移节点上的 pod
- 删除 node 相关服务及文件
- 从集群删除 node

### 操作步骤

``` bash
$ ezctl del-node test-k8s 192.168.1.11 # 假设待删除节点为 192.168.1.11
```

### 验证

略
