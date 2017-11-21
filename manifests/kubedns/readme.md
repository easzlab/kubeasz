## kubedns
### 部署文件
1. [kubedns-cm](./kubedns-cm.yaml)
1. [kubedns-controller](./kubedns-controller.yaml)
1. [kubedns-sa](./kubedns-sa.yaml)
1. [kubedns-svc](./kubedns-svc.yaml)

### pod继承node的dns解析
When running a pod, kubelet will prepend the cluster DNS server and search paths to the node’s own DNS settings. If the node is able to resolve DNS names specific to the larger environment, pods should be able to, also. See “Known issues” below for a caveat.

If you don’t want this, or if you want a different DNS config for pods, you can use the kubelet’s --resolv-conf flag. Setting it to “” means that pods will not inherit DNS. Setting it to a valid file path means that kubelet will use this file instead of /etc/resolv.conf for DNS inheritance.

### configmap配置私有dns服务器和上游dns服务器(未实验)
``` bash
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-dns
  namespace: kube-system
data:
  stubDomains: |
    {“acme.local”: [“1.2.3.4”]}
  upstreamNameservers: |
    [“8.8.8.8”, “8.8.4.4”]
```
使用上述特定配置，查询请求首先会被发送到kube-dns的DNS缓存层(Dnsmasq 服务器)。Dnsmasq服务器会先检查请求的后缀，带有集群后缀（例如：”.cluster.local”）的请求会被发往kube-dns，拥有存根域后缀的名称（例如：”.acme.local”）将会被发送到配置的私有DNS服务器[“1.2.3.4”]。最后，不满足任何这些后缀的请求将会被发送到上游DNS [“8.8.8.8”, “8.8.4.4”]里。
