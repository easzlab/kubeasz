# 在 Kubernetes 安装 KubeSphere 容器平台

## 什么是 KubeSphere

[KubeSphere](https://github.com/kubesphere/kubesphere) 是在 [Kubernetes](https://kubernetes.io) 之上构建的面向云原生应用的**开源容器平台**，支持多云与多集群管理，提供全栈的 IT 自动化运维能力，简化企业的 DevOps 工作流。它的架构可以非常方便地使第三方应用与云原生生态组件进行即插即用 (plug-and-play) 的集成。

KubeSphere 作为一个**全栈的多租户容器平台**，不仅支持**安装和纳管原生 Kubernetes**，还设计了一套完整的管理界面，方便开发者与运维人员在一个**统一的平台**中安装与管理最常用的云原生工具，**从业务视角提供一致的用户体验来降低复杂性**。目前最新的 3.0 版本提供以下功能：

|功能 |介绍 |
| --- | ---|
| Kubernetes 集群搭建与运维 | 支持在线 & 离线安装、升级与扩容 Kubernetes 集群，支持安装 “云原生全家桶” |
| Kubernetes 资源可视化管理 | 比 Kubernetes 原生 Dashboard 功能更丰富的控制面板，支持向导式创建与管理 Kubernetes 资源 |
| 基于 Jenkins 的 DevOps 系统 | 支持图形化与脚本两种方式构建 CI/CD 流水线，内置 Source/Binary to Image 等 CD 工具 |
| 应用商店与应用生命周期管理 | 内置 Redis、MySQL 等十五个常用应用，基于 Helm 提供应用上传、审核、发布、部署、下架等操作 |
| 基于 Istio 的微服务治理 (Service Mesh) | 提供可视化无代码侵入的**灰度发布、熔断机制、流量治理与流量拓扑、分布式链路追踪** |
| 多租户管理 | 提供基于角色的细粒度多租户统一认证，支持**对接企业 LDAP/AD**，提供多层级的权限管理 |
| 丰富的可观察性功能 | UI 提供集群/工作负载/Pod/容器等多维度的监控、事件/日志查询、告警与通知管理 |
| 存储管理 | 支持对接 Ceph、GlusterFS、NFS，支持可视化管理 PVC、PV、StorageClass |
| 网络管理 | 支持 Calico、Flannel，提供 Porter LB 帮助暴露物理环境 Kubernetes 集群的 LoadBalancer 服务 |
| GPU support | 集群支持添加 GPU 与 vGPU，可运行 TensorFlow 等 ML 框架 |


## 在 Kubernetes 与 Kubeasz 之上安装 KubeSphere

作为一个轻量化容器平台，KubeSphere 可以安装在任何私有或托管的 Kubernetes、虚拟机、裸机、本地环境、公有云、混合云之上，并且所有功能组件都是可插拔的。当使用 Kubeasz 完成 Kubernetes 集群的安装后，可参考以下步骤在 Kubernetes 上安装 KubeSphere。

**前提条件**

> - Kubernetes 版本必须是：1.15.x、1.16.x、1.17.x 或 1.18.x；
> - 您的机器满足最低硬件要求：CPU > 1 Core，可用内存 > 2 G；
> - 安装之前，Kubernetes 集群已配置**默认**存储类型 (StorageClass)；
> - 当使用 `--cluster-signing-cert-file` 和 `--cluster-signing-key-file` 参数启动时，在 `kube-apiserver` 中会激活 CSR 签名功能。请参见 [RKE 安装问题](https://github.com/kubesphere/kubesphere/issues/1925#issuecomment-591698309)；
> - 有关在 Kubernetes 上安装 KubeSphere 的准备工作，请参见[准备工作](https://kubesphere.io/zh/docs/installing-on-kubernetes/introduction/prerequisites/)。
>

1. 若待安装的环境满足以上条件，则可以执行以下命令部署 KubeSphere：

   ```yaml
   kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.0.0/kubesphere-installer.yaml
   
   kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.0.0/cluster-configuration.yaml
   ```

2. 等待安装成功（取决于您的网络状况，约十几至二十几分钟不等），运行以下命令查看安装日志：

   ```bash
   kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l app=ks-install -o jsonpath='{.items[0].metadata.name}') -f
   ```

  ![](https://pek3b.qingstor.com/kubesphere-docs/png/20191005195724.png)

3. 使用 `kubectl get pod --all-namespaces` 查看所有 Pod 是否在 KubeSphere 的相关命名空间中正常运行。如果是，请通过以下命令检查控制台的端口（默认为 `30880`）：

   ```bash
   kubectl get svc/ks-console -n kubesphere-system
   ```

4. 请确保在安全组中打开了端口 `30880`，并通过 NodePort `(IP:30880)` 使用默认帐户和密码 `(admin/P@88w0rd)` 访问 Web 控制台。

5. 登录控制台后，您可以在**服务组件**中查看不同组件的状态。如果要使用相关服务，可能需要等待某些组件启动并运行。

**Tips**：若要在 KubeSphere 中启用其他组件，请参见[启用可插拔组件](https://kubesphere.io/zh/docs/pluggable-components/)。开启安装前确认您的机器资源已符合[资源最低要求](https://kubesphere.io/zh/docs/pluggable-components/overview/)。

## 延伸阅读

- [安装 Kubeasz 与 KubeSphere](https://kubesphere.com.cn/forum/d/716-play-with-kubesphere-and-kubeasz)
- [在 Linux 完整安装 KubeSphere 与 Kubernetes](https://kubesphere.io/zh/docs/installing-on-linux/introduction/intro/)
- [KubeSphere 官网](https://kubesphere.io/zh/)
- [常见问题](https://kubesphere.io/zh/docs/faq/)


