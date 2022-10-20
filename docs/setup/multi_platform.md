# 多架构支持

kubeasz 3.4.1 以后支持多CPU架构，当前已支持linux amd64和linux arm64，更多架构支持根据后续需求来计划。

## 使用方式

kubeasz 多架构安装逻辑：根据部署机器（执行ezdown/ezctl命令的机器）的架构，会自动判断下载对应amd64/arm64的二进制文件和容器镜像，然后推送安装到整个集群。

- 暂不支持不同架构的机器加入到同一个集群。
- harbor目前仅支持amd64安装

## 架构支持备忘

#### k8s核心组件本身提供多架构的二进制文件/容器镜像下载，项目调整了下载二进制文件的容器dockerfile

- https://github.com/easzlab/dockerfile-kubeasz-k8s-bin

#### kubeasz其他用到的二进制或镜像，重新调整了容器创建dockerfile

- https://github.com/easzlab/dockerfile-kubeasz-ext-bin
- https://github.com/easzlab/dockerfile-kubeasz-ext-build
- https://github.com/easzlab/dockerfile-kubeasz-sys-pkg
- https://github.com/easzlab/dockerfile-kubeasz-mirrored-images
- https://github.com/easzlab/dockerfile-kubeasz
- https://github.com/easzlab/dockerfile-ansible

#### 其他组件(coredns/network plugin/dashboard/metrics-server等)一般都提供多架构的容器镜像，可以直接下载拉取


