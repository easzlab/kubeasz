# 公有云上部署 kubeasz

在公有云上使用`kubeasz`部署`k8s`集群需要注意以下几点：

1. 注意虚机的安全组规则配置，一般集群内部节点之间端口放开即可;

2. 部分`k8s`网络组件受限，一般可以选择 flannel (vxlan模式)、calico（开启ipinip）;

3. 无法自由创建`lb`节点，一般使用云负载均衡（内网）四层TCP负载模式;

4. 部分云厂商负载均衡使用四层负载模式时不支持添加进后端云服务器池的 ECS 既作为 Real Server，又作为客户端向所在的 SLB 实例发送请求；因此注意不要在 master节点执行 kubectl，会出现时通时不通的情况；

## 在公有云上部署多主多节点集群

- 单节点和单主多节点集群的节点规划与自有环境没有差异

- 多主多节点集群节点规划不需要lb节点

其他在公有云上的安装步骤与自有环境没有差异，节点规划可以参考 [example/hosts.cloud.example](../../example/hosts.cloud.example)，如下：（避免deploy节点同时作为master节点）

``` bash
# 集群部署节点：一般为运行ansible 脚本的节点
# 变量 NTP_ENABLED (=yes/no) 设置集群是否安装 chrony 时间同步, 公有云上虚机不需要
[deploy]
10.1.0.160 NTP_ENABLED=no

# etcd集群请提供如下NODE_NAME，注意etcd集群必须是1,3,5,7...奇数个节点
[etcd]
10.1.0.160 NODE_NAME=etcd1
10.1.0.161 NODE_NAME=etcd2
10.1.0.162 NODE_NAME=etcd3

[kube-master]
10.1.0.161
10.1.0.162

# 公有云上一般都有提供负载均衡产品，且不允许自己创建，lb 节点留空，仅保留组名
[lb]

[kube-node]
10.1.0.160
10.1.0.163

# 参数 NEW_INSTALL：yes表示新建，no表示使用已有harbor服务器
[harbor]
#10.1.0.8 HARBOR_DOMAIN="harbor.yourdomain.com" NEW_INSTALL=no

...
```
+ 创建云负载均衡，例如阿里云slb如下：

``` bash
1. 首先创建SLB，注意选择【可用区】，【实例类型】可以先选‘私网’，【网络类型】专有网络，【虚拟交换机】跟你k8s集群节点同一交换机
2. 配置【协议&监听】TCP 【端口】8443，【后端服务器】即 master 节点服务器，端口 6443
3. 配置完成，记下负载均衡的内部地址（例如 10.1.0.200）
```
+ 继续配置 ansible hosts，设置`MASTER_IP` 为刚才创建的SLB地址 

``` bash
[all:vars]
# ---------集群主要参数---------------
#集群部署模式：allinone, single-master, multi-master
DEPLOY_MODE=multi-master

# 创建内网云负载均衡，然后配置：前端监听 tcp 8443，后端 tcp 6443，后端节点即 master 节点
MASTER_IP="10.1.0.200"          # 即负载均衡内网地址
KUBE_APISERVER="https://{{ MASTER_IP }}:8443"

# 集群网络插件，目前支持calico, flannel
CLUSTER_NETWORK="flannel"

... 
```
+ 一步创建集群 `ansible-playbook /etc/ansible/90.setup.yml`

### 其他资料

另外由[li-sen](https://github.com/li-sen)分享的[kubeasz-阿里云vpc部署记录](https://li-sen.github.io/post/blog-wiki/2018-09-27-k8s-kubeasz-%E9%98%BF%E9%87%8C%E4%BA%91vpc%E9%83%A8%E7%BD%B2%E8%AE%B0%E5%BD%95/)：介绍了阿里云上自建高可用k8s集群碰过的问题与解决，主要是使用一台haproxy中转解决slb的限制问题。

