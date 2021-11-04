# K8S 集群备份与恢复

虽然 K8S 集群可以配置成多主多节点的高可用的部署，还是有必要了解下集群的备份和容灾恢复能力；在高可用k8s集群中 etcd集群保存了整个集群的状态，因此这里的备份与恢复重点就是：

- 从运行的etcd集群备份数据到磁盘文件
- 从etcd备份文件恢复数据，从而使集群恢复到备份时状态

## 备份与恢复操作说明

- 1.首先搭建一个测试集群，部署几个测试deployment，验证集群各项正常后，进行一次备份(假设集群名为k8s-01)：

``` bash
$ ezctl backup k8s-01
# 或者如下手动执行ansible命令
# ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml  playbooks/94.backup.yml
```

执行完毕可以在部署主机的备份目录下检查备份情况，示例如下：

```
/etc/kubeasz/clusters/k8s-01/backup/
├── snapshot_202106201205.db
├── snapshot_202106211406.db
└── snapshot.db
```
其中，snapshot.db始终为最近一次备份文件

- 2.模拟误删除操作（略）

- 3.恢复集群及验证

可以在 `roles/cluster-restore/defaults/main.yml` 文件中配置需要恢复的 etcd备份版本（从上述备份目录中选取），默认使用最近一次备份；执行恢复后，需要一定时间等待 pod/svc 等资源恢复重建。

``` bash
$ ezctl restore k8s-01
# 或者如下手动执行ansible命令
# ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml  playbooks/95.restore.yml
```
如果集群主要组件（master/etcd/node）等出现不可恢复问题，可以尝试使用如下步骤 [清理]() --> [创建]() --> [恢复]()

``` bash
$ ezctl clean k8s-01
# 或者如下手动执行ansible命令
# ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml  playbooks/99.clean.yml
$ ezctl setup k8s-01 01
$ ezctl setup k8s-01 02
$ ezctl setup k8s-01 03
$ ezctl setup k8s-01 04
$ ezctl setup k8s-01 05
...
$ ezctl restore k8s-01
# ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml  playbooks/95.restore.yml
```

## 参考

- https://etcd.io/docs/v3.4/op-guide/recovery/
