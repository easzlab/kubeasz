# kubeasz

`kubeasz`致力于提供快速部署高可用`k8s`集群的工具, 并且也努力成为`k8s`实践、使用的参考书；基于二进制方式部署和利用`ansible-playbook`实现自动化：即提供一键安装脚本, 也可以分步执行安装各个组件, 同时讲解每一步主要参数配置和注意事项。

**集群特性：`TLS`双向认证、`RBAC`授权、多`Master`高可用、支持`Network Policy`、备份恢复**

|组件|支持|
|:-|:-|
|OS|Ubuntu 16.04+, CentOS/RedHat 7|
|k8s|v1.8, v1.9, v1.10, v1.11, v1.12, v1.13|
|etcd|v3.1, v3.2, v3.3|
|docker|17.03.2-ce, 18.06.1-ce, 18.09.2|
|network|calico, cilium, flannel, kube-router|

- 注：集群用到的所有二进制文件已打包好供下载 [https://pan.baidu.com/s/1c4RFaA](https://pan.baidu.com/s/1c4RFaA)  

请阅读[项目TodoList](docs/mixes/TodoList.md)和[项目分支说明](docs/mixes/branch.md), 欢迎提[Issues](https://github.com/gjmzj/kubeasz/issues)和[PRs](docs/mixes/HowToContribute.md)参与维护项目。

## 快速指南

单机快速体验k8s集群的测试、开发环境--[AllinOne部署](docs/setup/quickStart.md)

## 安装指南

<table border="0">
    <tr>
        <td><a href="docs/setup/00-planning_and_overall_intro.md">00-规划集群和安装概览</a></td>
        <td><a href="docs/setup/02-install_etcd.md">02-安装etcd集群</a></td>
        <td><a href="docs/setup/04-install_kube_master.md">04-安装master节点</a></td>
        <td><a href="docs/setup/06-install_network_plugin.md">06-安装集群网络</a></td>
    </tr>
    <tr>
        <td><a href="docs/setup/01-CA_and_prerequisite.md">01-创建证书和安装准备</a></td>
        <td><a href="docs/setup/03-install_docker.md">03-安装docker服务</a></td>
        <td><a href="docs/setup/05-install_kube_node.md">05-安装node节点</a></td>
        <td><a href="docs/setup/07-install_cluster_addon.md">07-安装集群插件</a></td>
    </tr>
</table>

- 公有云部署请阅读 [使用kubeasz在公有云上创建k8s集群](docs/setup/kubeasz_on_public_cloud.md)

## 使用指南

<table border="0">
    <tr>
        <td><strong>常用插件</strong><a href="docs/guide/index.md">+</a></td>
        <td><a href="docs/guide/kubedns.md">DNS</a></td>
        <td><a href="docs/guide/dashboard.md">dashboard</a></td>
        <td><a href="docs/guide/metrics-server.md">metrics-server</a></td>
        <td><a href="docs/guide/prometheus.md">prometheus</a></td>
        <td><a href="docs/guide/efk.md">efk</a></td>
        <td><a href="docs/guide/ingress.md">ingress</a></td>
    </tr>
    <tr>
        <td><strong>集群管理</strong><a href="docs/op/op-index.md">+</a></td>
        <td><a href="docs/op/AddNode.md">增加node节点</a></td>
        <td><a href="docs/op/AddMaster.md">增加master节点</a></td>
        <td><a href="docs/op/op-etcd.md">管理etcd集群</a></td>
        <td><a href="docs/op/del_one_node.md">删除节点</a></td>
        <td><a href="docs/op/upgrade.md">升级集群</a></td>
        <td><a href="docs/op/cluster_restore.md">备份恢复</a></td>
    </tr>
    <tr>
        <td><strong>特性实验</strong></td>
        <td><a href="docs/guide/networkpolicy.md">NetworkPolicy</a></td>
        <td><a href="docs/guide/rollingupdateWithZeroDowntime.md">RollingUpdate</a></td>
        <td><a href="docs/guide/hpa.md">HPA</a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
    </tr>
    <tr>
        <td><strong>周边生态</strong></td>
        <td><a href="docs/guide/harbor.md">harbor</a></td>
        <td><a href="docs/guide/helm.md">helm</a></td>
        <td><a href="docs/guide/jenkins.md">jenkins</a></td>
        <td><a href="docs/guide/gitlab/readme.md">gitlab</a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
    </tr>
    <tr>
        <td><strong>应用实践</strong></td>
        <td><a href="docs/practice/java_war_app.md">java应用部署</a></td>
        <td><a href="docs/practice/es_cluster.md">elasticsearch集群</a></td>
        <td><a href="docs/practice/mariadb_cluster.md">mariadb集群</a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
        <td><a href=""></a></td>
    </tr>
</table>

## 沟通交流

- 微信群：k8s&kubeasz实践, 搜索微信号`badtobone`, 请备注（城市-github用户名）, 验证通过会加入群聊。
- 推荐阅读：[feisky-Kubernetes指南](https://github.com/feiskyer/kubernetes-handbook/blob/master/SUMMARY.md) [rootsongjc-Kubernetes指南](https://github.com/rootsongjc/kubernetes-handbook) [opsnull-安装教程](https://github.com/opsnull/follow-me-install-kubernetes-cluster)

## 贡献&致谢

感谢所有为项目提交 `Issues`和`PRs` 的贡献者！感谢[捐赠](docs/mixes/donate.md)鼓励！

- [如何 PR](docs/mixes/HowToContribute.md)

Copyright 2017 gjmzj (jmgaozz@163.com) Apache License 2.0, 详情见 [LICENSE](docs/mixes/LICENSE) 文件。
