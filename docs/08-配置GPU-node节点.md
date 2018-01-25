## 08-配置GPU-node节点.md

推荐阅读[官方GPU节点配置文档](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/)

### 1. 允许k8s系统使用Device Plugins

实现方式：通过修改kube-apiserver, kubelet, kube-proxy配置文件模板

因为GPU node是kube node的子集，正常执行`90.setup.yml`搭建k8s即可实现此目标

### 2. 配置GPU节点
在[官方驱动网页](http://www.nvidia.com/Download/index.aspx?lang=en-uk)下载对应操作系统与显卡的驱动，
改名为nvidia-diag-driver-local-repo.deb，
并放入到invertory file中{{ base_dir }}路线下的 bin 文件夹中

执行命令：`ansible-playbook 21.gpunode.yml`

相当于完成以下任务

#### 1). GPU 节点安装 nvidia driver 
建议在node上只安装驱动，而不安装CUDA包，所有CUDA包都放到镜像中去，否则容易出现版本不匹配的问题。



#### 2). GPU 节点安装 nvidia docker 2.0
这一部分可能失败，原因可能是系统所用docker版本与nvidia docker 2.0依赖的docker版本不一致，
具体参考[文档](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions)

    
可以通过命令`apt-cache madison nvidia-docker2 nvidia-container-runtime`查询可以安装的nvidia docker 2.0 
的 docker 版本
    
可以通过命令`dpkg -l | grep docker`来查看实际安装docker 版本
    
nvidia docker 安装失败，可以[手动更新docker](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/#upgrade-docker-ce),
再重新运行安装playbook


#### 3). GPU 节点配置 nvidia-container-runtime 为 docker default runtime

通过修改`/etc/docker/daemon.json`实现

### 3. 配置Nvidia device plugin
执行命令 `kubectl create -f manifests/gpu-device-plugin/v1.9/nvidia-device-plugin.yml`

可以通过执行`kubectl describe nodes | grep nvidia.com/gpu` 来查看GPU节点配置是否成功