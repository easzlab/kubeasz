## kubeasz-0.6.0 发布说明

- Note: 本次为 kubeasz-0.x 最后一次版本发布，它将被并入 release-0 分支，停止主要更新，仅做 bug 修复版本；后续 master 分支将开始 kubeasz-1.x 版本发布。
- Action Required: 本次更新修改 ansible hosts 文件，如需要更新已有项目使用，请按照 example 目录中的对应例子修改`/etc/ansible/hosts`文件。

CHANGELOG:
- 组件更新：
  - k8s: v1.13.3
  - calico v3.4.1
  - flannel v0.11.0-amd64
  - docker 18.09.2
  - harbor 1.6.3
  - helm/tiller: v2.12.3
- 集群安装：
  - **增加添加/删除 etcd 节点**脚本和[文档](https://github.com/gjmzj/kubeasz/blob/master/docs/op/op-etcd.md)
  - **增加可选配置附加负载均衡节点（ex-lb）**，可用于负载均衡 NodePort 方式暴露的服务
  - 更新删除节点脚本和[文档](https://github.com/gjmzj/kubeasz/blob/master/docs/op/del_one_node.md)
  - 优化增加 node 和增加 master 节点流程
  - 更新 harbor 安装流程和文档
  - 优化 prepare tasks，避免把证书和 kubeconfig 分发到不需要的节点
  - 更新 prometheus 告警发送钉钉配置和[文档](https://github.com/gjmzj/kubeasz/blob/master/docs/guide/prometheus.md#%E5%8F%AF%E9%80%89-%E9%85%8D%E7%BD%AE%E9%92%89%E9%92%89%E5%91%8A%E8%AD%A6)
  - 增加使用 helm 部署 mariadb 集群和文档
  - 增加 k8s 官方 mysql 集群示意配置
  - 增加使用 helm 部署 redis-ha 集群
  - 增加开机启动 k8s 相关内核模块配置
  - 更新 calico 3.4.1，并保留版本 3.3.x/3.2.x 可选
- 文档更新：
  - **增加 gitlab-ci 文档**, https://github.com/gjmzj/kubeasz/blob/master/docs/guide/gitlab/readme.md
  - 部分文档更新（helm/dns/chrony）
- 其他：
  - 修复为兼容k8s版本 <= 1.11，revert PR #440
  - 修复清除iptables规则时无法连接节点（PR #453 by PowerDos）
  - 添加开启docker远程API选项（默认关闭）（PR #444 by lusyoe）
  - 修复 calico 3.3.x rbac 配置（PR #447 by sunshanpeng）
  - 增加 coredns 和 calico 的 metrics 监控选项（PR #447 by sunshanpeng）
  - 添加 helm 离线安装方法说明（doc/guide/helm.md）（PR #443 by j4ckzh0u）
