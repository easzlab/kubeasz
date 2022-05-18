# 配置负载转发 ingress nodeport

向集群外暴露 ingress-controller 本身的服务端口（80/443/8080）一般有以下三种方法：

- 1.部署ingress-controller时使用`hostNetwork: true`，这样就可以直接使用上述端口，可能与host已listen端口冲突
- 2.部署ingress-controller时使用`LoadBalancer`类型服务，需要集群支持`LoadBalancer`
- 3.部署ingress-controller时使用`nodePort`类型服务，然后在集群外使用 haproxy/f5 等配置 virtual server 集群

本文档讲解使用 haproxy 配置 ingress的 VS 集群，前提是配置了自建`ex_lb`节点

## 1.配置 ex_lb 参数开启转发 ingress nodeport

``` bash
# 编辑 roles/ex-lb/defaults/main.yml，配置如下变量
INGRESS_NODEPORT_LB: "yes"
INGRESS_TLS_NODEPORT_LB: "yes"
```

## 2.重新配置启动LB节点服务

``` bash
$ ezctl setup ${集群名} ex-lb 
```

## 3.验证 ex_lb 节点的 haproxy 服务配置 `/etc/haproxy/haproxy.cfg` 包含如下配置

``` bash
... 前文省略
listen kube_master
        bind 0.0.0.0:8443
        mode tcp
        option tcplog
        balance roundrobin
        server 192.168.1.1 192.168.1.1:6443 check inter 2000 fall 2 rise 2 weight 1
        server 192.168.1.2 192.168.1.2:6443 check inter 2000 fall 2 rise 2 weight 1

listen ingress-node
        bind 0.0.0.0:80
        mode tcp
        option tcplog
        balance roundrobin
        server 192.168.1.3 192.168.1.3:23456 check inter 2000 fall 2 rise 2 weight 1
        server 192.168.1.4 192.168.1.4:23456 check inter 2000 fall 2 rise 2 weight 1

listen ingress-node-tls
        bind 0.0.0.0:443
        mode tcp
        option tcplog
        balance roundrobin
        server 192.168.1.3 192.168.1.3:23457 check inter 2000 fall 2 rise 2 weight 1
        server 192.168.1.4 192.168.1.4:23457 check inter 2000 fall 2 rise 2 weight 1
```

验证成功后，我们可以方便的去做[配置ingress](../guide/ingress.md)和[配置https ingress](../guide/ingress-tls.md)实验了。
