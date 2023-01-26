# 管理 kube_master 节点

## 1.增加 kube_master 节点

新增`kube_master`节点大致流程为：(参考ezctl 中add-master函数和playbooks/23.addmaster.yml)
- [可选]新节点安装 chrony 时间同步
- 新节点预处理 prepare
- 新节点安装 container runtime 
- 新节点安装 kube_master 服务
- 新节点安装 kube_node 服务
- 新节点安装网络插件相关
- 禁止业务 pod调度到新master节点
- 更新 node 节点 haproxy 负载均衡并重启

### 操作步骤

执行如下 (假设待增加节点为 192.168.1.11, 集群名称test-k8s)：

``` bash
# ssh 免密码登录
$ ssh-copy-id 192.168.1.11

# 部分操作系统需要配置python软链接
$ ssh 192.168.1.11 ln -s /usr/bin/python3 /usr/bin/python

# 新增节点
$ ezctl add-master test-k8s 192.168.1.11

# 同理，重复上面步骤再新增节点并自定义nodename
$ ezctl add-master test-k8s 192.168.1.12 k8s_nodename=master-03
```

### 验证

``` bash
# 在新节点master 服务状态
$ systemctl status kube-apiserver 
$ systemctl status kube-controller-manager
$ systemctl status kube-scheduler

# 查看新master的服务日志
$ journalctl -u kube-apiserver -f

# 查看集群节点，可以看到新 master节点 Ready, 并且禁止了POD 调度功能
$ kubectl get node
NAME           STATUS                     ROLES     AGE       VERSION
192.168.1.1    Ready,SchedulingDisabled   <none>    3h        v1.9.3
192.168.1.2    Ready,SchedulingDisabled   <none>    3h        v1.9.3
192.168.1.3    Ready                      <none>    3h        v1.9.3
192.168.1.4    Ready                      <none>    3h        v1.9.3
192.168.1.11   Ready,SchedulingDisabled   <none>    2h        v1.9.3	# 新增 master节点
```

## 2.删除 kube_master 节点


删除`kube_master`节点大致流程为：(参考ezctl 中del-master函数和playbooks/33.delmaster.yml)
- 检测是否可以删除
- 迁移节点 pod
- 删除 master 相关服务及文件
- 删除 node 相关服务及文件
- 从集群删除 node 节点
- 从 ansible hosts 移除节点
- 在 ansible 控制端更新 kubeconfig
- 更新 node 节点 haproxy 配置

### 操作步骤

``` bash
$ ezctl del-master test-k8s 192.168.1.11  # 假设待删除节点 192.168.1.11
```

### 验证

略

