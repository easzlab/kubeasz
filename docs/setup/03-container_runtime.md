# 03-安装容器运行时

项目根据k8s版本提供不同的默认容器运行时：

- k8s 版本 < 1.24 时，支持docker containerd 可选
- k8s 版本 >= 1.24 时，仅支持 containerd

## 安装containerd

作为 CNCF 毕业项目，containerd 致力于提供简洁、可靠、可扩展的容器运行时；它被设计用来集成到 kubernetes 等系统使用，而不是像 docker 那样独立使用。

- 安装指南 https://github.com/containerd/cri/blob/master/docs/installation.md
- 客户端 circtl 使用指南 https://github.com/containerd/cri/blob/master/docs/crictl.md
- man 文档 https://github.com/containerd/containerd/tree/master/docs/man

## kubeasz 集成安装 containerd

- 注意：k8s 1.24以后，项目已经设置默认容器运行时为 containerd，无需手动修改
- 执行安装：分步安装`ezctl setup xxxx 03`，一键安装`ezctl setup xxxx all`

## 命令对比

|命令           |docker         |crictl（推荐） |ctr                    |
|:-             |:-             |:-             |:-                     |
|查看容器列表   |docker ps      |crictl ps      |ctr -n k8s.io c ls     |
|查看容器详情   |docker inspect |crictl inspect |ctr -n k8s.io c info   |
|查看容器日志   |docker logs    |crictl logs    |无                     |
|容器内执行命令 |docker exec    |crictl exec    |无                     |
|挂载容器       |docker attach  |crictl attach  |无                     |
|容器资源使用   |docker stats   |crictl stats   |无                     |
|创建容器       |docker create  |crictl create  |ctr -n k8s.io c create |
|启动容器       |docker start   |crictl start   |ctr -n k8s.io run      |
|停止容器       |docker stop    |crictl stop    |无                     |
|删除容器       |docker rm      |crictl rm      |ctr -n k8s.io c del    |
|查看镜像列表   |docker images  |crictl images  |ctr -n k8s.io i ls     |
|查看镜像详情   |docker inspect |crictl inspecti|无                     |
|拉取镜像       |docker pull    |crictl pull    |ctr -n k8s.io i pull   |
|推送镜像       |docker push    |无             |ctr -n k8s.io i push   |
|删除镜像       |docker rmi     |crictl rmi     |ctr -n k8s.io i rm     |
|查看Pod列表    |无             |crictl pods    |无                     |
|查看Pod详情    |无             |crictl inspectp|无                     |
|启动Pod        |无             |crictl runp    |无                     |
|停止Pod        |无             |crictl stopp   |无                     |


[后一篇](04-install_kube_master.md)
