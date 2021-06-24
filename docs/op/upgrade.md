## k8s 集群升级

集群升级存在一定风险，请谨慎操作。 

- 在当前支持k8s大版本基础上升级任意小版本，比如当前安装集群为1.19.0，你可以方便的升级到任何1.19.x版本

### 备份etcd数据 

- 升级前对 etcd数据做备份，在任意 etcd节点上执行：

``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```

`kubeasz`项目也可以如下方便执行备份（假设集群名为k8s-01）

- `ezctl backup k8s-01`，详情阅读文档[备份恢复](cluster_restore.md)

### 快速k8s版本升级

快速升级是指只升级`k8s`版本，比较常见如`Bug修复` `重要特性发布`时使用。

- 首先去官网release下载待升级的k8s版本，例如`https://dl.k8s.io/v1.11.5/kubernetes-server-linux-amd64.tar.gz`
- 解压下载的tar.gz文件，找到如下`kube*`开头的二进制，复制替换ansible控制端目录`/etc/ansible/bin`对应文件
  - kube-apiserver
  - kube-controller-manager
  - kubectl
  - kubelet
  - kube-proxy
  - kube-scheduler

- 在ansible控制端执行`ezctl upgrade k8s-01` 即可完成k8s 升级，不会中断业务应用


### 其他升级说明

其他升级是指升级k8s组件包括：`etcd版本` `docker版本`，一般不需要用到，不建议升级，以下仅作说明。

- 1.下载所有组件相关新的二进制解压并替换 `/etc/ansible/bin/` 目录下文件

- 2.升级 etcd: `ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml -t upgrade_etcd playbooks/02.etcd.yml`

- 3.升级 docker （建议使用k8s官方支持的docker稳定版本）
  - 如果可以接受短暂业务中断，执行 `ansible-playbook -t upgrade_docker 03.docker.yml`
  - 如果要求零中断升级，执行 `ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml -t download_docker playbooks/03.runtime.yml`，然后手动执行如下
    - 待升级节点，先应用`kubectl cordon`和`kubectl drain`命令迁移业务pod
    - 待升级节点执行 `systemctl restart docker`
    - 恢复节点可调度 `kubectl uncordon`
