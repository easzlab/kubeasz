## kubeasz-0.4.0 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.12.1, v1.10.8, v1.9.11 [注意 v1.12.1 kubelet日志bug](https://github.com/kubernetes/kubernetes/issues/69503)
  - docker: 18.06.1-ce (选择k8s官方测试稳定的版本)
  - metrics-server: v0.3.1
  - coredns: 1.2.2, kube-dns 1.14.13
  - heapster v1.5.4
  - traefik 1.7.2
- 集群安装：
  - **更新 kubelet使用 webhook方式认证/授权** ，提高集群安全性
  - 调整安装步骤中 kubectl 命令的执行以兼容公有云部署
  - 调整部分安装步骤以兼容`ansible`执行节点与`deploy`节点分离
  - 更新节点的安全加固脚本[ansible-os-hardening 5.0.0](https://github.com/dev-sec/ansible-os-hardening)
- 文档更新：
  - 新增`elasticsearch`集群[部署实践](https://github.com/easzlab/kubeasz/blob/master/docs/practice/es_cluster.md)
  - 更新[kubeasz 公有云安装文档](https://github.com/easzlab/kubeasz/blob/master/docs/setup/kubeasz_on_public_cloud.md)
  - 调整集群安装步骤文档目录及修改使用英文文件名
  - 修改部分脚本内部注释为英文
- 其他：
  - 升级 promethus chart 7.1.4，grafana chart 1.16.0
  - 升级 jenkins 安全插件和 k8s 插件版本 (#325)
  - 修复 新增 master 节点时报变量未定义错误
  - 修复 ipvs 模式下网络组件偶尔连不上`kubernetes svc`的错误
  - 修复 Ansible 2.7 环境下 yum/apt 安装多个软件包的 DEPRECATION WARNING (#334)
  - 修复 chrony 与 ntp 共存冲突问题 (#341)
  - 修复 CentOS 下使用 ipvs 模式需依赖 conntrack-tools 软件包
  - 修复 tools/change_k8s_network.yml 脚本
