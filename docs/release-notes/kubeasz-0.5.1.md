## kubeasz-0.5.1 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.13.2, v1.12.4, v1.11.6, v1.10.12
  - calico v3.3.2
  - coredns 1.2.6
- 集群安装：
  - 更新 calico 3.3.2，并保留3.2.4可选
  - 修复特定环境下lb节点变量LB_IF自动设置错误
  - 移除 kube-node csr 请求批准部分（PR #399）
  - 添加支持 RedHat （PR #431）
  - 修改 docker 存储的目录设置（PR #436）
  - 更新 kube-schedule 监听参数 （PR #440）
  - 安装流程增加等待 ETCD 同步完成再返回成功（PR #420）
  - 增加 pod-infra-container 可选择配置
  - 增加 nginx-ingress manifests
- 文档更新：
  - **增加 [calico 设置route reflector文档](https://github.com/gjmzj/kubeasz/blob/master/docs/setup/network-plugin/calico-bgp-rr.md)**，大规模k8s集群使用calico网络必读
  - 部分文档更新优化，部分文档中内部链接修复（PR #429）
  - 增加 dashboard ingress [配置示例](https://github.com/gjmzj/kubeasz/blob/master/docs/guide/ingress-tls.md#%E9%85%8D%E7%BD%AE-dashboard-ingress)
- 其他：
  - 添加 helm tls 环境变量（PR #398）
  - 修复 dashboard ingress 配置（issue #403）
