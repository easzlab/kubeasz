## kubeasz-3.0.0 (the Beginning of Spring)

2021春快到了，kubeasz带来全新3.x版本，是继2.x基础上做了一些小优化和创新，力求更加整洁和实用。主要更新点：

### 优化多集群创建、管理逻辑

根目录新增 clusters 子目录，用于存放不同集群的配置；现在 ezctl 命令行天然支持多集群管理操作，统一创建、管理，互不影响；建议deploy节点独立出来，具体集群创建、管理操作可以参考 docs/setup/ezctl.md

### 配置集中，组件版本统一设置

模版配置文件 example/config.yml 是把原先 roles/xxxx/defaults/main.yml 配置合并后的全局配置文件；每创建一个集群会从这个模版派生一个实例集群的配置文件到 clusters/xxxx/config.yml；

ezdown 脚本统一设置组件、镜像版本；自动替换clusters/xxxx/config.yml 全局配置中相关版本

### 增加默认部署 node local dns
NodeLocal DNSCache在集群的上运行一个dnsCache daemonset来提高clusterDNS性能和可靠性。在K8S集群上的一些测试表明：相比于纯coredns方案，nodelocaldns + coredns方案能够大幅降低DNS查询timeout的频次，提升服务稳定性

参考官方文档：https://kubernetes.io/docs/tasks/administer-cluster/nodelocaldns/

### 客户端 kubeconfig 管理【强烈推荐】

经常遇到有人问某个kubeconfig(kubectl)泄露了怎么办？不同权限的kubeconfig怎么生成？这里利用cfssl签发自定义用户证书和k8s灵活的rbac权限绑定机制，ezctl 命令行封装了这个功能，非常方便、实用。

详细使用参考 docs/op/kcfg-adm.md 

### 更新 prometheus安装部署，自动集成安装

参考 example/config.yml 配置和 roles/cluster-addon/templates/prometheus/values.yaml.j2 模版配置文件，详细使用文档待更新

### 其他主要更新

- 更新支持 ansible 2.10.4
- 更新系统加固 os-harden 7.0.0
- 更新traefik 安装部署（helm charts）

### 组件更新

- k8s: v1.20.2, v.1.19.7, v1.18.15, v1.17.17

### 集群安装

- ca 安全管理，所有证书都在deploy节点创建后推送到需要的节点
- 移除 netaddr (pip安装) 依赖
- 修复ansible group命名不规范问题（group 'kube-node' --> group 'kube_node'）
- 更新 kube-ovn to 1.5.3 #958
- 调整cluster-addon安装方式
- 修复 calico 网络 backend 设置为 vxlan none 时，calico 部署失败 #959
- 调整默认nodePort范围为30000-32767
- 修复 calico backend config #973
- 修复 restore an etcd cluster #973
- 修复带自定义变量时增加/删除节点可能失败

### 工具脚本

- ezdown 替换原 tools/easzup
- ezctl 替换原 tools/easzctl

### 文档

- 大量文档更新（部分未完成）
