## kubeasz-0.3.1 发布说明

CHANGELOG:
- 组件更新：
  - k8s: v1.11.3, v1.10.7
  - kube-router: v0.2.0
  - dashboard: v1.10.0
  - docker: 17.03.2-ce (选择k8s官方测试稳定的版本)
- 集群安装：
  - **增加集群时间同步服务chrony** [说明](https://github.com/easzlab/kubeasz/blob/master/docs/guide/chrony.md)
  - **取消 Node节点 Bootstrap机制**，安装流程更稳定，配置更精简
  - 简化 ansible host 文件配置，移除etcd、harbor 相关变量
  - 拆分 prepare 阶段的安装脚本，增加设置系统 ulimit
  - 增加多lb节点（多于2节点）配置支持 (#286)
  - 增加可选配置lb 节点负载转发ingress controller NodePort service的功能 
  - 自定义 kubelet docker 存储目录 (#305)
  - 增加变量配置支持多网卡情况时安装 flannel calico 
- 文档更新：
  - 更新 kubeasz 公有云安装文档 https://github.com/easzlab/kubeasz/blob/master/docs/setup/kubeasz_on_public_cloud.md
  - 更新 java war应用部署实践 https://github.com/easzlab/kubeasz/blob/master/docs/practice/java_war_app.md
  - 更新 cilium 文档，翻译官方 cilium 安全策略例子（deathstar/starwar） 
  - 更新 harbor kubedns README 文档
  - 更新集群安装部分文档
- 其他：
  - 修复 calicoctl 配置，修复calico/node跑在LB 主节点时使用`vip`作为`bgp peer`地址问题
  - 修复 jq安装错误，补充ipset和ipvsadm安装
  - 修复清除单节点脚本 tools/clean_one_node.yml
  - 修复消除离线镜像不存在时安装的错误提示信息
  - 修复多节点（超过2节点时）lb 备节点 router_id重复问题
  - 锁定jenkins镜像tag、升级插件版本以及锁定安全插件 (#315)
