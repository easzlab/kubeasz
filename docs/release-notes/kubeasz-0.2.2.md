## kubeasz-0.2.2 发布说明

CHANGELOG:
- 组件更新：
  - k8s v1.11.0
  - etcd v3.3.8 
  - docker 18.03.1-ce
- 功能更新：
  - 更新使用ipvs 配置及[说明文档](https://github.com/gjmzj/kubeasz/blob/master/docs/guide/ipvs.md) 
  - 更新lb节点keepalived使用单播发送vrrp报文，预期兼容公有云上自建LB（待测试）
  - 废弃原 ansible hosts 中变量SERVICE_PROXY
  - 更新haproxy负载均衡算法配置
- 其他修复：
  - fix 变更集群网络的脚本和[文档](https://github.com/gjmzj/kubeasz/blob/master/docs/op/change_k8s_network.md)
  - fix 脚本99.clean.yml清理环境变量
  - fix metrics-server允许的client cert问题
  - fix #242: 添加CA有效期参数，设定CA有效期为15年(131400h) (#245)
  - fix helm安装出现Error: transport is closing (#248)
  - fix harbor点击tag界面出现\"发生未知错误,请稍后再试" (#250)
  - fix 脚本99.clean.yml清理 services softlink (#253)
  - fix kube-apiserver-v1.8 使用真实数量的 apiserver-count (#254)
  - fix 清理ipvs产生的网络接口
