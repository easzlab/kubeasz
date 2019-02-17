# etcd 集群管理的 playbook

etcd 集群的主要操作包括`备份数据`,`添加/删除节点`等，本文介绍使用`ansible playbook`方便地完成这些任务。

- NOTE: 操作 etcd 集群节点增加/删除存在一定风险，请先在测试环境操作练手！

## 备份 etcd 数据

可以根据需要进行定期备份（使用 crontab），或者手动在任意正常 etcd 节点上执行备份：

``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```
- `kubeasz`项目也可以方便执行 `ansible-playbook /etc/ansible/23.backup.yml`，请阅读文档[备份恢复](cluster_restore.md)

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

## 删除 etcd 集群节点

删除节点的操作步骤比较简单，运行：`ansible-playbook /etc/ansible/tools/remove_etcd_node.yml`后按照提示输入待删除节点的IP地址即可。

主要删除步骤：

- 提示/获取用户输入待删除节点IP，并判断是否可以删除
- 获取待删除 etcd 节点的 ID 和 NAME 信息
- 修改 ansible hosts 文件，把待删除节点从 etcd 组中删除
- 执行 etcdctl member remove 命令删除节点
- 删除节点的 etcd 数据目录
- 重新配置启动整个 etcd 集群

## 重置 k8s 连接 etcd 参数

上述步骤验证成功，确认新etcd集群工作正常后，可以重新配置运行apiserver，以让 k8s 集群能够识别新的etcd集群：

``` bash
# 重启 master 节点服务
$ ansible-playbook /etc/ansible/04.kube-master.yml -t restart_master

# 验证 k8s 能够识别新 etcd 集群
$ kubectl get cs
```

## 参考

- 官方文档 https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/runtime-configuration.md
