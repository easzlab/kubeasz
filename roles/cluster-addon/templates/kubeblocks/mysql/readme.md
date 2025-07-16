# 说明

本目录资源涵盖使用kubeblocks创建mysql集群的各种操作，包括：集群变更、备份、恢复、监控等等；

## 双机主备集群创建流程

- 1.修改mysql componentdefinitions，支持修改属性

kubectl annotate componentdefinitions.apps.kubeblocks.io mysql-8.0-1.0.0 apps.kubeblocks.io/skip-immutable-check=true

- 2.修改mysql componentdefinitions，支持mysql容器使用hostNetwork（或者通过nodePort方式把主/备mysql服务都暴露出来）

kubectl edit componentdefinitions.apps.kubeblocks.io mysql-8.0-1.0.0

配置runtime.hostNetwork: true 和 runtime.dnsPolicy: ClusterFirstWithHostNet

- 3.正常创建mysql集群，并验证主备状态和机器上host端口3306

kubectl apply -n db -f 001.semisync-cluster.yaml

- 4.登录数据库主节点，写入测试数据

```
 CREATE DATABASE test;
 USE test;
 CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(255));
 INSERT INTO t1 VALUES (1, 'John Doe');
```

- 5.主节点关机；这样整个k8s集群无法访问，但是备节点上的mysql容器仍旧运行，登录备节点数据库，手动切主，提供读写服务

```
mysql> STOP REPLICA;
mysql> SET GLOBAL super_read_only = OFF;
mysql> SET GLOBAL read_only = OFF;

# 写入新数据
 USE test;
 INSERT INTO t1 VALUES (2, 'after master down');

# 主节点已故障排查，准备迎接主节点重启
mysql> SET GLOBAL super_read_only = ON;
mysql> SET GLOBAL read_only = ON;
```

- 6.主节点重启成功，验证mysql集群恢复，原备节点变主

登录主节点mysql，验证：

```
mysql> show status like 'Rpl%_status';
+------------------------------+-------+
| Variable_name                | Value |
+------------------------------+-------+
| Rpl_semi_sync_replica_status | OFF   |
| Rpl_semi_sync_source_status  | ON    |
+------------------------------+-------+
2 rows in set (0.00 sec)
```

登录备节点mysql，验证：

```
mysql> show status like 'Rpl%_status';
+------------------------------+-------+
| Variable_name                | Value |
+------------------------------+-------+
| Rpl_semi_sync_replica_status | ON    |
| Rpl_semi_sync_source_status  | OFF   |
+------------------------------+-------+
2 rows in set (0.00 sec)
```

### 参考

- https://kubeblocks.io/docs/release-1_0/kubeblocks-for-mysql/04-operations/11-rebuild-replica
- https://kubeblocks.io/docs/release-1_0/kubeblocks-for-mysql/03-topologies/01-semisync
