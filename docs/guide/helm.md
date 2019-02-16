# Helm

`Helm`致力于成为k8s集群的应用包管理工具，希望像linux 系统的`RPM` `DPKG`那样成功；确实在k8s上部署复杂一点的应用很麻烦，需要管理很多yaml文件（configmap,controller,service,rbac,pv,pvc等等），而helm能够整齐管理这些文档：版本控制，参数化安装，方便的打包与分享等。  
- 建议积累一定k8s经验以后再去使用helm；对于初学者来说手工去配置那些yaml文件对于快速学习k8s的设计理念和运行原理非常有帮助，而不是直接去使用helm，面对又一层封装与复杂度。
- 本文参考 helm 官网安全实践启用 TLS 认证，参考 https://docs.helm.sh/using_helm/#securing-your-helm-installation 

## 安全安装 helm（在线）

在helm客户端和tiller服务器间建立安全的SSL/TLS认证机制；tiller服务器和helm客户端都是使用同一CA签发的`client cert`，然后互相识别对方身份。建议通过本项目提供的`ansible role`安装，符合官网上介绍的安全加固措施，在delpoy节点运行:  
``` bash
# 1.如果已安装非安全模式，使用 helm reset 清理
# 2.配置默认helm参数 vi  /etc/ansible/roles/helm/defaults/main.yml
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

### 注意因使用了TLS认证，所以helm命令执行分以下两种情况 

- 执行与tiller服务有关的命令，比如 `helm ls` `helm version` `helm install`等需要加`--tls`参数
- 执行其他命令，比如`helm search` `helm fetch` `helm home`等不需要加`--tls`
- helm v2.11.0及以上版本，启用环境变量 HELM_TLS_ENABLE=true，可以都不用加 --tls 参数

## 安全安装 helm（离线）
在内网环境中，由于不能访问互联网，无法连接repo地址，使用上述的在线安装helm的方式会报错。因此需要使用离线安装的方法来安装。
离线安装步骤：
```bash
# 1.创建本地repo
mkdir -p /opt/helm-repo
# 2.启动helm repo server,如果要其他服务器访问，改为本地IP
nohup helm serve --address 127.0.0.1:8879 --repo-path /opt/helm-repo &
# 3.更改helm 配置文件
将/etc/ansible/role/helm/default/main.yml中repo的地址改为 http://127.0.0.1:8879
cat <<EOF >/etc/ansible/role/helm/default/main.yml
helm_namespace: kube-system 
helm_cert_cn: helm001
tiller_sa: tiller
tiller_cert_cn: tiller001
tiller_image: jmgao1983/tiller:v2.9.1
#repo_url: https://kubernetes-charts.storage.googleapis.com
repo_url: http://127.0.0.1:8879
# 如果默认官方repo 网络访问不稳定可以使用如下的阿里云镜像repo
#repo_url: https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
EOF
# 4.运行安全helm命令
ansible-playbook /etc/ansible/role/helm/helm.yml 
```
## 使用helm安装应用到k8s上

请阅读本项目文档[helm安装prometheus监控](prometheus.md)
