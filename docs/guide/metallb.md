# metallb 网络负载均衡

本文档已过期，以下内容仅做介绍，安装请参考最新官方文档

`Metallb`是在自有硬件上（非公有云）实现 `Kubernetes Load-balancer`的工具，由`google`团队开源，值得推荐！项目[github主页](https://github.com/google/metallb)。

## metallb 简介

这里简单介绍下它的实现原理，具体可以参考[metallb官网](https://metallb.universe.tf/)，文档非常简洁、清晰。目前有如下的使用限制：

- `Kubernetes v1.9.0`版本以上，暂不支持`ipvs`模式
- 支持网络组件 (flannel/weave/romana), calico 部分支持
- `layer2`和`bgp`两种模式，其中`bgp`模式需要外部网络设备支持`bgp`协议

`metallb`主要实现了两个功能：地址分配和对外宣告

- 地址分配：需要向网络管理员申请一段ip地址，如果是layer2模式需要这段地址与node节点地址同个网段（同一个二层）；如果是bgp模式没有这个限制。
- 对外宣告：layer2模式使用arp协议，利用节点的mac额外宣告一个loadbalancer的ip（同mac多ip）；bgp模式下节点利用bgp协议与外部网络设备建立邻居，宣告loadbalancer的地址段给外部网络。

