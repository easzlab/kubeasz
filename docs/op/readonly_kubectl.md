# 配置 kubectl 只读访问权限

默认 k8s 集群安装后配置的 kubectl 客户端拥有所有的管理权限，而有时候我们需要把只读权限分发给普通开发人员，本文档将创建一个只读权限的kubectl 配置文档 kubeconfig。

## 创建

- 备份下原先 admin 权限的 kubeconfig 文件：`mv ~/.kube ~/.kubeadmin`
- 执行 `ansible-playbook /etc/ansible/01.prepare.yml -t create_kctl_cfg -e USER_NAME=read`，成功后查看~/.kube/config 即为只读权限

## 讲解

创建主要包括三个步骤：

- 创建 group:read rbac 权限
- 创建 read 用户证书和私钥
- 创建 kubeconfig

### read rbac 权限

所有权限控制魔法在`k8s`中由`rbac`实现，所谓`read`权限类似于集群自带的`clusterrole view`，具体查看：

`kubectl get clusterrole view -o yaml`

`read`权限配置`roles/deploy/files/read-group-rbac.yaml`是在`clusterrole view`基础上增加了若干读权限（Nodes/Persistent Volume Claims）

### read 用户证书

准备 read 证书请求：`read-csr.json`

``` bash
{
  "CN": "read",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "HangZhou",
      "L": "XS",
      "O": "group:read",
      "OU": "System"
    }
  ]
}
```
- 注意: O `group:read`，kube-apiserver 收到该证书后将请求的 Group 设置为`group:read`；之前步骤创建的 ClusterRoleBinding `read-clusterrole-binding`将 Group `group:read`与 ClusterRole `read-clusterrole`绑定，从而实现只读权限。

### read kubeconfig

kubeconfig 为与apiserver交互使用的认证配置文件，如脚本步骤需要：

- 设置集群参数，指定CA证书和apiserver地址
- 设置客户端认证参数，指定使用read证书和私钥
- 设置上下文参数，指定使用cluster集群和用户read
- 设置指定默认上下文

创建完成后生成默认配置文件为 `~/.kube/config`

## 恢复 admin 权限

- 可以恢复之前备份的`~/.kubeadmin`文件：`mv ~/.kube ~/.kuberead && mv ~/.kubeadmin ~/.kube`

## 参考

- [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [A Read Only Kubernetes Dashboard](https://blog.cowger.us/2018/07/03/a-read-only-kubernetes-dashboard.html)
