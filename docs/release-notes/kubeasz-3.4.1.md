## kubeasz 3.4.1 (Frost's Descent)

霜降水返壑，风落木归山。冉冉岁将宴，物皆复本源。kubeasz 3.4.1 发布更新支持多架构平台（amd64/arm64）

### 1.更新多架构支持

当前已支持linux amd64和linux arm64，更多架构支持根据后续需求来计划。

目前多架构安装逻辑：根据部署机器（执行ezdown/ezctl命令的机器）的架构，会自动判断下载对应amd64/arm64的二进制文件和容器镜像，然后推送安装到整个集群。

- 暂不支持不同架构的机器加入到同一个集群。
- harbor目前仅支持amd64安装

### 2.重写项目依赖组件的镜像构建流程，利用github-action自动构建、推送多架构的镜像

k8s核心组件本身提供多架构的二进制文件/容器镜像下载，项目调整了下载二进制文件的容器dockerfile

- https://github.com/easzlab/dockerfile-kubeasz-k8s-bin

kubeasz其他用到的二进制或镜像，重新调整了容器创建dockerfile

- https://github.com/easzlab/dockerfile-kubeasz-ext-bin
- https://github.com/easzlab/dockerfile-kubeasz-ext-build
- https://github.com/easzlab/dockerfile-kubeasz-sys-pkg
- https://github.com/easzlab/dockerfile-kubeasz-mirrored-images
- https://github.com/easzlab/dockerfile-kubeasz
- https://github.com/easzlab/dockerfile-ansible

### 3.去除master/node节点上的admin kubeconfig文件，这个文件拥有全部集群权限，需要谨慎保管，目前仅部署机器上保留，可以自行按需管理使用。

### 组件更新

- k8s: v1.25.3

### 其他

- fix: curl dns resolving problem in a rare case (#ab9603d509900919)
- cleaning some images/pics
- fix: logo url
