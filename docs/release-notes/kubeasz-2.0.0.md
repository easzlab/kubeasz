## kubeasz-2.0.0 发布说明

**IMPORTANT:** 本次更新为 HA-2x (#585) 架构的第一个版本，相比 HA-1x (#584) 架构优势在于：
- 高可用集群安装更简单，不再依赖外部负载均衡；自有环境和云厂商环境安装流程完全一致；
- 扩展集群更方便，从单节点集群可以方便地扩展到多主多节点集群；

**WARNNING:** 因为架构差异，已有集群（HA-1x）不支持升级到 2x 版本；只能在创建新 k8s 集群使用 2x 新版本；后续项目重心转移到维护 2x 版本，详见 [分支说明](https://github.com/easzlab/kubeasz/blob/master/docs/mixes/branch.md) 

CHANGELOG:
- 集群安装：
  - 废弃 ansible hosts 中 deploy 角色，精简保留2个预定义节点规划例子（example/hosts.xx）
  - 重构 prepare 安装流程（删除 deploy 角色，移除 lb 节点创建）
  - 调整 kube-master 安装流程
  - 调整 kube-node 安装流程（node 节点新增 haproxy 服务）
  - 调整 network 等其他安装流程
  - 精简 example hosts 配置文件及配置项
  - 调整 ex-lb 安装流程【可选】
  - 添加 docker/containerd 安装时互斥判断
  - 新增 role: clean，重写清理脚本 99.clean.yml
  - 废弃 tools/clean_one_node.yml
  - 调整 helm 安装流程
  - 调整 cluster-addon 安装流程（自动安装traefik，调整dashboard离线安装）
  - 替换 playbook 中 hosts: all 为具体节点组名称，防止操作扩大风险
  - 废弃百度盘下载方式，新增 easzup 下载工具
- easzctl 工具 
  - 废弃 clean-node 命令，调整为 del-master/del-node 命令
  - 调整 add-etcd/add-node/add-master 脚本以适应 HA-2x 架构
  - 调整 del-etcd/del-node/del-master 脚本
  - 修复 add-node/add-master/add-etcd 判断节点是否存在
- easzup 工具
  - 修复 centos 等可能存在 selinux 设置问题
  - 下载 docker 二进制时使用 curl 替换 wget
- 文档：
  - 集群安装相关大量文档更新
    - 快速指南安装文档
    - 集群规划与配置介绍
    - 公有云安装文档
    - node 节点安装文档
    - ...
  - 集群操作管理文档更新（docs/op/op-index.md）
  - 新增可选外部负载均衡文档（docs/setup/ex-lb.md）
  - 新增容器化系统服务 haproxy/chrony 文档（docs/practice/dockerize_system_service.md）
- 其他：
  - fix: 对已有集群进行安全加固时禁用 ip_forward 问题
  - fix: haproxy 最大连接数设置
  - fix: 容器化运行 kubeasz 时清理脚本
