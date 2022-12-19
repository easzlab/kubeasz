## k8s 集群升级

集群升级存在一定风险，请谨慎操作。 

- 支持k8s相同大版本基础上升级任意小版本，比如当前安装集群为1.25.0，你可以方便的升级到任何1.25.x版本
- 不建议跨大版本升级，一般大版本更新时k8s api有一些变动

### 备份etcd数据

- 自动备份

`kubeasz`项目也可以如下方便执行备份（假设集群名为k8s-01），详情阅读文档[备份恢复](cluster_restore.md)

```
dk ezctl backup k8s-01
```

- 手动备份 etcd数据，在任意 etcd节点上执行：

``` bash
# snapshot备份
$ ETCDCTL_API=3 etcdctl snapshot save backup.db
# 查看备份
$ ETCDCTL_API=3 etcdctl --write-out=table snapshot status backup.db
```

### k8s 升级小版本

快速升级`k8s`小版本，比较常见如`Bug修复` `特性发布`时使用。

- 首先去官网release下载待升级的k8s版本，例如`https://dl.k8s.io/v1.25.4/kubernetes-server-linux-amd64.tar.gz`
- 解压下载的tar.gz文件，找到如下`kube*`开头的二进制，复制替换kubeasz控制端目录`/etc/kubeasz/bin`对应文件
  - kube-apiserver
  - kube-controller-manager
  - kubectl
  - kubelet
  - kube-proxy
  - kube-scheduler

- 在kubeasz控制端执行`dk ezctl upgrade k8s-01` 即可完成k8s 升级，不会中断业务应用


### 其他升级说明

其他升级是指升级k8s组件包括：`etcd版本` `docker版本`，一般不需要用到，不建议升级，以下仅作说明。

- 1.下载所有组件相关新的二进制解压并替换 `/etc/kubeasz/bin/` 目录下文件

- 2.升级 etcd: `ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml -t upgrade_etcd playbooks/02.etcd.yml`

- 3.升级 docker （建议使用k8s官方支持的docker稳定版本）
  - 如果可以接受短暂业务中断，执行 `ansible-playbook -t upgrade_docker 03.docker.yml`
  - 如果要求零中断升级，执行 `ansible-playbook -i clusters/k8s-01/hosts -e @clusters/k8s-01/config.yml -t download_docker playbooks/03.runtime.yml`，然后手动执行如下
    - 待升级节点，先应用`kubectl cordon`和`kubectl drain`命令迁移业务pod
    - 待升级节点执行 `systemctl restart docker`
    - 恢复节点可调度 `kubectl uncordon`
