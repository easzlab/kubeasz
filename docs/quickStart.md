## 快速指南

以下为快速体验k8s集群的测试、开发环境--AllinOne部署，觉得比官方的minikube方便、简单很多。

### 1.准备一台虚机(推荐内存3G，硬盘20G以上)，最小化安装Ubuntu16.04 server，配置基础网络、更新源、SSH登陆等。
### 2.安装python2/git/python-pip/ansible
``` bash
# 文档中脚本默认均以root用户执行
apt-get update && apt-get upgrade -y && apt-get dist-upgrade -y
# 删除不要的默认安装
apt-get purge ufw lxd lxd-client lxcfs lxc-common
# 安装依赖工具
apt-get install python2.7 git python-pip
# Ubuntu16.04可能需要配置以下软连接
ln -s /usr/bin/python2.7 /usr/bin/python
# 安装ansible (国内如果安装太慢可以直接用pip阿里云加速)
#pip install pip --upgrade
#pip install ansible
pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
# 配置ansible ssh密钥登陆
ssh-keygen -t rsa -b 2048 回车 回车 回车
ssh-copy-id $IP #$IP为本虚机地址，按照提示输入yes 和root密码
```
### 3.安装kubernetes集群
``` bash
git clone https://github.com/gjmzj/kubeasz.git
mv kubeasz /etc/ansible
# 下载已打包好的binaries，并且解压缩到./bin目录
# 国内请从我分享的百度云链接下载 https://pan.baidu.com/s/1eSetFSA
# 如果你有合适网络环境也可以按照/down/download.sh自行从官网下载各种tar包到 ./down目录，并执行download.sh
tar zxvf k8s.184.tar.gz
mv bin/* /etc/ansible/bin
# 配置集群hosts
cd /etc/ansible
cp example/hosts.allinone.example hosts
然后根据实际情况修改此hosts文件，所有节点都是本虚机IP
# 采用分步安装或者一步安装
ansible-playbook 01.prepare.yml
ansible-playbook 02.etcd.yml
ansible-playbook 03.kubectl.yml
ansible-playbook 04.docker.yml
ansible-playbook 05.calico.yml
ansible-playbook 06.kube-master.yml
ansible-playbook 07.kube-node.yml
# ansible-playbook 90.setup.yml # 一步安装
```
如果执行成功，k8s集群就安装好了。

### 4.验证安装
``` bash
# 如果提示kubectl: command not found，退出重新ssh登陆一下，环境变量生效即可
kubectl version
kubectl get componentstatus # 可以看到scheduler/controller-manager/etcd等组件 Healthy
kubectl clusterinfo # 可以看到kubernetes master(apiserver)组件 running
kubectl get node # 可以看到单 node Ready状态
kubectl get pod --all-namespaces # 可以查看所有集群pod状态
kubectl get svc --all-namespaces # 可以查看所有集群服务状态
```
### 5.安装主要组件
``` bash
# 安装kubedns
kubectl create -f manifests/kubedns
# 安装heapster
kubectl create -f manifests/heapster
# 安装dashboard
kubectl create -f manifests/dashboard
```
