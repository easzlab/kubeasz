# 管理 etcd 集群

Etcd 集群支持在线改变集群成员节点，可以增加、修改、删除成员节点；不过改变成员数量仍旧需要满足集群成员多数同意原则（quorum），另外请记住集群成员数量变化的影响：

- 注意：如果etcd 集群有故障节点，务必先删除故障节点，然后添加新节点，[参考FAQ](https://etcd.io/docs/v3.4.0/faq/)
- 增加 etcd 集群节点, 提高集群稳定性
- 增加 etcd 集群节点, 提高集群读性能（所有节点数据一致，客户端可以从任意节点读取数据）
- 增加 etcd 集群节点, 降低集群写性能（所有节点数据一致，每一次写入会需要所有节点数据同步）

## 备份 etcd 数据

可以根据需要进行定期备份（使用 crontab），或者手动在任意正常 etcd 节点上执行备份：

``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```

## etcd 集群节点操作

首先确认配置 ssh 免密码登录，然后执行 (假设待操作节点为 192.168.1.11，集群名称test-k8s)：

- 增加 etcd 节点：`$ ezctl add-etcd test-k8s 192.168.1.11`
- 删除 etcd 节点：`$ ezctl del-etcd test-k8s 192.168.1.11`

具体操作流程参考 ezctl中 add-etcd/del-etcd 相关函数和playbooks/ 目录的操作剧本

### 验证 etcd 集群

``` bash
# 登录任意etcd节点验证etcd集群状态
$ export ETCDCTL_API=3 
$ etcdctl member list

# 验证所有etcd节点服务状态和日志
$ systemctl status etcd
$ journalctl -u etcd -f
```

## 参考

- 官方文档 https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/runtime-configuration.md
