## 快速指南

以下为快速体验k8s集群的测试、开发环境--allinone部署，国内环境下觉得比官方的minikube方便、简单很多。

### 1.基础系统配置

+ 推荐内存2G/硬盘30G以上
+ 最小化安装`Ubuntu 16.04 server`或者`CentOS 7 Minimal`
+ 配置基础网络、更新源、SSH登陆等

### 2.安装依赖工具

Ubuntu 16.04 请执行以下脚本:

``` bash
# 文档中脚本默认均以root用户执行
# 安装依赖工具
apt-get install python2.7 git python-pip
# Ubuntu16.04可能需要配置以下软连接
ln -s /usr/bin/python2.7 /usr/bin/python
```
CentOS 7 请执行以下脚本：

``` bash
# 文档中脚本默认均以root用户执行
# 安装 epel 源
yum install epel-release -y
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
将`/usr/bin/pip`做以下修改：

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

- 4.1 下载项目源码

``` bash
# 方式一：使用git clone
git clone --depth=1 https://github.com/gjmzj/kubeasz.git /etc/ansible

# 方式二：从发布页面 https://github.com/gjmzj/kubeasz/releases 下载源码解压到同样目录
```
- 4.2a 下载二进制文件  
请从分享的[百度云链接](https://pan.baidu.com/s/1c4RFaA)，下载解压到/etc/ansible/bin目录，如果你有合适网络环境也可以按照/down/download.sh自行从官网下载各种tar包

``` bash
tar xvf k8s.1-9-8.tar.gz	# 以安装k8s v1.9.8为例
mv bin/* /etc/ansible/bin
```
- 4.2b [可选]下载离线docker镜像  
服务器使用内部yum源/apt源，但是无法访问公网情况下，请下载离线docker镜像完成集群安装；从百度云盘把`basic_images_kubeasz_x.y.tar.gz` 下载解压到`/etc/ansible/down` 目录

``` bash
tar xvf basic_images_kubeasz_1.0.tar.gz -C /etc/ansible/down
```
- 4.3 配置集群参数
  - 4.3.1 必要配置：`cd /etc/ansible && cp example/hosts.allinone.example hosts`, 然后实际情况修改此hosts文件
  - 4.3.2 可选配置，初次使用可以不做修改，详见[配置指南](config_guide.md)
  - 4.3.3 验证ansible 安装：`ansible all -m ping` 正常能看到节点返回 SUCCESS

- 4.4 开始安装
如果你对集群安装流程不熟悉，请阅读项目首页 **安装步骤** 讲解后分步安装，并对 **每步都进行验证**

``` bash
# 分步安装
ansible-playbook 01.prepare.yml
ansible-playbook 02.etcd.yml
ansible-playbook 03.docker.yml
ansible-playbook 04.kube-master.yml
ansible-playbook 05.kube-node.yml
ansible-playbook 06.network.yml
ansible-playbook 07.cluster-addon.yml 
# 一步安装
#ansible-playbook 90.setup.yml
```

+ [可选]对集群节点进行操作系统层面的安全加固 `ansible-playbook roles/os-harden/os-harden.yml`，详情请参考[os-harden项目](https://github.com/dev-sec/ansible-os-hardening)

### 5.验证安装
如果提示kubectl: command not found，退出重新ssh登陆一下，环境变量生效即可

``` bash
kubectl version
kubectl get componentstatus # 可以看到scheduler/controller-manager/etcd等组件 Healthy
kubectl cluster-info # 可以看到kubernetes master(apiserver)组件 running
kubectl get node # 可以看到单 node Ready状态
kubectl get pod --all-namespaces # 可以查看所有集群pod状态，默认已安装网络插件、coredns、metrics-server等
kubectl get svc --all-namespaces # 可以查看所有集群服务状态
```
### 6.安装主要组件

``` bash
# 安装kubedns，默认已集成安装
#kubectl create -f /etc/ansible/manifests/kubedns
# 安装dashboard，默认已集成安装
#kubectl create -f /etc/ansible/manifests/dashboard
```
+ 登陆 `dashboard`可以查看和管理集群，更多内容请查阅[dashboard文档](../guide/dashboard.md)

### 7.清理集群

以上步骤创建的K8S开发测试环境请尽情折腾，碰到错误尽量通过查看日志、上网搜索、提交`issues`等方式解决；当然如果是彻底奔溃了，可以清理集群后重新创建。

``` bash
ansible-playbook 99.clean.yml
```

如果出现清理失败，类似报错：`... Device or resource busy: '/var/run/docker/netns/xxxxxxxxxx'`，需要手动umount该目录后清理

``` bash
$ umount /var/run/docker/netns/xxxxxxxxxx
$ rm -rf /var/run/docker/netns/xxxxxxxxxx
```
