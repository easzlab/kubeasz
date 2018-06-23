# Helm

`Helm`致力于成为k8s集群的应用包管理工具，希望像linux 系统的`RPM` `DPKG`那样成功；确实在k8s上部署复杂一点的应用很麻烦，需要管理很多yaml文件（configmap,controller,service,rbac,pv,pvc等等），而helm能够整齐管理这些文档：版本控制，参数化安装，方便的打包与分享等。  
- 建议积累一定k8s经验以后再去使用helm；对于初学者来说手工去配置那些yaml文件对于快速学习k8s的设计理念和运行原理非常有帮助，而不是直接去使用helm，面对又一层封装与复杂度。

## 安装（开发测试环境）

在开发环境安装（不考虑安全性问题）很方便，除了tiller官方镜像需要翻墙下载。  
- 客户端：从最新[release](https://github.com/kubernetes/helm/releases)，下载helm-v2.9.1-linux-amd64.tar.gz到本地，解压后把二进制直接放到环境PATH下即可
- 服务端：不能翻墙的可以使用docker hub上的转存镜像，比如jmgao1983/tiller:v2.9.1  
helm默认使用kubectl客户端相同的配置文件去访问k8s集群，因此只要在能使用kubectl的节点运行如下，即能进行安装。  
``` bash
$ helm init --tiller-image jmgao1983/tiller:v2.9.1
```
- 验证  
``` bash
$ kubectl get pod --all-namespaces|grep tiller
kube-system   tiller-deploy-7c6cd89d69-72r7j            1/1       Running   0          10h

$ helm version
Client: &version.Version{SemVer:"v2.9.1", GitCommit:"20adb27c7c5868466912eebdf6664e7390ebe710", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.9.1", GitCommit:"20adb27c7c5868466912eebdf6664e7390ebe710", GitTreeState:"clean"}
```

如果 `helm version` 出现如下错误，在每个节点安装 `socat`即可（如 apt install socat）  
``` bash
E0522 22:22:15.492436   24409 portforward.go:331] an error occurred forwarding 38398 -> 44134: error forwarding port 44134 to pod dc6da4ab99ad9c497c0cef1776b9dd18e0a612d507e2746ed63d36ef40f30174, uid : unable to do port forwarding: socat not found.
Error: cannot connect to Tiller
```

## 安全安装 helm

上述安装的tiller服务器默认允许匿名访问，那么k8s集群中的任何pod都能访问tiller，风险较大，因此需要在helm客户端和tiller服务器间建立安全的SSL/TLS认证机制；tiller服务器和helm客户端都是使用同一CA签发的`client cert`，然后互相识别对方身份。建议通过本项目提供的`ansible role`安装，符合官网上介绍的安全加固措施，在delpoy节点运行:  
``` bash
# 1.如果已安装非安全模式，使用 helm reset 清理
# 2.配置默认helm参数 vi /etc/ansible/roles/helm/defaults/main.yml
# 3.执行安装
$ ansible-playbook /etc/ansible/roles/helm/helm.yml
```

简单介绍下`/roles/helm/tasks/main.yml`中的步骤

- 1-下载最新release的helm客户端到/etc/ansible/bin目录下，再由它自动推送到deploy的{{ bin_dir }}目录下
- 2-由集群CA签发helm客户端证书和私钥
- 3-由集群CA签发tiller服务端证书和私钥
- 4-创建tiller专用的RBAC配置，只允许helm在指定的namespace查看和安装应用
- 5-安全安装tiller到集群，tiller服务启用tls验证
- 6-配置helm客户端使用tls方式与tiller服务端通讯
- 7-创建helms命令别名，方便使用，即alias helms='helm --tls --tiller-namespace {{ helm_namespace }}'

注：helms别名生效请执行：`source ~/.bashrc`，或者退出后重新登陆shell  
- 使用`helms`执行与tiller服务有关的命令，比如 `helms ls` `helms version` `helms install`等
- 使用`helm`执行其他命令，比如`helm search` `helm fetch` `helm home`等

## 使用helm安装应用到k8s上

请阅读本项目文档[helm安装prometheus监控](prometheus.md)
