## kubeasz-0.1.0 发布说明

`kubeasz`项目第一个独立版本发布，使用`ansible playbook`自动化安装k8s集群（目前支持v1.8/v1.9/v1.10）和主要插件，方便部署和灵活配置集群；

CHANGELOG:
- 组件更新：
  - kubernetes v1.10.4, v1.9.8, v1.8.12
  - etcd v3.3.6
- 安全更新：
  - 修复kubelet匿名访问漏洞（感谢 cqspirit #192 提醒）
- 功能更新：
  - 增加helm安全部署及说明
  - 增加prometheus部署及说明
  - 增加jenkins部署及说明（感谢 lusyoe #208 ）
- 脚本更新：
  - 精简 inventory（/etc/ansible/hosts）配置项
    - 移动calico/flannel配置至对应的roles/defaults/main.yml
    - 取消变量NODE_IP，使用内置变量inventory_hostname代替
    - 取消lb组变量设置，自动完成
    - 取消etcd相关集群变量设置，自动完成
  - 增加集群版本K8S_VER变量，为兼容k8s v1.8安装
  - 增加修改AIO部署的系统IP的脚本和说明(docs/op/change_ip_allinone.md)
  - 增加设置node角色
  - 修改OS安全加固脚本为可选安装
- 其他：
  - 修复calico-controller多网卡问题
  - 修改manifests/apiserver参数兼容k8s v1.8
  - 简化新增master/node节点步骤
  - 优化ansible配置参数
  - 更新 harbor 1.5.1及文档修复（感谢 lusyoe #224 ）
  - 更新 kube-dns 1.14.10
  - 丰富dashboard文档（ #182 ）
  - 修复selinux关闭（ #194 ）
