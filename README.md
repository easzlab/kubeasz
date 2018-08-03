# 使用ansible-playbook部署kubernetes集群

![docker](./pics/docker.jpg) ![kube](./pics/kube.jpg) ![ansible](./pics/ansible.jpg)

本项目致力于提供快速部署高可用`k8s`集群的工具，并且也努力成为`k8s`实践、使用的参考书；基于二进制方式部署和利用`ansible-playbook`实现自动化：既提供一键安装脚本，也可以分步执行安装各个组件，同时讲解每一步主要参数配置和注意事项；二进制方式部署有助于理解系统各组件的交互原理和熟悉组件启动参数，有助于快速排查解决实际问题。

**集群特性：`TLS` 双向认证、`RBAC` 授权、多`Master`高可用、支持`Network Policy`**

项目基于`Ubuntu 16.04/CentOS 7`，需要了解基础`kubernetes` `docker` `linux` 知识，关于`ansible`建议阅读 [ansible超快入门](http://weiweidefeng.blog.51cto.com/1957995/1895261) 。

请阅读[项目TodoList](docs/mixes/TodoList.md)和[项目分支说明](docs/branch.md)，欢迎提[Issues](https://github.com/gjmzj/kubeasz/issues)和[PRs](docs/mixes/HowToContribute.md)参与维护项目。

|组件|更新|支持|
|:-|:-|:-|
|OS|&nbsp;|Ubuntu 16.04+, CentOS 7|
|k8s|v1.11.0|v1.8,v1.9,v1.10,v1.11|
|etcd|v3.3.8|v3.1,v3.2,v3.3|
|docker|18.03.1-ce|17.*.*-ce, 18.*.*-ce|
|network|&nbsp;|calico v3.0, flannel v0.10, kube-router|

-注：集群用到的所有二进制文件已打包好供下载 [https://pan.baidu.com/s/1c4RFaA](https://pan.baidu.com/s/1c4RFaA)  

## 快速指南

单机快速体验k8s集群的测试、开发环境--[AllinOne部署](docs/quickStart.md)

## 安装步骤

<table border="0">
    <tr>
        <td><a href="docs/00-集群规划和基础参数设定.md">00-规划集群和安装概览</a></td>
        <td><a href="docs/02-安装etcd集群.md">02-安装etcd集群</a></td>
        <td><a href="docs/04-安装kube-master节点.md">04-安装master节点</a></td>
        <td><a href="docs/06-安装网络组件.md">06-安装集群网络</a></td>
    </tr>
    <tr>
        <td><a href="docs/01-创建CA证书和环境配置.md">01-创建证书和安装准备</a></td>
        <td><a href="docs/03-安装docker服务.md">03-安装docker服务</a></td>
        <td><a href="docs/05-安装kube-node节点.md">05-安装node节点</a></td>
        <td><a href="docs/07-安装集群插件.md">07-安装集群插件</a></td>
    </tr>
</table>

## 使用指南

<table border="0">
    <tr>
        <td><strong>常用插件</strong></td>
        <td><a href="docs/guide/kubedns.md">kube-dns</a></td>
        <td><a href="docs/guide/dashboard.md">dashboard</a></td>
        <td><a href="docs/guide/metrics-server.md">metrics-server</a></td>
        <td><a href="docs/guide/prometheus.md">prometheus</a></td>
        <td><a href="docs/guide/index.md">更多...</a></td>
    </tr>
    <tr>
        <td><strong>集群运维</strong></td>
        <td><a href="docs/op/AddNode.md">增加node</a></td>
        <td><a href="docs/op/AddMaster.md">增加master</a></td>
        <td><a href="docs/op/upgrade.md">升级集群</a></td>
        <td><a href="docs/op/cluster_restore.md">备份恢复</a></td>
        <td><a href="docs/op/op-index.md">更多...</a></td>
    </tr>
    <tr>
        <td><strong>特性实验</strong></td>
        <td><a href="docs/guide/networkpolicy.md">NetworkPolicy</a></td>
        <td><a href="docs/guide/rollingupdateWithZeroDowntime.md">RollingUpdate</a></td>
        <td><a href="docs/guide/hpa.md">HPA</a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
    </tr>
    <tr>
        <td><strong>其他</strong></td>
        <td><a href="docs/guide/harbor.md">harbor部署</a></td>
        <td><a href="docs/guide/ingress.md">ingress入口</a></td>
        <td><a href="docs/guide/helm.md">helm工具</a></td>
        <td><a href="docs/guide/jenkins.md">jenkins集成</a></td>
        <td><a href=""></a></td>
    </tr>
</table>

## 沟通交流

- 微信群：k8s&kubeasz实践，搜索微信号`badtobone`，请备注（城市-github用户名），验证通过会加入群聊。
- 推荐阅读：[rootsongjc-Kubernetes指南](https://github.com/rootsongjc/kubernetes-handbook) [feisky-Kubernetes指南](https://github.com/feiskyer/kubernetes-handbook/blob/master/zh/SUMMARY.md) [opsnull-安装教程](https://github.com/opsnull/follow-me-install-kubernetes-cluster)

## 贡献&致谢

感谢所有为项目提交 `Issues`和`PRs` 的贡献者！

- [如何贡献](docs/mixes/HowToContribute.md)

Copyright 2017 gjmzj (jmgaozz@163.com) Apache License 2.0，详情见 [LICENSE](docs/mixes/LICENSE) 文件。
