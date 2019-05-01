### <center>腾讯云 CVM CentOS 系统 kubernetes 集群部署</center>

文档中脚本默认均以 root 用户执行

#### 高可用集群所需节点配置如下
---- 

|角色|数量|描述|
|:-|:-|:-|
|deploy节点|1|运行这份 ansible 脚本的节点|
|etcd节点|3|注意etcd集群必须是1,3,5,7...奇数个节点|
|master节点|3|共用etcd节点，master VIP(虚地址)在云管理后台创建，可根据需要提升机器配置或增加节点数|
|node节点|2|运行应用负载的节点，可根据需要提升机器配置或增加节点数|


#### 环境准备
---- 

##### 创建 CVM 实例
- 准备6台虚机，搭建一个多主高可用集群，`Node` 节点内存不低于 4GB
- 生产环境一个节点只担任一个角色
- 1个 `deploy` 节点 网段 10.0.0.3/21
- 3个 `master` 节点 网段 10.0.8.0/21 ，建议采用 SSD 类型磁盘
- 2个或以上 `node` 节点 网段 10.0.8.0/21

##### 创建 master vip
- 云负载均衡中创建传统型内网 CLB, 区域与 CVM 相同，取名为 `k8s-master-lb`，假定为 `10.0.8.12`
- 创建 TCP 类型监听器，前端监听 `8443` 端口，转发后端 `6443` 端口
- 绑定 `master` 节点到监听器

##### 创建 ingress vip （收费类型 可在集群创建成功后操作）
- 云负载均衡中创建应用型外网 CLB, 区域与 CVM 相同，取名为 `k8s-ingress-lb`
- 创建 TCP 类型监听器，前端监听 `23457` 端口，转发后端 `23457` 端口
- 创建 TCP 类型监听器，前端监听 `23456` 端口，转发后端 `23456` 端口
- 绑定 `node` 节点到监听器


#### 部署步骤
---- 


##### 0. 基础系统配置

+ 使用自定义系统镜像 `k8s-node` 安装系统
+ 配置基础网络、更新源、SSH登陆等
+ 腾讯云后台创建 CLB


##### 1. 以 `CentOS 7.x 64bit` 镜像初始化 CVM 实例安装 deploy 节点

- 更新本节点主机名
  ```bash
  hostnamectl set-hostname deploy
  # 重新登录
  ```

- 更新主机节点列表  
  编辑文件
  ```bash
  vi /etc/hosts
  ```

  删除云主机自动创建的回环主机名映射行，如下类似
  ```
  127.0.0.1 VM_0_15_centos VM_0_15_centos
  ::1 VM_0_15_centos VM_0_15_centos
  ```

  根据服务器配置添加主机列表，添加
  ```
  10.0.0.3 deploy
  10.0.8.2 master01
  10.0.8.3 master02
  10.0.8.4 master02
  10.0.8.10 node01
  10.0.8.11 node02
  ```


##### 2. 在deploy节点安装及准备ansible

pip 安装 ansible
``` bash
yum install python-pip -y

# pip安装ansible （腾讯云服务器自带加速）
pip install pip --upgrade
pip install ansible

# pip安装ansible(国内如果安装太慢可以直接用pip阿里云加速)
pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
```


##### 3. 资源准备

- 克隆源码
  ```bash
  git clone --depth=1 -b cvm https://github.com/waitingsong/kubeasz.git /etc/ansible
  ```

- 下载 k8s 二进制文件
从分享的[百度云链接](https://pan.baidu.com/s/1c4RFaA)，下载解压到 `/etc/ansible/bin` 目录
  ```bash
  # 以安装k8s v1.13.5为例
  tar -xvf k8s.1-13-5.tar.gz -C /etc/ansible/
  ```

- 生成离线 docker 镜像  
  ```bash
  wget https://raw.githubusercontent.com/waitingsong/blog/master/201904/assets/make_basic_images_bundle.sh
  wget https://raw.githubusercontent.com/waitingsong/blog/master/201904/assets/make_extra_images_bundle.sh
  wget https://raw.githubusercontent.com/waitingsong/blog/master/201904/assets/make_istio_images_bundle.sh
  chmod a+x make_basic_images_bundle.sh
  chmod a+x make_extra_images_bundle.sh
  chmod a+x make_istio_images_bundle.sh

  # 根据需要执行脚本进行下载并打包 xz 格式压缩时间比较长
  # 分别生成以下文件
  # /tmp/basic_images_kubeasz_1.0.tar.xz
  # /tmp/extra_images_kubeasz_1.0.tar.xz
  # /tmp/istio_images_bundle_1.1.3.tar.xz
  ./make_basic_images_bundle.sh dump
  ./make_extra_images_bundle.sh dump
  ./make_istio_images_bundle.sh dump
  ```

- 下载离线 docker 镜像  
将上一步生成的文件和脚本文件拷贝到 deploy 节点服务器相同目录下执行  
istio 安装见文档 [istio_install.md](./istio_install.md)
  ```bash
  ./make_basic_images_bundle.sh extract
  ./make_extra_images_bundle.sh extract
  ```

##### 4. 配置集群参数
```bash
cd /etc/ansible && cp example/hosts.cloud.example hosts
``` 

编辑此 hosts 文件
```bash
vi /etc/ansible/hosts
```

更新以下内容
```
# deploy 节点的地址
10.0.0.3 NTP_ENABLED=yes

[etcd]
10.0.8.2 NODE_NAME=etcd1
10.0.8.3 NODE_NAME=etcd2
10.0.8.4 NODE_NAME=etcd3

[kube-master]
10.0.8.2
10.0.8.3
10.0.8.4

[kube-node]
10.0.8.10
10.0.8.11

MASTER_IP="10.0.8.12"                # 即 master vip 负载均衡内网地址
```


##### 5. 编排k8s安装

如果你对集群安装流程不熟悉，请阅读项目首页 **安装步骤** 讲解后分步安装，并对 **每步都进行验证**  

验证 ansible 执行 正常能看到所有节点返回 SUCCESS
```bash
ansible all -m ping
``` 

执行安装
```bash
cd /etc/ansible
# 一步安装
ansible-playbook 90.setup.yml
# 分步安装
ansible-playbook 01.prepare.yml
ansible-playbook 02.etcd.yml
ansible-playbook 03.docker.yml
ansible-playbook 04.kube-master.yml
ansible-playbook 05.kube-node.yml
ansible-playbook 06.network.yml
ansible-playbook 07.cluster-addon.yml

# 把 k8s 集群 ca 证书加入本机信任列表
cp /etc/kubernetes/ssl/ca.pem /etc/pki/ca-trust/source/anchors/ && update-ca-trust
```

#### 查看集群状态
```bash
kubectl cluster-info
kubectl get cs
kubectl get node
kubectl get pod,svc --all-namespaces -o wide
kubectl top node
```

#### 资源
- [kubeasz](https://github.com/gjmzj/kubeasz)
- [镜像打包脚本](https://github.com/waitingsong/blog/tree/master/201904/assets)

origin by [waitingsong](https://github.com/waitingsong/blog/blob/master/201904/k8s_cvm_intro.md)
