# 容器化运行 kubeasz

## TL;DR;

- 1.本机安装 docker (略) 
- 2.配置 ssh 密钥登陆集群节点

``` bash
ssh-keygen -t rsa -b 2048 回车 回车 回车
ssh-copy-id $IP  # $IP 为所有节点地址包括自身，按照提示输入 yes 和 root 密码
```

- 3.下载 kubeasz docker 镜像并运行 (可能需较长时间下载镜像jmgao1983/kubeasz)

``` bash
curl -sfL https://github.com/gjmzj/kubeasz/releases/download/1.0.0/kubeasz-docker-1.0.0 | bash -
```

- 4.在 kubeasz 容器中创建 k8s 集群，步骤与非容器方式创建类似，快速创建单节点集群如下

``` bash
docker exec -it kubeasz easzctl start-aio
```

## 验证

使用容器化安装成功后，可以在 **容器内** 或者 **宿主机** 上执行 kubectl 命令验证集群状态。

## kubeasz 镜像介绍

镜像描述文件 dockerfiles/kubeasz/Dockerfile，它基于 ansible 镜像（dockerfiles/ansible/Dockerfile），主要包含 kubeasz 项目代码和 k8s 集群安装所需二进制文件。

- 在本地创建 kubeasz 镜像，由于镜像较大，可以按以下步骤在本地创建

``` bash
cd /etc/ansible/dockerfiles/kubeasz
# 克隆代码
git clone --depth=1 https://github.com/gjmzj/kubeasz.git
# 手动下载二进制文件放入上述 git clone 完成目录 kubeasz/bin 
docker build -t kubeasz:$TAG .
```

## 容器运行讲解

容器启动脚本详见文件 tools/kubeasz_docker

``` bash
docker run --detach \
      --name kubeasz \
      --restart always \
      --env HOST_IP=$host_ip \
      --volume /etc/ansible:/etc/ansible \
      --volume /root/.kube:/root/.kube \
      --volume /root/.ssh/id_rsa:/root/.ssh/id_rsa:ro \
      --volume /root/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub:ro \
      --volume /root/.ssh/known_hosts:/root/.ssh/known_hosts:ro \
      $KUBEASZ_DOCKER_VER
```

- --env HOST_IP=$host_ip 传递这个参数是为了快速在本机安装aio集群
- --volume /etc/ansible:/etc/ansible 挂载本地目录，这样可以在宿主机上修改集群配置，然后在容器内执行 ansible 安装
- --volume /root/.kube:/root/.kube 容器内与主机共享 kubeconfig，这样都可以执行 kubectl 命令
- --volume /root/.ssh/id_rsa:/root/.ssh/id_rsa:ro 等三个 volume 挂载保证：如果宿主机配置了免密码登陆所有集群节点，那么容器内也可以免密码登陆所有节点

## 参考

- ansible 容器镜像制作: https://github.com/William-Yeh/docker-ansible
