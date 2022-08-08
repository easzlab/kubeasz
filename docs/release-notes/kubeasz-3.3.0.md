## kubeasz 3.3.0 (Grain in Ear)

泽草所生，种之芒种。kubeasz 3.3.0 发布更新支持k8s 1.24 版本。

### 主要更新

#### 1.容器运行时

修改默认容器运行时为containerd，kubeasz 3.3.0 暂未适配docker 和其他容器运行时；集群使用containerd作为运行时，确实更简单、稳定；至于docker，镜像打包、单机运行容器等等真好用；各自发挥所长吧；kubeasz项目中在离线资源下载、安装中使用docker非常方便，还可以避免在部署机器上安装ansible等麻烦事，推荐使用。

#### 2.去除安装ingress插件

ingress一般是具体业务强相关的，属于上层组件；鉴于维护人力和频率，项目中仅保留历史相关ingress文档，不再继续更新；请移步相关ingress组件官网获取更新部署方式；kubeasz 今后将更加关注底层集群组件的更新和维护。做简单，做好一件事。

#### 3.更新prometheus安装套件

监控组件属于底层功能，将持续更新；项目使用kube-prometheus-stack helm chart 默认部署，需要自定义设置请参考项目 https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack

### 组件更新

- k8s: v1.24.1
- etcd: v3.5.4
- containerd: 1.6.4
- calico: v3.19.4
- cni-plugins: v1.1.1
- dashboard: v2.5.1

### 其他

- 调整kube-controller-manager启动配置文件
- 调整kubelet启用配置文件
- 修复部分系统首次执行安装失败 (缺失 '/usr/bin/python')
- 修复'ezdown‘运行可能会遗留容器导致再次运行失败
- 部分文档更新
- fix: get secret tokens for dashboard login in v1.24
