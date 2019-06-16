# 容器化运行 kubeasz

## TL;DR;

- 1.准备一台全新虚机（ansible控制端）
```
$ curl -sfL https://github.com/easzlab/kubeasz/releases/download/1.3.0/easzup
$ ./easzup -D
``` 
- 2.配置 ssh 密钥登陆集群节点
``` bash
ssh-keygen -t rsa -b 2048 回车 回车 回车
ssh-copy-id $IP  # $IP 为所有节点地址包括自身，按照提示输入 yes 和 root 密码
```
- 3.容器化运行 kubeasz，然后执行安装 k8s 集群（举例aio集群）

``` bash
$ ./easzup -S
$ docker exec -it kubeasz easzctl start-aio
# 若需要自定义集群创建，如下进入容器，然后配置/etc/ansible/hosts，执行创建即可
# docker exec -it kubeasz sh
```

## 验证

使用容器化安装成功后，可以在 **容器内** 或者 **宿主机** 上执行 kubectl 命令验证集群状态。

## easzup 工具介绍

初始化工具 tools/easzup 主要用于：

- 下载 kubeasz 项目代码/k8s 二进制文件/其他所需二进制文件/离线docker镜像等
- 【可选】容器化运行 kubeasz

详见脚本内容

### 容器化运行 kubeasz

容器启动脚本详见文件 tools/easzup 中函数`start_kubeasz_docker`

``` bash
  docker run --detach \
      --name kubeasz \
      --restart always \
      --env HOST_IP="$host_ip" \
      --volume /etc/ansible:/etc/ansible \
      --volume /root/.kube:/root/.kube \
      --volume /root/.ssh/id_rsa:/root/.ssh/id_rsa:ro \
      --volume /root/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub:ro \
      --volume /root/.ssh/known_hosts:/root/.ssh/known_hosts:ro \
      easzlab/kubeasz:${KUBEASZ_VER}
```

- --env HOST_IP="$host_ip" 传递这个参数是为了快速在本机安装aio集群
- --volume /etc/ansible:/etc/ansible 挂载本地目录，这样可以在宿主机上修改集群配置，然后在容器内执行 ansible 安装
- --volume /root/.kube:/root/.kube 容器内与主机共享 kubeconfig，这样都可以执行 kubectl 命令
- --volume /root/.ssh/id_rsa:/root/.ssh/id_rsa:ro 等三个 volume 挂载保证：如果宿主机配置了免密码登陆所有集群节点，那么容器内也可以免密码登陆所有节点

