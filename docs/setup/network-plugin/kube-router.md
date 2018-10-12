# kube-router 网络组件

kube-router是一个简单、高效的网络插件，它提供一揽子解决方案：  
- 基于GoBGP 提供Pod 网络互联（Routing）
- 使用ipsets优化的iptables 提供网络策略支持（Firewall/NetworkPolicy）
- 基于IPVS/LVS 提供高性能服务代理（Service Proxy）

更多介绍请前往`https://github.com/cloudnativelabs/kube-router`

## 配置

本项目提供多种网络插件可选，如果需要安装kube-router，请在/etc/ansible/hosts文件中设置变量 `CLUSTER_NETWORK="kube-router"`，更多设置请查看`roles/kube-router/defaults/main.yml`

- kube-router需要在所有master节点和node节点安装

## 安装

- 单步安装已经集成：`ansible-playbook 90.setup.yml`
- 分步安装请执行：`ansible-playbook 06.network.yml`

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

- 5.ipset列表查看

``` bash
$ ipset list
...
Name: kube-router-pod-subnets
Type: hash:net
Revision: 6
Header: family inet hashsize 1024 maxelem 65536 timeout 0
Size in memory: 672
References: 2
Members:
172.20.1.0/24 timeout 0
172.20.2.0/24 timeout 0
172.20.0.0/24 timeout 0

Name: kube-router-node-ips
Type: hash:ip
Revision: 4
Header: family inet hashsize 1024 maxelem 65536 timeout 0
Size in memory: 416
References: 1
Members:
192.168.1.1 timeout 0
192.168.1.2 timeout 0
192.168.1.3 timeout 0
...
```

- 6.ipvs虚拟服务器查看 (roles/kube-router/defaults/main.yml 需配置`SERVICE_PROXY: "true"`)

``` bash
# 首先创建测试应用
$ kubectl run nginx --image=nginx --replicas=3 --port=80 --expose

# 查看ipvsadm输出
$ ipvsadm
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.68.0.1:https rr persistent 10800	# 这个kubernetes虚拟服务地址
  -> 192.168.1.1:6443            Masq    1      0          0         
TCP  10.68.199.39:http rr			# 这个是测试应用nginx的虚拟服务地址
  -> 172.20.1.5:http              Masq    1      0          0         
  -> 172.20.2.6:http              Masq    1      0          0         
  -> 172.20.2.8:http              Masq    1      0          0         
```
