## 升级注意事项

集群更新存在一定风险，请谨慎操作。 

- 项目分支`master`安装的集群可以在k8s 1.8/1.9/1.10任意小版本、大版本间升级（特别注意如果跨大版本升级需要修改/etc/ansible/hosts文件中的参数K8S_VER）
- 项目分支`closed`（已停止更新）安装的集群目前只能进行小版本1.8.x的升级

### 备份etcd数据 

- 升级前对 etcd数据做镜像备份  
``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```
- 从备份恢复可以参考：[备份恢复](cluster_restore.md)

### 升级步骤

- 1.下载最新项目代码 `git pull origin master`（注意手动更新现有hosts的配置项与example中的实例一致）
- 2.下载新的二进制解压并替换 `/etc/ansible/bin/` 目录下文件
- 3a.如果可以接受短暂业务中断，执行 `ansible-playbook -t upgrade_k8s,restart_dockerd 22.upgrade.yml` 即可
- 3b.如果要求零中断升级集群
  - 首先执行 `ansible-playbook -t upgrade_k8s 22.upgrade.yml` (该步骤不会影响k8s上的业务应用)
  - 然后逐个升级重启每个node节点的dockerd服务
    - 待重启节点，先应用`kubectl cordon`和`kubectl drain`命令
    - 待重启节点执行 `systemctl restart docker`
    - 恢复待重启节点可调度 `kubectl uncordon`
