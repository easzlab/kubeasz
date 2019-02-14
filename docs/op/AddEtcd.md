## 增加 etcd 集群节点

etcd 集群支持在线改变集群成员节点，可以增加、修改、删除成员节点；不过改变成员数量仍旧需要满足集群成员多数同意原则（quorum），另外请记住集群成员数量变化的影响：

- 增加 etcd 集群节点，提高集群稳定性
- 增加 etcd 集群节点，提高集群读性能（所有节点数据一致，客户端可以从任意节点读取数据）
- 增加 etcd 集群节点, 降低集群写性能（所有节点数据一致，每一次写入会需要所有节点数据同步）

新增`new-etcd`节点大致流程为：
- 在原有集群节点执行 member add 命令
- 新节点预处理 prepare
- 新节点安装 etcd 服务运行
- 原有集群节点以新配置重启服务运行
- 操作修改ansible hosts 文件

### 操作步骤

按照本项目说明，首先确保deploy节点能够ssh免密码登陆新增节点，然后在**deploy**节点执行两步：

- 修改ansible hosts 文件，在 [new-etcd] 组编辑需要新增的节点，例如：

``` bash
...
# 预留组，后续添加etcd节点使用
[new-etcd]
192.168.1.6      #新增etcd节点
...
```
- 执行安装脚本

``` bash
$ ansible-playbook /etc/ansible/19.addetcd.yml
```

### 验证

``` bash
# 登陆任意etcd节点验证etcd集群状态
$ export ETCDCTL_API=3 
$ etcdctl member list

# 验证所有etcd节点服务状态和日志
$ systemctl status etcd
$ journalctl -u etcd -f

# 检查ansible hosts文件，脚本执行成功后会自动把[new-etcd]组的新节点移动至[etcd]组

```

- 注意：etcd 集群一次只能添加一个节点，如果你在[new-etcd]组中添加了2个新节点，那么需要执行两次 `ansible-playbook /etc/ansible/19.addetcd.yml`

### [可选]后续

上述步骤验证成功，确认新etcd集群工作正常后，可以重新配置运行apiserver，以让 k8s 集群能够识别新的etcd节点：

``` bash
# 重启 master 节点服务
$ ansible-playbook /etc/ansible/04.kube-master.yml -t restart_master

# 验证 k8s 能够识别新 etcd 节点
$ kubectl get cs
```

### 参考

- 官方文档 https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/runtime-configuration.md
