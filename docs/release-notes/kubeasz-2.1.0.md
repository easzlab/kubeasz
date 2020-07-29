## kubeasz-2.1.0 发布说明

【Warnning】PROXY_MODE变量定义转移到ansible hosts #688，对于已有的ansible hosts文件需要手动增加该定义，参考example/hosts.*

CHANGELOG:
- 组件更新
  - k8s: v1.16.2 v1.15.5 v1.14.8 v1.13.12
  - docker: 18.09.9
  - coredns v1.6.2
  - metrics-server v0.3.6
  - kube-ovn: 0.8.0 #708
  - dashboard v2.0.0-beta5
- 集群安装
  - 更新/清理 APIs version，支持 k8s v1.16
  - 增加临时启停集群脚本 91.start.yml 92.stop.yml
  - 更新只读权限 read rbac role
- 工具脚本
  - 更新 tools/easzup 
- 文档
  - 增加go web应用部署实践 docs/practice/go_web_app
  - 增加go项目dockerfile示例 docs/practice/go_web_app/Dockerfile-more
  - 更新 log-pilot 日志方案 docs/guide/log-pilot.md
  - 更新主页【推荐工具栏】kuboard k9s octant
- 其他
  - fix: 增加kube-proxy参数--cluster-cidr #663
  - fix: 删除etcd服务不影响node服务 #690
  - fix: deploy阶段pip安装netaddr包
  - fix: 仅非容器化运行ansible需要安装 #658
  - fix: ipvs-connection-timeout-issue 
  - fix: heapster无法读取节点度量数据
  - fix: tcp_tw_recycle settings issue #714
  - fix: 文档文字“登陆”->“登录” #720
