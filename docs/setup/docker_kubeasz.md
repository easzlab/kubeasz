# 容器化运行 kubeasz

## easzup 工具介绍

初始化工具 tools/easzup 主要用于：

- 下载 kubeasz 项目代码/k8s 二进制文件/其他所需二进制文件/离线docker镜像等
- 【可选】容器化运行 kubeasz

详见脚本内容

### 容器化运行 kubeasz

kubeasz 容器启动脚本详见文件 tools/easzup 中函数`start_kubeasz_docker`

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
- --volume /root/.ssh/id_rsa:/root/.ssh/id_rsa:ro 等三个 volume 挂载保证：如果宿主机配置了免密码登录所有集群节点，那么容器内也可以免密码登录所有节点

## 容器化安装集群

项目[快速指南](quickStart.md)，就是利用 kubeasz 容器快速安装单节点k8s集群的例子。

## 验证

使用容器化安装成功后，可以在 **容器内** 或者 **宿主机** 上执行 kubectl 命令验证集群状态。

## 清理

登录管理节点，按照如下步骤清理（清理后可以重新安装测试）

- 1.清理集群 `$ docker exec -it kubeasz easzctl destroy`
- 2.清理管理节点
  - 清理运行的容器 `$ easzup -C`
  - 清理容器镜像 `$ docker system prune -a`
  - 停止docker服务 `$ systemctl stop docker`
  - 删除下载文件 `$ rm -rf /etc/ansible /etc/docker /opt/kube`
  - 删除docker文件 
```
$ umount /var/run/docker/netns/default
$ umount /var/lib/docker/overlay
$ rm -rf /var/lib/docker /var/run/docker
```
