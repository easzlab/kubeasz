# 公有云上部署 kubeasz

在公有云上使用`kubeasz`部署`k8s`集群需要注意以下几个常见问题。

### 安全组

注意虚机的安全组规则配置，一般集群内部节点之间端口全部放开即可；

### 网络组件

一般公有云对网络限制较多，跨节点 pod 通讯需要使用 OVERLAY 添加报头；默认配置详见example/config.yml

- flannel 使用 vxlan 模式：`FLANNEL_BACKEND: "vxlan"`
- calico 开启 ipinip：`CALICO_IPV4POOL_IPIP: "Always"`
- kube-router 开启 ipinip：`OVERLAY_TYPE: "full"`

### 节点公网访问

可以在安装时每个节点绑定`弹性公网地址`(EIP)，装完集群解绑；也可以开通NAT网关，或者利用iptables自建上网网关等方式

### 负载均衡

一般云厂商会限制使用`keepalived+haproxy`自建负载均衡，你可以根据云厂商文档使用云负载均衡（内网）四层TCP负载模式；

- kubeasz 2x 版本已无需依赖外部负载均衡实现apiserver的高可用，详见 [2x架构](https://github.com/easzlab/kubeasz/blob/dev2/docs/setup/00-planning_and_overall_intro.md#ha-architecture)
- kubeasz 1x 及以前版本需要负载均衡实现apiserver高可用，详见 [1x架构](https://github.com/easzlab/kubeasz/blob/dev1/docs/setup/00-planning_and_overall_intro.md#ha-architecture)

### 时间同步

一般云厂商提供的虚机都已默认安装时间同步服务，无需自行安装。 

### 访问 APISERVER

在公有云上安装完集群后，需要在公网访问集群 apiserver，而我们在安装前可能没有规划公网IP或者公网域名；而 apiserver 肯定需要 https 方式访问，在证书创建时需要加入公网ip/域名；可以参考这里[修改 APISERVER（MASTER）证书](../op/ch_apiserver_cert.md)

## 在公有云上部署多主高可用集群

处理好以上讨论的常见问题后，在公有云上使用 kubeasz 安装集群与自有环境没有差异。

- 使用 kubeasz 2x 版本安装单节点、单主多节点、多主多节点 k8s 集群，云上云下的预期安装体验完全一致
