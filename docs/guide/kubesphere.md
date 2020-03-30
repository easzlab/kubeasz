# 在 Kubernetes 安装 KubeSphere 容器平台

## 什么是 KubeSphere

[KubeSphere](https://github.com/kubesphere/kubesphere) 是在 [Kubernetes](https://kubernetes.io) 之上构建的**开源企业级容器平台**，提供全栈的 IT 自动化运维的能力，简化企业的 DevOps 工作流。KubeSphere 作为一个**全栈的容器平台**，不仅支持**安装和纳管原生 Kubernetes**，还设计了一套完整的管理界面方便开发者与运维人员在一个**统一的平台** 中安装与管理最常用的云原生工具，**从业务视角提供一致的用户体验来降低复杂性**。目前版本提供以下功能：

|功能 |介绍 |
| --- | ---|
| Kubernetes 集群搭建与运维 | 支持在线 & 离线安装、升级与扩容 K8s 集群，支持安装 “云原生全家桶” |
| Kubernetes 资源可视化管理 | 比 K8s 原生 Dashboard 功能更丰富的控制面板，支持向导式创建与管理 K8s 资源 |
| 基于 Jenkins 的 DevOps 系统 | 支持图形化与脚本两种方式构建 CI/CD 流水线，内置 Source/Binary to Image 等 CD 工具 |
| 应用商店与应用生命周期管理 | 内置 Redis、MySQL 等十个常用应用，基于 Helm 提供应用上传、审核、发布、部署、下架等操作 |
| 基于 Istio 的微服务治理 (Service Mesh) | 提供可视化无代码侵入的 **灰度发布、熔断、流量治理与流量拓扑、分布式 Tracing** |
| 多租户管理 | 提供基于角色的细粒度多租户统一认证，支持 **对接企业 LDAP/AD**，提供多层级的权限管理 |
| 丰富的可观察性功能 | UI 提供集群/工作负载/Pod/容器等多维度的监控、日志、告警与通知 |
| 存储管理 | 支持对接 Ceph、GlusterFS、NFS，支持可视化管理 PVC、PV、StorageClass |
| 网络管理 | 支持 Calico、Flannel，提供 Porter LB 帮助暴露物理环境 K8s 集群的 LoadBalancer 服务 |
| GPU support | 集群支持添加 GPU 与 vGPU，可运行 TensorFlow 等 ML 框架 |


## 在 Kubernetes 与 Kubeasz 之上安装 KubeSphere

KubeSphere 可以安装在任何私有或托管的 Kubernetes、私有云、公有云、VM 或物理环境之上，并且所有功能组件都是可插拔的。当使用 Kubeasz 完成 K8s 集群的安装后，可参考以下步骤在 Kubernetes 上安装 KubeSphere v2.1.0。

**前提条件**

> - `Kubernetes 版本` ： `1.13.0 ≤ K8s version < 1.16`；
> - `Helm 版本`: `2.10.0 ≤ Helm ＜ 3.0.0`，且已安装了 Tiller（v3.0 支持 Helm v3）；参考 [如何安装与配置 Helm](https://devopscube.com/install-configure-helm-kubernetes/)；
> - 集群的可用 CPU > 1 C，可用内存 > 2 G；且集群能够访问外网
> - 集群已有默认的存储类型（StorageClass）；
>
> 以上四项可参考 [前提条件](https://kubesphere.io/docs/v2.1/zh-CN/installation/prerequisites/) 进行验证。

1. 若待安装的环境满足以上条件则可以通过一条命令部署 KubeSphere。

```yaml
$ kubectl apply -f https://raw.githubusercontent.com/kubesphere/ks-installer/master/kubesphere-minimal.yaml
```

2. 查看 ks-installer 安装过程中产生的动态日志，等待安装成功（约 10 min 左右）：

```bash
$ kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l app=ks-install -o jsonpath='{.items[0].metadata.name}') -f
```

![](https://pek3b.qingstor.com/kubesphere-docs/png/20191005195724.png)

3. 当 KubeSphere 的所有 Pod 都为 Running 则说明安装成功。使用 `http://IP:30880` 访问 Dashboard，默认账号为 `admin/P@88w0rd`。


**Tips**：KubeSphere 在 K8s 默认仅开启 **最小化安装**，执行以下命令开启可插拔功能组件的安装，开启安装前确认您的机器资源已符合 [资源最低要求](https://kubesphere.io/docs/v2.1/zh-CN/installation/intro/#%E5%8F%AF%E6%8F%92%E6%8B%94%E5%8A%9F%E8%83%BD%E7%BB%84%E4%BB%B6%E5%88%97%E8%A1%A8)。

```
$ kubectl edit cm -n kubesphere-system ks-installer
```

## 延伸阅读

- [安装 Kubeasz 与 KubeSphere](https://kubesphere.com.cn/forum/d/716-play-with-kubesphere-and-kubeasz)
- [在 Linux 完整安装 KubeSphere 与 Kubernetes](https://kubesphere.com.cn/docs/v2.1/zh-CN/installation/intro/)
- [KubeSphere 官网](https://kubesphere.com.cn/zh-CN/)


