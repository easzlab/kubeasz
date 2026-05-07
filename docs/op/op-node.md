# 管理 node 节点

目录
- 1.增加 kube_node 节点
- 2.批量增加 kube_node 节点
- 3.增加非标准ssh端口节点
- 4.删除 kube_node 节点

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

## 2.批量增加 kube_node 节点

当一次需要新增多台 work 节点时，可使用 `add-nodes` 一次性提交多个 IP，复用 Ansible 自身的并行能力，避免对每台节点重复 fork `ansible-playbook`。底层调用的仍然是 `playbooks/22.addnode.yml`，流程与 `add-node` 完全一致。

### 操作步骤

执行如下 (假设待增加节点为 192.168.1.21、192.168.1.22、192.168.1.23，k8s集群名为 test-k8s)：

``` bash
# ssh 免密码登录（每个新节点都需要先做）
$ ssh-copy-id 192.168.1.21
$ ssh-copy-id 192.168.1.22
$ ssh-copy-id 192.168.1.23

# 一次性批量新增节点
$ ezctl add-nodes test-k8s 192.168.1.21 192.168.1.22 192.168.1.23
```

### `add-node` 与 `add-nodes` 如何选择

| 场景 | 选用 |
| --- | --- |
| 只新增一台节点 | `add-node` |
| 需要给新节点附带 host_vars（如 `ansible_ssh_port=10022`、`k8s_nodename=worker-03` 等） | `add-node`（参见 §3） |
| 一次新增多台节点，且无需为每台节点单独配置 host_vars | `add-nodes` |

> 说明：`add-nodes` 的所有位置参数都按 IP 校验，不接受 inventory 行内变量。如果批量新增的节点中有部分需要 host_vars，请对这部分节点单独使用 `add-node` 添加。

## 3.增加非标准ssh端口节点

假设待添加节点192.168.2.1，ssh 端口 10022；然后执行 

``` bash
$ ssh-copy-id -p 10022 192.168.2.1
$ ezctl add-node test-k8s 192.168.2.1 ansible_ssh_port=10022
```

- 注意：如果在添加节点时需要设置其他个性化变量，可以同理在后面不断添加


## 4.删除 kube_node 节点

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
