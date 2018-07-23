# K8S 集群备份与恢复

虽然 K8S 集群可以配置成多主多节点的高可用的部署，还是有必要了解下集群的备份和容灾恢复能力；在高可用k8s集群中 etcd集群保存了整个集群的状态，因此这里的备份与恢复重点就是：

- 从运行的etcd集群备份数据到磁盘文件
- 从etcd备份文件恢复数据到运行的etcd集群，然后据此重建整个集群

## 前提

k8s 集群可能因为软硬件故障或者误操作出现了不可自愈的问题，这个时候需要考虑集群从备份中恢复重建；使用kubeasz项目创建的集群如需恢复前提如下：

- 集群正常状态下的etcd 备份文件（etcd V3数据）
- 创建集群时使用的 CA证书相关文件
- 创建集群时使用的 ansible hosts文件

## 备份与恢复手动操作说明

首先用kubeasz 搭建一个测试集群，部署几个测试deployment，验证集群各项正常后，进行一次备份：

- 1.在一个etcd节点上执行数据备份，把产生的备份文件`snapshot.db`复制到所有etcd集群节点

``` bash
$ mkdir -p /backup/k8s/ && cd /backup/k8s
$ ETCDCTL_API=3 etcdctl snapshot save snapshot.db
```

- 2.在deploy节点把 CA证书相关备份出来 

``` bash
$ mkdir -p /backup/k8s/ && cp /etc/kubernetes/ssl/ca* /backup/k8s/
```

- 3.在deploy节点清理集群，模拟集群完全崩溃

``` bash
$ ansible-playbook /etc/ansible/99.clean.yml
```

- 4.在deploy节点开始一步步重建集群

``` bash
# 恢复原集群的CA 证书相关
$ mkdir -p /etc/kubernetes/ssl/ && cp /backup/k8s/* /etc/kubernetes/ssl/

# 然后执行集群恢复步骤，安装至 kube-node完成阶段
$ cd /etc/ansible
$ ansible-playbook 01.prepare.yml
$ ansible-playbook 02.etcd.yml 
$ ansible-playbook 03.docker.yml
$ ansible-playbook 04.kube-master.yml
$ ansible-playbook 05.kube-node.yml

# 以上步骤验证正常后，停止etcd集群服务，并清空新etcd集群数据目录
$ ansible etcd -m service -a 'name=etcd state=stopped'
$ asnible etcd -m file -a 'name=/var/lib/etcd/member/ state=absent'
```

- 5.手动分别登陆每个etcd节点进行数据备份恢复，每个etcd都要如下操作

``` bash
# 参照本etcd节点/etc/systemd/system/etcd.service的服务文件，替换如下{{}}中变量后执行
$ cd /backup/k8s/
$ ETCDCTL_API=3 etcdctl snapshot restore snapshot.db \
  --name {{ NODE_NAME }} \
  --initial-cluster {{ ETCD_NODES }} \
  --initial-cluster-token etcd-cluster-0 \
  --initial-advertise-peer-urls https://{{ inventory_hostname }}:2380

# 以上执行完后，会生成{{ NODE_NAME }}.etcd的文件夹，将它里面的member 拷贝到etcd数据目录中
$ cp -r {{ NODE_NAME }}.etcd/member /var/lib/etcd/

$ systemctl restart etcd
```

- 6.在deploy节点执行网络重建

``` bash
$ ansible-playbook /etc/ansible/tools/change_k8s_network.yml
```

执行完之后，可以验证整个集群是否恢复正常，之前的测试应用部署是否全部恢复。

- 参考：https://github.com/coreos/etcd/blob/master/Documentation/op-guide/recovery.md

## 备份恢复自动脚本操作指南

- 一.集群备份

``` bash
$ ansible-playbook /etc/ansible/23.backup.yml
```

执行完毕可以在目录 `/etc/ansible/roles/cluster-backup/files`下检查备份情况，示例如下：

``` bash
roles/cluster-backup/files/
├── ca			# 集群CA 相关备份
│   ├── ca-config.json
│   ├── ca.csr
│   ├── ca-csr.json
│   ├── ca-key.pem
│   └── ca.pem
├── hosts		# ansible hosts备份
│   ├── hosts		# 最近的备份
│   └── hosts-201807231642
├── readme.md
└── snapshot		# etcd 数据备份
    ├── snapshot-201807231642.db
    └── snapshot.db	# 最近的备份
```

- 二.模拟集群故障

``` bash
$ ansible-playbook /etc/ansible/99.clean.yml
```

**注意** 为了模拟集群彻底崩溃，这里清理整个集群；实际操作中，在有备份前提下，也建议彻底清理集群后再尝试去恢复

- 三.集群恢复

可以在 `roles/cluster-restore/defaults/main.yml` 文件中配置需要恢复的 etcd备份版本，默认使用最近一次备份

``` bash
$ ansible-playbook /etc/ansible/24.restore.yml
$ ansible-playbook /etc/ansible/tools/change_k8s_network.yml
```

执行完成可以验证整个集群是否恢复如初！
