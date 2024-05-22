## kubeasz 3.6.4

kubeasz 3.6.4 发布：支持k8s v1.30版本，组件更新和一些bugfix。

### 版本更新

- k8s: v1.30.1
- etcd: v3.5.12
- containerd: 1.7.17
- runc: v1.1.12
- calico: v3.26.4
- cilium: 1.15.5
- cni: v1.4.1
- harbor: v2.10.2
- metrics-server: v0.7.1

### 重要更新

- 安全更新：to solve CVE-2024-21626: update containerd, runc
- 安装流程：role 'prepare' 阶段增加设置hostname，这样当网络组件为calico时不会因为主机名相同而出错；同时在example/config.yml 中增加配置开关`ENABLE_SETTING_HOSTNAME`
- 操作系统：增加测试支持 Ubuntu 2404
  - 已知在ubuntu 2404上使用网络插件calico v3.26.4不兼容，提示：ipset v7.11: Kernel and userspace incompatible
  - 使用cilium 组件没有问题

### 其他

- 21376465de7f44d1ec997bde096afc7404ce45c5 fix: cilium ui images settings
- c40548e0e33cab3c4e5742aacce11101ac0c7366 #1343, 恢复podPidsLimit=-1默认设置
- 
