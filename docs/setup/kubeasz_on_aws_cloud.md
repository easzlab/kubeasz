### <center>AWS EC2 Amazon Linux 系统 kubernetes 集群高可用部署</center>

文档中脚本默认均以 root 用户执行  
Amazon Linux默认不能以root用户登录，只能先登录到`ec2-user`，再通过`sudo su - root`切换root用户


#### 高可用集群所需节点配置如下
---- 

|角色|数量|IP|描述|
|:-|:-|:-|:-|
|deploy节点|1|172.31.16.81|运行这份 ansible 脚本的节点|
|etcd节点|3|172.31.9.100, 172.31.9.192, 172.31.11.185|注意etcd集群必须是1,3,5,7...奇数个节点|
|master节点|2|172.31.9.100, 172.31.9.192|共用etcd节点，master VIP(虚地址)采用内部ELB域名代替，可根据需要提升机器配置或增加节点数|
|node节点|2|172.31.11.185, 172.31.14.4|运行应用负载的节点，可根据需要提升机器配置或增加节点数|


#### 环境准备
---- 

##### 创建 EC2 实例
- 准备5台虚机，搭建一个多主高可用集群，`Node` 节点内存不低于 4GB
- 生产环境一个节点只担任一个角色
- 1个 `deploy` 节点 网段 172.31.16.0/20
- 2个 `master` 节点 网段 172.31.0.0/20 ，建议采用 SSD 类型磁盘
- 2个 `node` 节点 网段 172.31.0.0/20
- 2个弹性IP、1个NAT网关，其中1个弹性IP绑定到`deploy`节点、另1个绑定到`NAT网关`
- `master`节点和`node`节点不分配公网IP，仅通过`NAT`连接外网

**注意：** 
- master和node所在的子网需添加路由表到NAT网关, 否则无法连接到外网
- NAT子网需与deploy节点相同，因为最终都是通过同一个Internet网关访问外网，在主路由表


##### 创建 master ELB
- 云负载均衡中创建 **经典型内网ELB**, 区域与 EC2 相同，取名为 `k8s-master-lb`，假定内网域名为 `internal-k8s-master-lb-42488333.xxxx.elb.amazonaws.com`
- 创建 TCP 类型监听器，前端监听 `8443` 端口，转发后端 `6443` 端口
- 绑定 `master` 节点到监听器，子网选择master和node节点相同

##### 创建 ingress ELB （收费类型 可在集群创建成功后操作）
- 云负载均衡中创建 **经典型外网ELB**, 区域与 EC2 相同，取名为 `k8s-ingress-lb`
- 创建 TCP 类型监听器，前端监听 `80` 端口，转发后端 `23456` 端口
- 绑定 `master和node` 节点到监听器(这里也可以只负载到master或node)
- 运行状态检查，Ping协议: `TCP`，Ping端口: `23456`
- 如果想要访问`traefik dashboard`，需要再创建一个 TCP 类型监听器，前端后端都监听`traefik admin暴露的nodePort端口`

##### 开启安全组
- 集群节点间的安全组要所有协议都开放

#### 部署步骤
---- 

##### 0. 基础系统配置

+ 使用社区AMI amzn2-ami 安装系统
+ 配置SSH root免密登陆等
+ 下载kubernetes对应版本bin文件

##### 1. 在deploy节点安装及准备ansible

pip 安装 ansible
``` bash
yum install python-pip -y

# pip安装ansible
pip install pip --upgrade
pip install ansible

# pip安装ansible(国内如果安装太慢可以直接用pip阿里云加速)
pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
```


##### 2. 资源准备

- 克隆源码
  ```bash
  git clone --depth=1 https://github.com/easzlab/kubeasz.git /etc/ansible
  ```

- 下载 k8s 二进制文件
从分享的[百度云链接](https://pan.baidu.com/s/1c4RFaA)，下载解压到 `/etc/ansible/bin` 目录
  ```bash
  # 以安装k8s v1-14-1为例
  tar -xvf k8s.1-14-1.tar.gz -C /etc/ansible/
  ```

##### 3. 配置集群参数
```bash
cd /etc/ansible && cp example/hosts.cloud.example hosts
``` 

编辑此 hosts 文件
```bash
vi /etc/ansible/hosts
```

更新以下内容
```
[deploy]
172.31.16.81 NTP_ENABLED=yes

# etcd集群请提供如下NODE_NAME，注意etcd集群必须是1,3,5,7...奇数个节点
[etcd]
172.31.9.100 NODE_NAME=etcd1
172.31.9.192 NODE_NAME=etcd2
172.31.11.185 NODE_NAME=etcd3

[kube-master]
172.31.9.100
172.31.9.192

[kube-node]
172.31.11.185
172.31.14.4

MASTER_IP="internal-k8s-master-lb-42488333.xxxx.elb.amazonaws.com"     # 即 master vip 负载均衡内网地址
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

#### ingress访问
可访问`外部ELB`自动分配的域名地址，正式环境可以将域名做个`CNAME`到该地址。  
测试示例：
``` yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:alpine

---
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  ports:
  - name: nginx-port
    port: 80
  selector:
    app: nginx

---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: nginx
spec:
  rules:
  - host: aws.xxxx.cn
    http:
      paths:
      - path: /
        backend:
          serviceName: nginx
          servicePort: nginx-port
```
