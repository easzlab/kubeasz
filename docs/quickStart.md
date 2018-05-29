## 快速指南

以下为快速体验k8s集群的测试、开发环境--allinone部署，国内环境下觉得比官方的minikube方便、简单很多。

### 1.基础系统配置

+ 推荐内存2G/硬盘20G以上
+ 最小化安装`Ubuntu 16.04 server`或者`CentOS 7 Minimal`
+ 配置基础网络、更新源、SSH登陆等

### 2.安装依赖工具

Ubuntu 16.04 请执行以下脚本:

``` bash
# 文档中脚本默认均以root用户执行
apt-get update && apt-get upgrade -y && apt-get dist-upgrade -y
# 安装依赖工具
apt-get install python2.7 git python-pip
# Ubuntu16.04可能需要配置以下软连接
ln -s /usr/bin/python2.7 /usr/bin/python
```
CentOS 7 请执行以下脚本：

``` bash
# 文档中脚本默认均以root用户执行
# 安装 epel 源并更新
yum install epel-release -y
yum update
# 安装依赖工具
yum install git python python-pip -y
```
### 3.ansible安装及准备

``` bash
# 安装ansible (国内如果安装太慢可以直接用pip阿里云加速)
#pip install pip --upgrade
#pip install ansible
pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
# 配置ansible ssh密钥登陆
ssh-keygen -t rsa -b 2048 回车 回车 回车
ssh-copy-id $IP #$IP为本虚机地址，按照提示输入yes 和root密码
```

在`Ubuntu 16.04`中，如果出现以下错误:

``` bash
Traceback (most recent call last):
  File "/usr/bin/pip", line 9, in <module>
    from pip import main
ImportError: cannot import name main
```
将`/usr/bin/python`做以下修改：

``` bash
#原代码
from pip import main
if __name__ == '__main__':
    sys.exit(main())

#修改后
from pip import __main__
if __name__ == '__main__':
    sys.exit(__main__._main())
```


### 4.安装kubernetes集群
``` bash
git clone https://github.com/gjmzj/kubeasz.git
mkdir -p /etc/ansible
mv kubeasz/* /etc/ansible
# 下载已打包好的binaries，解压到/etc/ansible/bin目录
# 国内请从分享的百度云链接下载 https://pan.baidu.com/s/1c4RFaA
# 如果你有合适网络环境也可以按照/down/download.sh自行从官网下载各种tar包到 ./down目录，并执行download.sh
tar zxvf k8s.193.tar.gz
mv bin/* /etc/ansible/bin
# 配置ansible的hosts文件
cd /etc/ansible
cp example/hosts.allinone.example hosts # 然后根据实际情况修改此hosts文件，所有节点改成本虚机IP
# 开始集群安装，如果你对集群安装流程不熟悉，请阅读分步安装讲解后一步一步安装，并对每步都进行验证
# 分步安装
ansible-playbook 01.prepare.yml
ansible-playbook 02.etcd.yml
ansible-playbook 03.docker.yml
ansible-playbook 04.kube-master.yml
ansible-playbook 05.kube-node.yml
ansible-playbook 06.network.yml
# 一步安装
#ansible-playbook 90.setup.yml
```
如果执行成功，k8s集群就安装好了。详细分步讲解请查看项目目录 `/docs` 下相关文档

+ [可选]对节点进行操作系统层面的安全加固 `ansible-playbook roles/os-harden/os-harden.yml`，详情请参考[os-harden项目](https://github.com/dev-sec/ansible-os-hardening)

### 5.验证安装
``` bash
# 如果提示kubectl: command not found，退出重新ssh登陆一下，环境变量生效即可
kubectl version
kubectl get componentstatus # 可以看到scheduler/controller-manager/etcd等组件 Healthy
kubectl cluster-info # 可以看到kubernetes master(apiserver)组件 running
kubectl get node # 可以看到单 node Ready状态
kubectl get pod --all-namespaces # 可以查看所有集群pod状态
kubectl get svc --all-namespaces # 可以查看所有集群服务状态
```
### 6.安装主要组件
``` bash
# 安装kubedns
kubectl create -f /etc/ansible/manifests/kubedns
# 安装heapster
kubectl create -f /etc/ansible/manifests/heapster
# 安装dashboard
kubectl create -f /etc/ansible/manifests/dashboard
```
+ 登陆 `dashboard`可以查看和管理集群，更多内容请查阅[dashboard文档](guide/dashboard.md)

### 7.清理集群

以上步骤创建的K8S开发测试环境请尽情折腾，碰到错误尽量通过查看日志、上网搜索、提交`issues`等方式解决；当然如果是彻底奔溃了，可以清理集群后重新创建。

一步清理：`ansible-playbook 99.clean.yml`
