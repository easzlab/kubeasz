## kubeasz-2.0.2 发布说明

**WARNNING:** 从 kubeasz 2.0.1 版本开始，项目仅支持 kubernetes 社区最近的4个大版本，当前为 v1.12/v1.13/v1.14/v1.15，更老的版本不保证兼容性。

CHANGELOG:
- 组件更新
  - docker: 18.09.7 
- 集群安装
  - **系统软件离线安装** 全面测试支持 Ubuntu1604/1804 CentOS7 Debian9/10 操作系统
  - kubelet 分离配置文件至 /var/lib/kubelet/config.yaml
  - containerd/docker 增加配置项是否启用容器仓库镜像
  - 修复 helm 安装时使用已有 namespace 执行报错
  - 调整部分基础软件安装
  - 调整 apiserver 部分参数 0ca5f7fdd9dc97c72ac
  - 调整清理脚本不再进行虚拟网卡、路由表、iptalbes/ipvs规则等清理，并提示清理脚本执行后重启节点
- easzup 工具
  - 添加配置项是否启用docker仓库CN镜像和选择合适的docker二进制下载链接
  - 修复docker已经安装时运行失败问题
  - update versions and minor fixes
- 文档
  - 离线安装文档更新
  - 集群安装相关文档更新
- 其他
  - new logo
  - fix: 执行roles/cluster-storage/cluster-storage.yml 报错不存在`deploy`
  - fix: 部分os启用kube-reserved出错（提示/sys/fs/cgroup只读）
  - fix: ex-lb 组少量 keepalived 相关配置
  - fix: 偶然出现docker安装时提示找不到变量`docker_ver`
  - fix: Ubuntu1804 pod内dns解析不到外网
  - fix: k8s 相关服务在接收SIGPIPE信号停止后不重启问题 #631 thx to gj19910723
