## 增加 kube-node 节点

新增`kube-node`节点大致流程为：
- 1. 节点预处理 prepare
- 2. 安装 kubectl (可选)
- 3. 安装 docker
- 4. 安装 kubelet 和 kube-proxy
- 5. 准备网络插件 calico 或 flannel
- 6. 批准新节点 kubectl certificate approve

### 操作步骤

按照本项目说明，首先确保deploy节点能够ssh免密码登陆新增节点，然后在**deploy**节点执行两步：

- 1. 修改ansible hosts 文件，在 [new-node] 组编辑需要新增的节点，例如：

``` bash
...
# 预留组，后续添加node节点使用
[new-node]
192.168.1.6 NODE_ID=node6 NODE_IP="192.168.1.6"
...
```
- 2. 执行安装脚本

``` bash
$ cd /etc/ansible && ansible-playbook 20.addnode.yml
```

### 验证

``` bash
# 验证新节点状态
$ kubectl get node

# 验证新节点的网络插件calico 或flannel 的Pod 状态
$ kubectl get pod -n kube-system

# 验证新建负载能否调度到新节点，略
```

### 后续

上述步骤验证成功，确认新节点工作正常后，为了方便后续再次添加节点，在ansible hosts文件中，把 [new-node] 组下的节点全部复制到 [kube-node] 组下，并清空 [new-node] 组的节点。
