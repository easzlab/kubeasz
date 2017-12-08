## 部署 kubedns

kubedns 是 k8s 集群首先需要部署的，集群中的其他 pods 使用它提供域名解析服务；主要可以解析 `集群服务名` 和 `Pod hostname`；

配置文件参考 `https://github.com/kubernetes/kubernetes` 项目目录 `kubernetes/cluster/addons/dns` 

### 安装

kubectl create -f /etc/ansible/manifests/kubedns/[kubedns.yaml](../../manifests/kubedns/kubedns.yaml)




