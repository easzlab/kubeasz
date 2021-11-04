## kubeasz-3.0.1 (Pure Brightness)

清明降至，踏青郊游，祭祖缅怀。kubeasz 3.0.1 版本发布，主要更新点：

### 技术上通过k8s一致性认证的所有测试项

kubeasz 用户可以确认集群各项功能符合预期，既符合k8s设计标准；下一步会继续走流程正式申请成为官方认证的部署工具；正式PR在此：https://github.com/cncf/k8s-conformance/pull/1326。

### 推荐群里大佬的k8s架构师免费视频课程

作者花很多心思和精力去构思文档、视频录制，并且把工作中的实践经验分享出来；值得参考学习

https://www.toutiao.com/c/user/token/MS4wLjABAAAA0YFomuMNm87NNysXeUsQdI0Tt3gOgz8WG_0B3MzxsmI/?tab=article

### 更新harbor 安装流程 
重写 harbor 安装流程，利用easzlab/harbor-offline:v2.1.3 仓库加速离线安装文件下载，增加可选安装组件。


### 组件更新

- k8s: v1.20.5, v.1.19.9, v1.18.17
- docker: 20.10.5
- dashboard: v2.2.0
- harbor: v2.1.3

### 集群安装

- 修复默认集群内部dns域名后缀 
- 调整etcd集群配置参数
- 更新kube-scheduler部署使用配置文件 kube-scheduler-config.yaml
- 更新集群存储插件 nfs-provisioner
- 修复安装外部负载均衡服务 ./ezctl setup ${集群名} ex-lb
- 修复清理LB（haproxy/keepalived）服务可能报错问题
- 修复worker节点安装时无法推送dnscache镜像
