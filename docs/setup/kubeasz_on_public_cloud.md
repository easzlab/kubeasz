# 公有云上部署 kubeasz

在公有云上使用`kubeasz`部署`k8s`集群需要注意以下几点：

1. 注意虚机的安全组规则配置，一般集群内部节点之间端口放开即可;

2. 部分`k8s`网络组件受限，一般可以选择 flannel (vxlan模式)、calico（开启ipinip）;

3. 无法自由创建`lb`节点，一般使用云负载均衡（内网）四层TCP负载模式;

4. 安装时各个节点开通公网访问（可以使用绑定EIP/开通NAT网关/利用iptables自建上网网关等方式）

5. 部分云厂商负载均衡使用四层模式时不支持添加进后端云服务器池的 ECS 既作为 Real Server，又作为客户端向所在的 SLB 实例发送请求；因此在 master节点执行 kubectl 访问 apiserver VIP 地址时，会出现时通时不通的情况；但是不影响 kubeasz-k8s 集群正常工作。

## 在公有云上部署多主高可用集群

- 单节点、单主多节点集群的规划及安装与自有环境没有差异

- 多主高可用集群不需要lb，节点规划可以参考 [example/hosts.cloud.example](../../example/hosts.cloud.example)

- [阿里云部署 kubeasz 举例](kubeasz_on_aliyun.md)

- [腾讯云部署 kubeasz 举例](kubeasz_on_tencent_cloud.md)

- [百度云部署 kubeasz 举例](kubeasz_on_baidu_cloud.md)

- [AWS 部署 kubeasz 举例](kubeasz_on_aws_cloud.md)

