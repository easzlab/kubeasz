## 快速指南

以下为快速体验k8s集群的测试、开发环境--AllinOne部署，觉得比官方的minikube方便、简单很多。

### 1.准备一台虚机(推荐内存3G，硬盘20G以上)，安装Ubuntu16.04，配置基础网络、更新源、SSH登陆等。
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

