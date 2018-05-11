## 升级注意事项

集群更新存在一定风险，请谨慎操作。 

- 项目分支1.8安装的集群目前只能进行小版本1.8.x的升级
- 项目分支1.9和master安装的集群可以任意小版本、大版本的升级，即1.9.x升级至1.10.x也可以

### 备份etcd数据 

### 升级步骤

+ 1.下载最新项目代码 `git pull origin master`
+ 2.下载新的二进制解压并覆盖 `/etc/ansible/bin/` 目录下文件
+ 3.更新集群 `ansible-playbook -t upgrade_k8s 22.upgrade.yml`  
注：上述步骤升级过程中不会中断集群已有业务，如果同时需要升级docker版本，可以在每个node节点手工重启docker服务（docker服务重启会中断业务，可以结合`kubectl cordon`和`kubectl drain`命令实现零中断升级）
