# 管理 node 节点

目录
- 1.增加 kube-node 节点
- 2.增加非标准ssh端口节点
- 3.删除 kube-node 节点

## 1.增加 kube-node 节点

新增`kube-node`节点大致流程为：tools/02.addnode.yml
- [可选]新节点安装 chrony 时间同步
- 新节点预处理 prepare
- 新节点安装 docker 服务
- 新节点安装 kube-node 服务
- 新节点安装网络插件相关

### 操作步骤

首先配置 ssh 免密码登录新增节点，然后执行 (假设待增加节点为 192.168.1.11)：

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

## 2.增加非标准ssh端口节点

目前 easzctl 暂不支持自动添加非标准 ssh 端口的节点，可以手动操作如下：

- 假设待添加节点192.168.2.1，ssh 端口 10022；配置免密登录 ssh-copy-id -p 10022 192.168.2.1，按提示输入密码
- 在 /etc/ansible/hosts文件 [kube-node] 组下添加一行：
```
192.168.2.1 ansible_ssh_port=10022
```
- 最后执行 `ansible-playbook /etc/ansible/tools/02.addnode.yml -e NODE_TO_ADD=192.168.2.1`

## 3.删除 kube-node 节点

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
