## 快速指南

以下为快速体验k8s集群的测试、开发环境--单节点部署(aio)，国内环境下比官方的minikube方便、简单很多。

### 1.基础系统配置

- 准备一台虚机配置内存2G/硬盘30G以上
- 最小化安装`Ubuntu 16.04 server`或者`CentOS 7 Minimal`
- 配置基础网络、更新源、SSH登录等

### 2.下载文件

``` bash
# 下载工具脚本easzup，举例使用kubeasz版本2.0.2
export release=2.0.2
curl -C- -fLO --retry 3 https://github.com/easzlab/kubeasz/releases/download/${release}/easzup
chmod +x ./easzup
# 使用工具脚本下载
./easzup -D
```

上述脚本运行成功后，所有文件（kubeasz代码、二进制、离线镜像）均已整理好放入目录`/etc/ansible`

- `/etc/ansible` 包含 kubeasz 版本为 ${release} 的发布代码
- `/etc/ansible/bin` 包含 k8s/etcd/docker/cni 等二进制文件
- `/etc/ansible/down` 包含集群安装时需要的离线容器镜像
- `/etc/ansible/down/packages` 包含集群安装时需要的系统基础软件

### 3.配置 ssh 免密登录

``` bash
ssh-keygen -t rsa -b 2048 -N '' -f ~/.ssh/id_rsa
ssh-copy-id $IP  # $IP 为所有节点地址包括自身，按照提示输入 yes 和 root 密码
```

### 4.安装集群

- 4.1 容器化运行 kubeasz，详见[文档](docker_kubeasz.md)

```
./easzup -S
```

- 4.2 使用默认配置安装 aio 集群

```
docker exec -it kubeasz easzctl start-aio
```

### 5.验证安装

如果提示kubectl: command not found，退出重新ssh登录一下，环境变量生效即可

``` bash
$ kubectl version                   # 验证集群版本     
$ kubectl get componentstatus       # 验证 scheduler/controller-manager/etcd等组件状态
$ kubectl get node                  # 验证节点就绪 (Ready) 状态
$ kubectl get pod --all-namespaces  # 验证集群pod状态，默认已安装网络插件、coredns、metrics-server等
$ kubectl get svc --all-namespaces  # 验证集群服务状态
```
- 登录 `dashboard`可以查看和管理集群，更多内容请查阅[dashboard文档](../guide/dashboard.md)

### 6.清理

以上步骤创建的K8S开发测试环境请尽情折腾，碰到错误尽量通过查看日志、上网搜索、提交`issues`等方式解决；当然你也可以清理集群后重新创建。

在宿主机上，按照如下步骤清理

- 1.清理集群 `docker exec -it kubeasz easzctl destroy` 或 `docker exec -it kubeasz ansible-playbook /etc/ansible/99.clean.yml`
- 2.清理管理节点
  - 清理运行的容器 `easzup -C`
  - 清理容器镜像 `docker system prune -a`
  - 停止docker服务 `systemctl stop docker`
  - 删除下载文件 `rm -rf /etc/ansible /etc/docker /opt/kube`
  - 删除docker文件
```
$ umount /var/run/docker/netns/default
$ umount /var/lib/docker/overlay
$ rm -rf /var/lib/docker /var/run/docker
```

上述清理脚本执行成功后，建议重启节点，以确保清理残留的虚拟网卡、路由等信息。
