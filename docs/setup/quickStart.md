## 快速指南

本文档适用于kubeasz 3.3.1以上版本，部署单节点集群(aio)，作为快速体验k8s集群的测试环境。

### 1.基础系统配置

- 准备一台虚机配置内存2G/硬盘30G以上
- 最小化安装`Ubuntu 16.04 server或者CentOS 7 Minimal`
- 配置基础网络、更新源、SSH登录等

**注意:** 确保在干净的系统上开始安装，不能使用曾经装过kubeadm或其他k8s发行版的环境

### 2.下载文件

- 下载工具脚本ezdown，举例使用kubeasz版本3.5.0

``` bash
export release=3.5.0
wget https://github.com/easzlab/kubeasz/releases/download/${release}/ezdown
chmod +x ./ezdown
```

- 使用工具脚本下载（更多关于ezdown的参数，运行./ezdown 查看）

下载kubeasz代码、二进制、默认容器镜像

``` bash
# 国内环境
./ezdown -D
# 海外环境
#./ezdown -D -m standard
```

【可选】下载额外容器镜像（cilium,flannel,prometheus等）

``` bash
./ezdown -X
```

【可选】下载离线系统包 (适用于无法使用yum/apt仓库情形)

``` bash
./ezdown -P
```

上述脚本运行成功后，所有文件（kubeasz代码、二进制、离线镜像）均已整理好放入目录`/etc/kubeasz`

- `/etc/kubeasz` 包含 kubeasz 版本为 ${release} 的发布代码
- `/etc/kubeasz/bin` 包含 k8s/etcd/docker/cni 等二进制文件
- `/etc/kubeasz/down` 包含集群安装时需要的离线容器镜像
- `/etc/kubeasz/down/packages` 包含集群安装时需要的系统基础软件

### 3.安装集群

- 容器化运行 kubeasz

```
./ezdown -S
```

- 使用默认配置安装 aio 集群

```
docker exec -it kubeasz ezctl start-aio
# 如果安装失败，查看日志排除后，使用如下命令重新安装aio集群
# docker exec -it kubeasz ezctl setup default all
```

### 4.验证安装

``` bash
$ source ~/.bashrc
$ kubectl version         # 验证集群版本     
$ kubectl get node        # 验证节点就绪 (Ready) 状态
$ kubectl get pod -A      # 验证集群pod状态，默认已安装网络插件、coredns、metrics-server等
$ kubectl get svc -A      # 验证集群服务状态
```

- 登录 `dashboard`可以查看和管理集群，更多内容请查阅[dashboard文档](../guide/dashboard.md)

### 5.清理

以上步骤创建的K8S开发测试环境请尽情折腾，碰到错误尽量通过查看日志、上网搜索、提交`issues`等方式解决；当然你也可以清理集群后重新创建。

在宿主机上，按照如下步骤清理

- 清理集群 `docker exec -it kubeasz ezctl destroy default`
- 重启节点，以确保清理残留的虚拟网卡、路由等信息
