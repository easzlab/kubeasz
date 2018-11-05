## 升级注意事项

集群更新存在一定风险，请谨慎操作。 

- 项目分支`master`安装的集群可以在k8s 1.8/1.9/1.10/1.11/1.12 任意小版本、大版本间升级（特别注意如果跨大版本升级需要修改/etc/ansible/hosts文件中的参数K8S_VER）
- 项目分支`closed`（已停止更新）安装的集群目前只能进行小版本1.8.x的升级

### 备份etcd数据 

- 升级前手动对 etcd数据做镜像备份，在任意 etcd节点上执行：

``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```
- `kubeasz`项目也可以方便执行 `ansible-playbook /etc/ansible/23.backup.yml`，详情阅读文档[备份恢复](cluster_restore.md)

### 升级步骤

- 1.下载新的二进制解压并替换 `/etc/ansible/bin/` 目录下文件

- 2a.如果不需要升级 docker版本：执行 `ansible-playbook -t upgrade_k8s 22.upgrade.yml` 即可完成k8s 升级，不会中断业务应用
  - 注：建议使用稳定版本 docker

- 2b.如果可以接受短暂业务中断，执行 `ansible-playbook -t upgrade_k8s,upgrade_docker 22.upgrade.yml` 即可升级 k8s和 docker(如果有新的docker二进制)

- 2c.如果要求零中断升级 k8s和 docker
  - i 执行 `ansible-playbook -t upgrade_k8s,download_docker 22.upgrade.yml` (该步骤不会影响k8s上的业务应用)
  - ii 逐个升级重启每个node节点的dockerd服务
    - 待重启节点，先应用`kubectl cordon`和`kubectl drain`命令迁移业务pod
    - 待重启节点执行 `systemctl restart docker`
    - 恢复节点可调度 `kubectl uncordon`
