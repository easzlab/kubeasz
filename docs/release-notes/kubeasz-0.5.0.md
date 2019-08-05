## kubeasz-0.5.0 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.12.3, v1.11.5, v1.10.11
  - calico v3.2.4
  - helm v2.11.0
  - traefik 1.7.4
- 集群安装：
  - 更新集群升级脚本和[文档](https://github.com/easzlab/kubeasz/blob/master/docs/op/upgrade.md)，关注[安全漏洞](https://mp.weixin.qq.com/s/Q8XngAr5RuL_irRscbVbKw)
  - 集成 metallb 作为自有硬件 k8s 集群的 LoadBalancer 实现
  - 支持[修改 APISERVER 证书](https://github.com/easzlab/kubeasz/blob/master/docs/op/ch_apiserver_cert.md)
  - 增加 ingress nodeport 负载转发的脚本与[文档](https://github.com/easzlab/kubeasz/blob/master/docs/op/loadballance_ingress_nodeport.md)
  - 增加 https ingress 配置和[文档](https://github.com/easzlab/kubeasz/blob/master/docs/guide/ingress-tls.md)
  - 增加 kubectl 只读访问权限配置和[文档](https://github.com/easzlab/kubeasz/blob/master/docs/op/readonly_kubectl.md)
  - 增加 apiserver 配置支持 istio sidecar自动注入webhook (#375)
  - 初始化集群节点设置 net.netfilter.nf_conntrack_max=1000000
  - 取消多主集群LB_IF参数设置，自动生成以避免人为配置疏忽
- 文档更新：
  - 更新[kubeasz 公有云安装文档](https://github.com/easzlab/kubeasz/blob/master/docs/setup/kubeasz_on_public_cloud.md)
  - 更新[metallb 文档](https://github.com/easzlab/kubeasz/blob/master/docs/guide/metallb.md)
  - 更新[dashboard 文档](https://github.com/easzlab/kubeasz/blob/master/docs/guide/dashboard.md)，支持只读权限设置
  - 新增istio安装说明
- 其他：
  - 修复内核4.19加载nf_conntrack (#366)
  - 修复 calico controller 中 NodePorts 的自动配置
  - 取消 helms 别名设置
  - 升级jenkins-lts版本和插件版本 (#358)
  - 修复阿里云nas动态pv脚本
