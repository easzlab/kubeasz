## 06-安装网络组件

首先回顾下K8S网络设计原则，在配置集群网络插件或者实践K8S 应用/服务部署请牢记这些原则：

- 1.每个Pod都拥有一个独立IP地址，Pod内所有容器共享一个网络命名空间
- 2.集群内所有Pod都在一个直接连通的扁平网络中，可通过IP直接访问
  - 所有容器之间无需NAT就可以直接互相访问
  - 所有Node和所有容器之间无需NAT就可以直接互相访问
  - 容器自己看到的IP跟其他容器看到的一样
- 3.Service cluster IP只可在集群内部访问，外部请求需要通过NodePort、LoadBalance或者Ingress来访问

`Container Network Interface (CNI)`是目前CNCF主推的网络模型，它由两部分组成：

- CNI Plugin负责给容器配置网络，它包括两个基本的接口
  - 配置网络: AddNetwork(net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
  - 清理网络: DelNetwork(net *NetworkConfig, rt *RuntimeConf) error
- IPAM Plugin负责给容器分配IP地址

Kubernetes Pod的网络是这样创建的：
- 0. 每个Pod除了创建时指定的容器外，都有一个kubelet启动时指定的`基础容器`，即`pause`容器 
- 1. kubelet创建`基础容器`生成network namespace
- 2. kubelet调用网络CNI driver，由它根据配置调用具体的CNI 插件
- 3. CNI 插件给`基础容器`配置网络
- 4. Pod 中其他的容器共享使用`基础容器`的网络

本项目基于CNI driver 调用各种网络插件来配置kubernetes的网络，常用CNI插件有 `flannel` `calico` `cilium`等等，这些插件各有优势，也在互相借鉴学习优点，比如：在所有node节点都在一个二层网络时候，flannel提供hostgw实现，避免vxlan实现的udp封装开销，估计是目前最高效的；calico也针对L3 Fabric，推出了IPinIP的选项，利用了GRE隧道封装；因此这些插件都能适合很多实际应用场景。

项目当前内置支持的网络插件有：`calico` `cilium` `flannel` `kube-ovn` `kube-router`

### 安装讲解

- [安装calico](network-plugin/calico.md)
- [安装cilium](network-plugin/cilium.md)  
- [安装flannel](network-plugin/flannel.md)
- [安装kube-ovn](network-plugin/kube-ovn.md) 暂未更新
- [安装kube-router](network-plugin/kube-router.md) 暂未更新

### 参考
- [kubernetes.io networking docs](https://kubernetes.io/docs/concepts/cluster-administration/networking/) 
- [feiskyer-kubernetes指南网络章节](https://github.com/feiskyer/kubernetes-handbook/blob/master/zh/network/network.md)


[后一篇](07-install_cluster_addon.md)
