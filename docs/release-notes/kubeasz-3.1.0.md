## kubeasz-3.1.0 (Grain Rain)

春归谷雨，夏至未至。kubeasz 3.1.0 版本发布，主要更新点：

### 正式通过k8s一致性认证

kubeasz 用户可以确认集群各项功能符合预期，既符合k8s设计标准；

- v1.20 [已完成](https://github.com/cncf/k8s-conformance/pull/1326)

### 推荐群里大佬的k8s架构师免费视频课程

作者花很多心思和精力去构思文档、视频录制，并且把工作中的实践经验分享出来；值得参考学习

https://www.toutiao.com/c/user/token/MS4wLjABAAAA0YFomuMNm87NNysXeUsQdI0Tt3gOgz8WG_0B3MzxsmI/?tab=article

### 重写chrony/ex-lb/kube-lb等组件的安装

原先使用yum/apt方式安装依赖于各操作系统软件源，版本无法统一管理，并且离线安装也不方便；现使用源码最小化依赖编译安装，仅依赖基本库，生成的二进制文件可以运行于受支持的各种操作系统，可以方便的版本管理、配置管理和离线安装。

原node节点haproxy 由仅四层转发的nginx替代(kube-lb)，负责集群内部负载均衡apiservers；简化部署逻辑，现在每个节点均会运行一个轻量kube-lb进程。

原ex-lb组件keepalived+haproxy由 keepalived+l4lb替代，l4lb同样是仅支持四层转发的nginx源码编译的。

### 修改有条件使用`systemd` `cgroup driver`

当容器运行时选择containerd，或者docker version >= 20.10时，容器运行时和kubelet使用`systemd`做资源管理和限制，这是官方文档建议的方式，一定程度上能增加稳定性；
当选择docker version < 20.10时，使用`cgroupfs`；主要因为部分操作系统不支持dockerd使用cgroup=systemd，会提示报错：`OCI runtime create failed: systemd cgroup flag passed, but systemd support for managing cgroups is not available: unknown`。

### 组件更新

- k8s: v1.21.0, v1.20.6, v.1.19.10, v1.18.18
- containerd: v1.4.4 (runc: v1.0.0-rc93)
- coredns: 1.8.0
- dns-node-cache: 1.17.0
- pause: 3.4.1

### 其他

- fix:增加/删除节点时ansible hosts文件更新错误 
- fix:kube-scheduler healthz/metrics listening setting
- fix:restart ex-lb when master nodes change
- fix:多条默认路由网卡自动识别问题
- fix:安装aio集群时docker cgroupdriver设置问题
- fix:add scheme:https to service-account-issuer
- fix:容器化aio安装时选择containerd运行时失败
- feat:增加可选配置apiserver安全端口
- feat:允许修改配置ingress port #999
- feat:增加支持ubuntu 20.04
- feat:增加ezctl setup支持传入额外参数 #1007
- 更新ansible.cfg
- 更新get-pip.py下载地址 #1006
