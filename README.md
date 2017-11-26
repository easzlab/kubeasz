# 利用Ansible部署kubernetes集群

![docker](./pics/docker.jpg) ![kube](./pics/kube.jpg) ![ansible](./pics/ansible.jpg)

本文档记录自己实践部署高可用k8s集群的过程，利用ansible-playbook简化二进制方式部署过程。

网上有很多类似shell脚本和ansible部署版本，要不看得太复杂，或者久未更新，所以这里自己造轮子吧。

二进制方式手动部署，将有助于理解系统各组件的交互原理和熟悉组建启动参数，进而能快速解决实际问题。

1. 建议阅读 [feisky.gitbooks](https://feisky.gitbooks.io/kubernetes/) 原理和部署章节。
1. 建议阅读 [opsnull教程](https://github.com/opsnull/follow-me-install-kubernetes-cluster) 二进制手工部署。

本文是按照上述文档，更新组件实践修饰而成，修改了部分安全特性。

## 特性

1. 截至2017-10-4 最新组件版本，参见[down版本](./down/download.sh) 文件。
1. 因本人部署节点IP属于同一网段，使用flannel新后端[host-gw](https://github.com/coreos/flannel/blob/master/Documentation/backends.md) 提升部分性能。

## 快速指南

以下为快速体验k8s集群的测试、开发环境--AllinOne部署，觉得比官方的minikube方便、简单很多。

### 1.准备一台虚机(推荐内存3G，CPU 2，硬盘 30G以上)，安装Ubuntu16.04，配置基础网络、更新源、SSH登陆等。
### 2.安装python2/git/python-pip/ansible
``` bash
# 更新
apt-get update && apt-get upgrade -y && apt-get dist-upgrade -y
# 删除不要的默认安装
apt-get purge ufw lxd lxd-client lxcfs lxc-common
# 安装依赖工具
apt-get install python2.7 git python-pip
# 安装ansible
pip install pip --upgrade
pip install ansible
# 国内加速
pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
# 配置ansible ssh密钥登陆
ssh-keygen -t rsa -b 2048 回车 回车 回车
ssh-copy-id $IP //$IP为本虚机地址
```
### 3.安装k8s
``` bash
git clone https://gitee.com/netmon/deploy-k8s-with-ansible.git
mv deploy-k8s-with-ansible/ /etc/ansible
cd /etc/ansible
# 配置ansible
cp example/ansible.cfg.example ansible.cfg
# 配置集群hosts
cp example/hosts.allinone.example hosts
然后根据实际情况修改此hosts文件
# 准备二进制安装包
按照down/download.sh文件提示先手工下载各种tar包到 ./down目录
sh down/download.sh
# 开始安装(一步安装)
ansible-playbook 90.setup.yml
# 或者采用分步安装
ansible-playbook 01.prepare.yml
ansible-playbook 02.etcd.yml
...
```
如果执行成功，k8s集群就安装好了

### 4.验证安装
``` bash
kubectl version
kubectl get componentstatus # 可以看到scheduler/controller-manager/etcd等组件 Healthy
kubectl clusterinfo # 可以看到kubernetes master(apiserver)组件 running
kubectl get node # 可以看到单 node Ready状态
kubectl get pod --all-namespaces # 可以查看所有集群pod状态
kubectl get svc --all-namespaces # 可以查看所有集群服务状态
```

## 多节点指南(文档更新中...)
1. 准备4台虚机(物理机也可，虚机实验更方便)，安装Ubuntu16.04(centos7理论上一样，不想ansible脚本太多条件判断)
1. 准备一台部署机(可以复用上述4台虚机)，安装ansible，配置到4台目标机器ssh无密码登陆等
1. 准备外部负载均衡，准备master节点的vip地址
1. 规划集群节点，完成ansible inventory文件[参考](hosts)
1. 其他安装步骤同单节点安装

