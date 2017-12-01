## 03-配置kubectl命令行工具.md

kubectl使用~/.kube/config 配置文件与kube-apiserver进行交互，且拥有完全权限[可配置]，因此尽量避免安装在不必要的节点上，这里为了演示方便，将它安装在master/node/deploy节点。
`cat ~/.kube/config`可以看到配置文件包含 kube-apiserver 地址、证书、用户名等信息。

``` bash
roles/kubectl
├── tasks
│   └── main.yml
└── templates
    └── admin-csr.json.j2
```
请在另外窗口打开[roles/kubectl/tasks/main.yml](../roles/kubectl/tasks/main.yml) 文件，对照看以下讲解内容。

### 准备kubectl使用的admin 证书签名请求 [admin-csr.json.j2](../roles/kubectl/templates/admin-csr.json.j2)

``` bash
{
  "CN": "admin",
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
      "O": "system:masters",
      "OU": "System"
    }
  ]
}

```
+ 后续我们在安装`master`节点时候会启用 `RBAC`特性，它在v1.8.x中已是稳定版本，推荐[RBAC官方文档](https://kubernetes.io/docs/admin/authorization/rbac/)
+ 证书请求中 `O` 指定该证书的 Group 为 `system:masters`，而 `RBAC` 预定义的 `ClusterRoleBinding` 将 Group `system:masters` 与 ClusterRole `cluster-admin` 绑定，这就赋予了kubectl**所有集群权限**

kubectl get clusterrolebinding cluster-admin -o yaml

``` bash
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  creationTimestamp: 2017-11-30T01:33:10Z
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: cluster-admin
  resourceVersion: "76"
  selfLink: /apis/rbac.authorization.k8s.io/v1/clusterrolebindings/cluster-admin
  uid: 6c9dd451-d56e-11e7-8ed6-525400103a5d
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:masters
```
### 创建admin 证书和私钥

### 创建 kubectl kubeconfig 文件

#### 设置集群参数，指定CA证书和apiserver地址

``` bash
{{ bin_dir }}/kubectl config set-cluster kubernetes \
        --certificate-authority={{ ca_dir }}/ca.pem \
        --embed-certs=true \
        --server={{ KUBE_APISERVER }}
```

#### 设置客户端认证参数，指定使用admin证书和私钥

``` bash
{{ bin_dir }}/kubectl config set-credentials admin \
        --client-certificate={{ ca_dir }}/admin.pem \
        --embed-certs=true \
        --client-key={{ ca_dir }}/admin-key.pem
```

#### 设置上下文参数，说明使用cluster集群和用户admin

``` bash
{{ bin_dir }}/kubectl config set-context kubernetes \
        --cluster=kubernetes --user=admin
```

#### 选择默认上下文

``` bash
{{ bin_dir }}/kubectl config use-context kubernetes
```
+ 注意{{ }}中参数与ansible hosts文件中设置对应
+ 以上生成的 kubeconfig 自动保存到 ~/.kube/config 文件
