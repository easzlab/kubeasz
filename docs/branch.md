## 项目分支说明

目前项目分支为 `master` `v1.9` `v1.8`，说明如下：

- `master` 分支将尽量使用最新版k8s (v1.10)和相关组件，网络使用`DaemonSet Pod`方式安装，目前提供`calico` `flannel` 可选
- `v1.9` 分支将尽量使用k8s v1.9的最新小版本和相关组件，网络使用`DaemonSet Pod`方式安装，目前提供`calico` `flannel` 可选
- `v1.8` 分支将尽量使用k8s v1.8的最新小版本和相关组件，使用`systemd service`方式安装 `calico`网络

## 项目分支与百度网盘离线包关系

- `master` 分支请使用 `k8s.110x.tar.gz` 的安装包
- `v1.9` 分支请使用 `k8s.19x.tar.gz` 的安装包
- `v1.8` 分支请使用 `k8s.18x.tar.gz` 的安装包

## 更新频率和内容 

- `master` 更新频繁：**相关文档**，**功能特性**，BUG修复，组件更新
- `v1.9` `v1.8` 较少更新：BUG修复，组件更新
