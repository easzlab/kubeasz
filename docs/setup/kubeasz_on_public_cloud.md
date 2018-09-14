# 公有云上部署 kubeasz

在公有云上使用`kubeasz`部署`k8s`集群需要注意以下几点：

1. 注意虚机的安全组规则配置，一般集群内部节点之间端口放开即可;

2. 部分`k8s`网络组件受限，一般可以选择 flannel (vxlan模式)、calico（开启ipinip）;

3. 无法自由创建`lb`节点，一般使用云负载均衡（内网）四层TCP负载模式;

4. 部分云厂商负载均衡使用四层负载模式时不支持添加进后端云服务器池的 ECS 既作为 Real Server，又作为客户端向所在的 SLB 实例发送请求；因此注意不要在 master节点执行 kubectl，会出现时通时不通的情况；

其他在公有云上的安装步骤与自有环境没有差异，节点规划可以参考 [example/hosts.cloud.example](../../example/hosts.cloud.example)

具体某个云厂商的问题，后续发现了会及时更新。
