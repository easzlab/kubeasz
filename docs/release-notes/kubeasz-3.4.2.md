## kubeasz 3.4.2 (Great Snow)

晚来天欲雪，能饮一杯无。kubeasz 3.4.2 发布，小版本更新以及一些bugfix。

### 小版本更新

- k8s: v1.25.4
- etcd: v3.5.5

### 更新国内容器镜像站

mirrors.ustc.edu.cn 站点已经停止服务，当前替换为docker.nju.edu.cn，提升国内网络环境下载国外容器镜像的速度。

### 新增命令强制更新集群CA及其他证书

此命令使用需要小心谨慎，确保了解功能背景和可能的结果；执行后，它会重新创建集群CA证书以及由它颁发的所有其他证书；一般适合于集群admin.conf不小心泄露，为了避免集群被非法访问，重新创建CA，从而使已泄漏的admin.conf失效。更新过程中会中断集群服务，详见 docs/op/force_ch_certs.md 使用说明。

### 新增机器人清理过期的issue

### 优化etcd 备份和恢复流程，by itswl (#1191 #1193)

### 修复默认配置的 Kubernetes CA 证书 by ffutop (#1197)

### 调整cluster-addon 组件安装流程

### 其他

- fix：系统架构判断，replace 'uname -p(non-portable)' to 'uname -m'
- fix：/var/lib/etcd单独分区时，删除集群报"Device or resource busy"错误 #1159
- fix：更新 roles/kube-master/main.yml 修改证书时复制新证书到Master节点 by liyu36 (#1186)
- fix：kube-apiserver 访问 kubelet的权限
- fix：shell加载环境变量 by itswl (#1202 #1203)
- fix：离线安装系统软件包 (38925ccc56134e4d007fec2a71691828dd15e9d5)
