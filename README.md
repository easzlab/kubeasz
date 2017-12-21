# 利用Ansible部署kubernetes集群

![docker](./pics/docker.jpg) ![kube](./pics/kube.jpg) ![ansible](./pics/ansible.jpg)

本系列文档致力于提供快速部署高可用`k8s`集群的工具，并且也努力成为`k8s`实践、使用的参考书；基于二进制方式部署和利用`ansible-playbook`实现自动化：既提供一键安装脚本，也可以分步执行安装各个组件，同时讲解每一步主要参数配置和注意事项。

**集群特性：`TLS` 双向认证、`RBAC` 授权、多`Master`高可用、支持`Network Policy`**

**二进制方式部署优势：有助于理解系统各组件的交互原理和熟悉组件启动参数，有助于快速排查解决实际问题**

文档基于`Ubuntu 16.04/CentOS 7`，其他系统需要读者自行替换部分命令；由于使用经验有限和简化脚本考虑，已经尽量避免`ansible-playbook`的高级特性和复杂逻辑。

你可能需要掌握基本`kubernetes` `docker` `linux shell` 知识，关于`ansible`建议阅读 [ansible超快入门](http://weiweidefeng.blog.51cto.com/1957995/1895261) 基本够用。

欢迎提`Issues`和`PRs`参与维护项目。

## 组件版本

1. kubernetes	v1.9.0
1. etcd		v3.2.11
1. docker	17.09.1-ce
1. calico/node	v2.6.3

+ 附：集群用到的所有二进制文件已打包好供下载 [https://pan.baidu.com/s/1c4RFaA](https://pan.baidu.com/s/1c4RFaA)
+ 注：`Kubernetes v1.8.x` 版本请切换到项目分支 `v1.8`, 若你需要从v1.8 升级至 v1.9，请参考 [升级注意](docs/upgrade.md)

## 快速指南

单机快速体验k8s集群的测试、开发环境--[AllinOne部署](docs/quickStart.md)；在国内的网络环境下要比官方的minikube方便、简单很多。

## 安装步骤

1. [集群规划和基础参数设定](docs/00-集群规划和基础参数设定.md)
1. [创建CA证书和环境配置](docs/01-创建CA证书和环境配置.md)
1. [安装etcd集群](docs/02-安装etcd集群.md)
1. [配置kubectl命令行工具](docs/03-配置kubectl命令行工具.md)
1. [安装docker服务](docs/04-安装docker服务.md)
1. [安装calico网络组件](docs/05-安装calico网络组件.md)
1. [安装kube-master节点](docs/06-安装kube-master节点.md)
1. [安装kube-node节点](docs/07-安装kube-node节点.md)

## 使用指南

基本k8s集群安装完成后，需要安装一些常用插件(`kubedns` `dashboard` `ingress`等)；接着介绍一些集群操作场景和思路；然后介绍一些应用部署实践，请根据这份[目录](docs/guide/index.md)阅读你所感兴趣的内容。尚在更新中...

## 参考阅读

1. 建议阅读 [rootsongjc-Kubernetes指南](https://github.com/rootsongjc/kubernetes-handbook) 原理和实践指南。
1. 建议阅读 [feisky-Kubernetes指南](https://github.com/feiskyer/kubernetes-handbook/blob/master/SUMMARY.md) 原理和部署章节。
1. 建议阅读 [opsnull-安装教程](https://github.com/opsnull/follow-me-install-kubernetes-cluster) 二进制手工部署。

## 版权

Copyright 2017 gjmzj (jmgaozz@163.com)

Apache License 2.0，详情见 [LICENSE](LICENSE) 文件。

如果觉得这份文档对你有帮助，请支付宝扫描下方的二维码进行捐赠，谢谢！

![donate](./pics/alipay.png) 
