# IPVS 服务负载均衡

kube-proxy 组件监听 API server 中 service 和 endpoint 的变化情况，从而为 k8s 集群内部的 service 提供动态负载均衡。在v1.10之前主要通过 iptables来实现，是稳定、推荐的方式，但是当服务多的时候会产生太多的 iptables 规则，大规模情况下有明显的性能问题；在v1.11 GA的 ipvs高性能负载模式，采用增量式更新，并可以保证 service 更新期间连接的保持。

- NOTE: k8s v1.11.0 CentOS7下使用ipvs模式会有问题（见 kubernetes/kubernetes#65461），测试 k8s v1.10.2 CentOS7 可以。

## 启用 ipvs

建议 k8s 版本1.13 及以后启用 ipvs，只要在 kube-proxy 启动参数（或者配置文件中）中增加 `--proxy-mode=ipvs`:

``` bash
[Unit]
Description=Kubernetes Kube-Proxy Server
After=network.target

[Service]
WorkingDirectory=/var/lib/kube-proxy
ExecStart={{ bin_dir }}/kube-proxy \
  --bind-address={{ NODE_IP }} \
  --hostname-override={{ NODE_IP }} \
  --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig \
  --logtostderr=true \
  --proxy-mode=ipvs
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
