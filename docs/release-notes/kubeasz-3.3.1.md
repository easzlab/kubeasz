## kubeasz 3.3.1 (Slight Heat)

倏忽温风至，因循小暑来。竹喧先觉雨，山暗已闻雷。kubeasz 3.3.1 发布，稳定性和新特性更新说明如下。

### 主要更新

#### 1.新增local insecure registry

为方便集群离线安装，新增本地镜像仓库，用于预存放集群安装所需的组件镜像；建议仅做集群安装时临时使用，不适合对外用作业务应用的镜像存储（harbor 可以作为企业内部镜像仓库应用）。调整ezdown 下载脚本，区分默认组件镜像（`ezdown -D`）自动下载和可选组件镜像（`ezdown -X`）下载并推送至该本地镜像仓库。

#### 2.新增网络检测工具/插件

集群初始安装后，或者运行很久时，非常需要有个工具能够简单检测当前集群各个节点网络是否正常；受 cilium connectivity-check 启发，利用cronjob 检测集群各种网络访问方式是否正常。详解[组件说明](https://github.com/easzlab/kubeasz/blob/master/docs/setup/network-plugin/network-check.md) 

#### 3.更新calico组件支持自动安装calico route reflector

calico 是最流行的网络组件之一；但是当集群节点达到一定数量后，默认的bgp全互联拓扑会导致每个节点需要维护大量BGP邻居信息；本次更新集成了calico-route-reflector自动安装，建议当节点数大于50时必须开启，详见[文档说明](https://github.com/easzlab/kubeasz/blob/master/docs/setup/network-plugin/calico-bgp-rr.md)

#### 4.更新重写cilium组件安装

cilium 可算是最酷的网络组件之一；拥有eBPF光环，以及炫酷的cilium network policy(比k8s原生network policy增强很多)，还有可观测性... 后续项目会加大对cilium组件的更新支持。

#### 5.增加github action 自动同步仓库

自动同步kubeasz项目到国内gitee仓库，方便国内网络环境下访问。


### 组件更新

- k8s: v1.24.2
- coredns: 1.9.3
- pause: 3.7

### 其他

- 大量安装文档更新
- 修复add-node等添加节点时自动添加`/usr/bin/python`软链接
