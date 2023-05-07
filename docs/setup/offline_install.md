# 离线安装集群

使用kubeasz 离线安装 k8s集群需要下载四个部分：

- kubeasz 项目代码
- 二进制文件（k8s、etcd、containerd等组件）
- 容器镜像文件（calico、coredns、metrics-server等容器镜像）
- 系统软件安装包（ipset、libseccomp2等，仅无法使用本地yum/apt源时需要）

## 离线文件准备

在一台能够访问互联网的服务器上执行：

- 下载工具脚本ezdown，举例使用kubeasz版本3.6.0

``` bash
export release=3.6.0
wget https://github.com/easzlab/kubeasz/releases/download/${release}/ezdown
chmod +x ./ezdown
```

- 使用工具脚本下载（更多关于ezdown的参数，运行./ezdown 查看）

下载kubeasz代码、二进制、默认容器镜像

``` bash
# 国内环境
./ezdown -D
```

[可选]如果需要更多组件，请下载额外容器镜像（cilium,flannel,prometheus等）

``` bash
./ezdown -X
```

下载离线系统包 (适用于无法使用yum/apt仓库情形)

``` bash
# 如果操作系统是ubuntu 22.04
./ezdown -P ubuntu_22
```

上述脚本运行成功后，所有文件（kubeasz代码、二进制、离线镜像）均已整理好放入目录`/etc/kubeasz`

- `/etc/kubeasz` 包含 kubeasz 版本为 ${release} 的发布代码
- `/etc/kubeasz/bin` 包含 k8s/etcd/docker/cni 等二进制文件
- `/etc/kubeasz/down` 包含集群安装时需要的离线容器镜像
- `/etc/kubeasz/down/packages` 包含集群安装时需要的系统基础软件

## 离线安装

上述下载完成后，把`/etc/kubeasz`整个目录复制到目标离线服务器相同目录，然后在离线服务器/etc/kubeasz目录下执行：

- 离线安装 docker，检查本地文件，正常会提示所有文件已经下载完成，并上传到本地私有镜像仓库

```
./ezdown -D
./ezdown -X
```

- 启动 kubeasz 容器

```
./ezdown -S
```

- 设置参数允许离线安装系统软件包

```
sed -i 's/^INSTALL_SOURCE.*$/INSTALL_SOURCE: "offline"/g' /etc/kubeasz/example/config.yml 
```

- 举例安装单节点集群，参考 https://github.com/easzlab/kubeasz/blob/master/docs/setup/quickStart.md

``` bash
source ~/.bashrc
dk ezctl start-aio
# 或者执行 docker exec -it kubeasz ezctl start-aio
```

- 多节点集群，进入kubeasz 容器内 `docker exec -it kubeasz bash`，参考https://github.com/easzlab/kubeasz/blob/master/docs/setup/00-planning_and_overall_intro.md 进行集群规划和设置后使用./ezctl 命令安装

