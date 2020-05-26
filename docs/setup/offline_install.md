# 离线安装集群

kubeasz 2.0.1 开始支持**完全离线安装**，目前已测试 `Ubuntu1604|1804` `CentOS7` `Debian9|10` 系统。

## 离线文件准备

在一台能够访问互联网的服务器上执行：

- 下载工具脚本easzup，举例使用kubeasz版本2.3.0

``` bash
export release=2.3.0
curl -C- -fLO --retry 3 https://github.com/easzlab/kubeasz/releases/download/${release}/easzup
chmod +x ./easzup
```

- 使用工具脚本下载

默认下载最新推荐k8s/docker等版本，使用命令`./easzup` 查看工具脚本的帮助信息

``` bash
# 举例使用 k8s 版本 v1.18.2，docker 19.03.5
./easzup -D -d 19.03.5 -k v1.18.2
# 下载离线系统软件包
./easzup -P
```

执行成功后，所有文件均已整理好放入目录`/etc/ansible`，只要把该目录整体复制到任何离线的机器上，即可开始安装集群，离线文件包括：

- `/etc/ansible` 包含 kubeasz 版本为 ${release} 的发布代码
- `/etc/ansible/bin` 包含 k8s/etcd/docker/cni 等二进制文件
- `/etc/ansible/down` 包含集群安装时需要的离线容器镜像
- `/etc/ansible/down/packages` 包含集群安装时需要的系统基础软件

离线文件不包括：

- 管理端 ansible 安装，但可以使用 kubeasz 容器运行 ansible 脚本
- 其他更多 kubernetes 插件镜像

## 离线安装

上述下载完成后，把`/etc/ansible`整个目录复制到目标离线服务器相同目录，然后在离线服务器上运行：

- 离线安装 docker，检查本地文件，正常会提示所有文件已经下载完成

```
./easzup -D
```

- 启动 kubeasz 容器

```
./easzup -S
```

- 设置参数允许离线安装

```
sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' /etc/ansible/roles/chrony/defaults/main.yml
sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' /etc/ansible/roles/ex-lb/defaults/main.yml
sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' /etc/ansible/roles/kube-node/defaults/main.yml
sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' /etc/ansible/roles/prepare/defaults/main.yml
```

- 举例安装单节点集群，参考 https://github.com/easzlab/kubeasz/blob/master/docs/setup/quickStart.md

```
docker exec -it kubeasz easzctl start-aio
```

- 多节点集群，进入kubeasz 容器内 `kubectl exec -it kubeasz bash`，参考https://github.com/easzlab/kubeasz/blob/master/docs/setup/00-planning_and_overall_intro.md 进行集群规划和设置后安装

```
#ansible-playbook 90.setup.yml
```
