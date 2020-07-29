## kubeasz-2.0.1 发布说明

**WARNNING:** 从 kubeasz 2.0.1 版本开始，项目仅支持 kubernetes 社区最近的4个大版本，当前为 v1.12/v1.13/v1.14/v1.15，更老的版本不保证兼容性。

CHANGELOG:
- 组件更新
  - k8s: v1.15.0
  - metrics-server: v0.3.3
- 集群安装
  - **系统软件离线安装** 支持 chrony/ipvsadm/ipset/haproxy/keepalived 等系统软件，目前已测试 Ubuntu1604/Ubuntu1804/CentOS7 操作系统
  - **修复及简化** 集群备份/恢复脚本及文档 
  - 调整 kubelet 默认禁用`--system-reserved`参数配置
  - 修复 kubelet v1.15 删除参数`--allow-privileged`
  - 修复升级集群时重新配置 k8s 服务文件
- easzup 工具
  - 增加自动下载系统软件包
  - 增加离线保存及加载kubeasz镜像
  - 修复下载docker安装包的位置
  - 增加预处理安装前准备ssh key pair和确保$PATH下python可执行
- 文档
  - 增加**完全离线安装**集群文档
  - 添加非标ssh端口节点文档说明 docs/op/op-node.md 
- 其他
  - role:deploy 增加是否容器化运行ansible脚本的判断
  - fix: 容器化运行deploy任务时删除kubeconfig报错
  - fix: 节点做网卡bonding时获取host_ip错误（thx beef9999, ISSUE #607）
  - fix: 节点操作系统ulimit设置
