## kubeasz-1.3.0 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.14.3 v1.13.7
  - calico: v3.4.4
  - kube-router: v0.3.1
  - traefik v1.7.12
- 集群安装：
  - 重写容器化运行kubeasz的脚本 tools/easzup
  - 增加网络插件 kube-ovn 及说明 docs/setup/network-plugin/kube-ovn.md 
  - 增加支持 containerd 离线镜像导入
  - 增加 kubelet 可选配置 kube-reserved/system-reserved 资源预留（默认开启）
  - 移除 dockerfiles 内容，转移至 https://github.com/kubeasz/dockerfiles
  - 添加安装 docker/containerd 互斥判断
  - 废弃 ansible hosts 中变量CLUSTER_DNS_SVC_IP等，配置更精简方便
  - 增加默认自动安装 traefik
  - 优化haproxy 最大连接数为50000
- 文档：
  - 更新ha-1x架构等文档
  - 更新容器化运行 kubeasz 文档 docs/setup/docker_kubeasz.md
- 其他：
  - fix: 双网卡下 apiserver endpoint 可能错误 #479
  - fix: 安全加固（roles/os-harden）默认允许ip_forward
  - fix: docker/containerd 安装检查 
  - fix: docker运行kubeasz时判断ansible控制端与deploy节点分离的条件
  - fix: docker运行kubeasz时清理脚本 99.clean.yml
