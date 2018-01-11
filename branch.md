## 项目分支说明

目前项目分支为 `master` `v1.9` `v1.8`，说明如下：

- `master` 分支将尽量使用最新版k8s和相关组件，网络使用`DaemonSet Pod`方式安装，目前提供`calico` `flannel` 可选
- `v1.9` 分支将尽量使用k8s v1.9的最新小版本和相关组件，使用`systemd service`方式安装 `calico`网络
- `v1.8` 分支将尽量使用k8s v1.8的最新小版本和相关组件，使用`systemd service`方式安装 `calico`网络
