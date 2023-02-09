<img alt="kubeasz-logo" width="320" height="100" src="pics/kubeasz.svg">  <a href="docs/mixes/conformance.md"><img align="right" alt="conformance-icon" width="75" height="100" src="https://www.cncf.io/wp-content/uploads/2020/07/certified_kubernetes_color-1.png"></a>

**kubeasz** 致力于提供快速部署高可用`k8s`集群的工具, 同时也努力成为`k8s`实践、使用的参考书；基于二进制方式部署和利用`ansible-playbook`实现自动化；既提供一键安装脚本, 也可以根据`安装指南`分步执行安装各个组件。

**kubeasz** 从每一个单独部件组装到完整的集群，提供最灵活的配置能力，几乎可以设置任何组件的任何参数；同时又为集群创建预置一套运行良好的默认配置，甚至自动化创建适合大规模集群的[BGP Route Reflector网络模式](docs/setup/network-plugin/calico-bgp-rr.md)。

- **集群特性** [Master高可用](docs/setup/00-planning_and_overall_intro.md#ha-architecture)、[离线安装](docs/setup/offline_install.md)、[多架构支持(amd64/arm64)](docs/setup/multi_platform.md)
- **集群版本** kubernetes v1.22, v1.23, v1.24, v1.25, v1.26
- **操作系统** CentOS/RedHat 7/8, Debian 9/10, Ubuntu 16.04/18.04/20.04/22.04
- **运行时** [containerd](docs/setup/03-container_runtime.md) v1.6.x, docker 20.10.x (k8s < 1.24)
- **网络** [calico](docs/setup/network-plugin/calico.md), [cilium](docs/setup/network-plugin/cilium.md), [flannel](docs/setup/network-plugin/flannel.md), [kube-ovn](docs/setup/network-plugin/kube-ovn.md), [kube-router](docs/setup/network-plugin/kube-router.md)


**[news]** kubeasz 通过cncf一致性测试 [详情](docs/mixes/conformance.md)

推荐版本对照

<table>
  <thead>
    <tr>
      <td>Kubernetes version</td>
      <td>1.22</td>
      <td>1.23</td>
      <td>1.24</td>
      <td>1.25</td>
      <td>1.26</td>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>kubeasz version</td>
      <td>3.1.1</td>
      <td>3.2.0</td>
      <td>3.3.4</td>
      <td>3.4.5</td>
      <td>3.5.2</td>
    </tr>
  </tbody>
</table>

## 快速指南

单机快速体验k8s集群的测试环境--[AllinOne部署](docs/setup/quickStart.md)

## 安装指南

<table border="0">
    <tr>
        <td><a href="docs/setup/00-planning_and_overall_intro.md">00-规划集群和配置介绍</a></td>
        <td><a href="docs/setup/02-install_etcd.md">02-安装etcd集群</a></td>
        <td><a href="docs/setup/04-install_kube_master.md">04-安装master节点</a></td>
        <td><a href="docs/setup/06-install_network_plugin.md">06-安装集群网络</a></td>
    </tr>
    <tr>
        <td><a href="docs/setup/01-CA_and_prerequisite.md">01-创建证书和安装准备</a></td>
        <td><a href="docs/setup/03-container_runtime.md">03-安装容器运行时</a></td>
        <td><a href="docs/setup/05-install_kube_node.md">05-安装node节点</a></td>
        <td><a href="docs/setup/07-install_cluster_addon.md">07-安装集群插件</a></td>
    </tr>
</table>

## 使用指南

<table border="0">
    <tr>
        <td><strong>常用插件</strong><a href="docs/guide/index.md">+</a></td>
        <td><a href="docs/guide/kubedns.md">DNS</a></td>
        <td><a href="docs/guide/dashboard.md">dashboard</a></td>
        <td><a href="docs/guide/metrics-server.md">metrics-server</a></td>
        <td><a href="docs/guide/prometheus.md">prometheus</a></td>
        <td><a href="docs/guide/efk.md">efk</a></td>
    </tr>
    <tr>
        <td><strong>集群管理</strong><a href="docs/op/op-index.md">+</a></td>
        <td><a href="docs/op/op-node.md">管理node节点</a></td>
        <td><a href="docs/op/op-master.md">管理master节点</a></td>
        <td><a href="docs/op/op-etcd.md">管理etcd节点</a></td>
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
    </tr>
    <tr>
        <td><strong>周边生态</strong></td>
        <td><a href="docs/guide/harbor.md">harbor</a></td>
        <td><a href="docs/guide/helm.md">helm</a></td>
        <td><a href="docs/guide/jenkins.md">jenkins</a></td>
        <td><a href="docs/guide/gitlab/readme.md">gitlab</a></td>
        <td><a href=""></a></td>
    </tr>
</table>

## 沟通交流

- 微信群：k8s&kubeasz实践, 搜索微信号`badtobone`, 请备注（城市-github用户名）, 验证通过会加入群聊。
- 推荐阅读
  - [kubernetes架构师课程](https://www.toutiao.com/c/user/token/MS4wLjABAAAA0YFomuMNm87NNysXeUsQdI0Tt3gOgz8WG_0B3MzxsmI/?tab=article)
  - [kubernetes-the-hard-way](https://github.com/kelseyhightower/kubernetes-the-hard-way)
  - [feisky-Kubernetes 指南](https://github.com/feiskyer/kubernetes-handbook/blob/master/SUMMARY.md)
  - [opsnull 安装教程](https://github.com/opsnull/follow-me-install-kubernetes-cluster)

## 贡献&致谢

欢迎提[Issues](https://github.com/easzlab/kubeasz/issues)和[PRs](docs/mixes/HowToContribute.md)参与维护项目！感谢您的关注与支持！
- [如何 PR](docs/mixes/HowToContribute.md)
- [如何捐赠](docs/mixes/donate.md)

Copyright 2017 gjmzj (jmgaozz@163.com) Apache License 2.0, 详情见 [LICENSE](docs/mixes/LICENSE) 文件。
