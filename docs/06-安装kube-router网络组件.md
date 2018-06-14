# kube-router 网络组件




## 验证

- 1.pod间网络联通性：略

- 2.host路由表

``` bash
# master上路由
root@master1:~$ ip route
...
172.20.1.0/24 via 192.168.1.2 dev ens3  proto 17 
172.20.2.0/24 via 192.168.1.3 dev ens3  proto 17 
...

# node3上路由
root@node3:~$ ip route
... 
172.20.0.0/24 via 192.168.1.1 dev ens3  proto 17 
172.20.1.0/24 via 192.168.1.2 dev ens3  proto 17 
172.20.2.0/24 dev kube-bridge  proto kernel  scope link  src 172.20.2.1 
...
```

- 3.bgp连接状态

``` bash
# master上
root@master1:~$ netstat -antlp|grep router|grep LISH|grep 179
tcp        0      0 192.168.1.1:179        192.168.1.3:58366      ESTABLISHED 26062/kube-router
tcp        0      0 192.168.1.1:42537      192.168.1.2:179        ESTABLISHED 26062/kube-router

# node3上
root@node3:~$ netstat -antlp|grep router|grep LISH|grep 179
tcp        0      0 192.168.1.3:58366      192.168.1.1:179        ESTABLISHED 18897/kube-router
tcp        0      0 192.168.1.3:179        192.168.1.2:43928      ESTABLISHED 18897/kube-router

```

- 4.NetworkPolicy有效性，验证参照[这里](guide/networkpolicy.md)

