## 增加 kube-master 节点

注意：目前仅支持按照本项目`多主模式`(hosts.m-masters.example/hosts.cloud.example)部署的`k8s`集群增加`master`节点

新增`kube-master`节点大致流程为：
- [可选]新节点安装 chrony 时间同步
- 新节点预处理 prepare
- 新节点安装 docker 服务
- 新节点安装 kube-master 服务
- 新节点安装 kube-node 服务
- 新节点安装网络插件相关
- 禁止业务 pod调度到新master节点
- 更新配置 haproxy 负载均衡并重启

### 操作步骤

首先配置 ssh 免密码登陆新增节点，然后执行 (假设待增加节点为 192.168.1.11)：

``` bash
$ easzctl add-master 192.168.1.11
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

