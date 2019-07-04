# 离线安装集群

kubeasz 2.0.1 开始支持完全离线安装，目前已测试 Ubuntu16.04/Ubuntu18.04/CentOS7 系统。

## 离线文件准备

在一台能够访问互联网的服务器上执行：

```
$ curl -C- -fLO --retry 3 https://github.com/easzlab/kubeasz/releases/download/2.0.1/easzup
$ chmod +x ./easzup
# 使用工具脚本下载
$ ./easzup -D
```

执行成功后，所有文件均已整理好放入目录`/etc/ansilbe`，只要把该目录整体复制到任何离线的机器上，即可开始安装集群，离线文件包括：

- kubeasz 项目代码 --> /etc/ansible
- kubernetes 集群组件二进制 --> /etc/ansible/bin
- 其他集群组件二进制（etcd/CNI等）--> /etc/ansible/bin
- 操作系统基础依赖软件包（haproxy/ipvsadm/ipset/socat等）--> /etc/ansible/down/packages
- 集群基本插件镜像（coredns/dashboard/metrics-server等）--> /etc/ansible/down

离线文件不包括：

- 管理端 ansible 安装，但可以使用 kubeasz 容器运行 ansible 脚本
- 其他更多 kubernetes 插件镜像

## 离线安装

上述下载完成后，执行

``` bash
# 启动 kubeasz 容器
$ ./easzup -S

# 进入容器
$ docker exec -it kubeasz sh

# 设置参数启用离线安装
$ cd /etc/ansible
$ sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' roles/chrony/defaults/main.yml
$ sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' roles/ex-lb/defaults/main.yml
$ sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' roles/kube-node/defaults/main.yml
$ sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' roles/prepare/defaults/main.yml

# 按照文档 https://github.com/easzlab/kubeasz/blob/master/docs/setup/00-planning_and_overall_intro.md 集群规划后安装
$ ansible-playbook 90.setup.yml
```
