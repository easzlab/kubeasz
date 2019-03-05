## 增加 kube-master 节点

注意：目前仅支持按照本项目`多主模式`(hosts.m-masters.example)部署的`k8s`集群增加`master`节点

新增`kube-master`节点大致流程为：
- LB节点重新配置 haproxy并重启 haproxy服务
- 新节点预处理 prepare
- 新节点安装 docker 服务
- 新节点安装 kube-master 服务
- 新节点安装 kube-node 服务
- 新节点安装网络插件相关
- 禁止业务 pod调度到新master节点
- 修改hosts文件，把 new-master 组成员转移到 kube-master 组

### 操作步骤

按照本项目说明，首先确保deploy节点能够ssh免密码登陆新增节点，然后在**deploy**节点执行两步：

- 修改ansible hosts 文件，在 [new-master] 组添加新增的节点，举例如下：

``` bash
...
[new-master]
192.168.1.5                 	# 新增 master节点

```

- 执行安装脚本

``` bash
$ ansible-playbook /etc/ansible/21.addmaster.yml
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
192.168.1.5    Ready,SchedulingDisabled   <none>    2h        v1.9.3	# 新增 master节点
```

### 后续

上述步骤验证成功，确认新节点工作正常后，为了方便后续再次添加节点，在ansible hosts文件中，把 [new-master] 组下的节点全部复制到 [kube-master] 组下，并清空 [new-master] 组中的节点。

- 注：新版本 kubeasz 已经自动完成 new-master 组成员转移到 kube-master 组
