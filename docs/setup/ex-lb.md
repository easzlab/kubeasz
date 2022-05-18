## EX-LB 负载均衡部署

根据[HA 2x架构](00-planning_and_overall_intro.md)，k8s集群自身高可用已经不依赖于外部 lb 服务；但是有时我们要从外部访问 apiserver（比如 CI 流程），就需要 ex_lb 来请求多个 apiserver；

还有一种情况是需要[负载转发到ingress服务](../op/loadballance_ingress_nodeport.md)，也需要部署ex_lb；

**注意：当遇到公有云环境无法自建 ex_lb 服务时，可以配置对应的云负载均衡服务**

### ex_lb 服务组件

更新：kubeasz 3.0.2 重写了ex-lb服务安装，利用最小化依赖编译安装的二进制文件，不依赖于linux发行版；优点是可以统一版本和简化离线安装部署，并且理论上能够支持更多linux发行版


ex_lb 服务由 keepalived 和 l4lb 组成：
- l4lb：是一个精简版（仅支持四层转发）的nginx编译二进制版本
- keepalived：利用主备节点vrrp协议通信和虚拟地址，消除l4lb的单点故障；keepalived保持存活，它是基于VRRP协议保证所谓的高可用或热备的，这里用来预防l4lb的单点故障。

keepalived与l4lb配合，实现master的高可用过程如下：

+ 1.keepalived利用vrrp协议生成一个虚拟地址(VIP)，正常情况下VIP存活在keepalive的主节点，当主节点故障时，VIP能够漂移到keepalived的备节点，保障VIP地址高可用性。
+ 2.在keepalived的主备节点都配置相同l4lb负载配置，并且监听客户端请求在VIP的地址上，保障随时都有一个l4lb负载均衡在正常工作。并且keepalived启用对l4lb进程的存活检测，一旦主节点l4lb进程故障，VIP也能切换到备节点，从而让备节点的l4lb进行负载工作。
+ 3.在l4lb的配置中配置多个后端真实kube-apiserver的endpoints，并启用存活监测后端kube-apiserver，如果一个kube-apiserver故障，l4lb会将其剔除负载池。

#### 安装l4lb

#### 配置l4lb (roles/ex-lb/templates/l4lb.conf.j2)

配置由全局配置和三个upstream servers配置组成：
- apiservers 用于转发至多个apiserver
- ingress-nodes 用于转发至node节点的ingress http服务，[参阅](../op/loadballance_ingress_nodeport.md)
- ingress-tls-nodes 用于转发至node节点的ingress https服务

#### 安装keepalived

#### 配置keepalived主节点 [keepalived-master.conf.j2](../../roles/ex-lb/templates/keepalived-master.conf.j2)

``` bash
global_defs {
}

vrrp_track_process check-l4lb {
    process l4lb
    weight -60
    delay 3
}

vrrp_instance VI-01 {
    state MASTER
    priority 120
    unicast_src_ip {{ inventory_hostname }}
    unicast_peer {
{% for h in groups['ex_lb'] %}{% if h != inventory_hostname %}
        {{ h }}
{% endif %}{% endfor %}
    }
    dont_track_primary
    interface {{ LB_IF }}
    virtual_router_id {{ ROUTER_ID }}
    advert_int 3
    track_process {
        check-l4lb
    }
    virtual_ipaddress {
        {{ EX_APISERVER_VIP }}
    }
}
```
+ vrrp_track_process 定义了监测l4lb进程是否存活，如果进程不存在，根据`weight -60`设置将主节点优先级降低60，这样原先备节点将变成主节点。
+ vrrp_instance 定义了vrrp组，包括优先级、使用端口、router_id、心跳频率、检测脚本、虚拟地址VIP等
+ 特别注意 `virtual_router_id` 标识了一个 VRRP组，在同网段下必须唯一，否则出现 `Keepalived_vrrp: bogus VRRP packet received on eth0 !!!`类似报错
+ 配置 vrrp 协议通过单播发送

#### 配置keepalived备节点 [keepalived-backup.conf.j2](../../roles/ex-lb/templates/keepalived-backup.conf.j2)

+ 备节点的配置类似主节点，除了优先级和检测脚本，其他如 `virtual_router_id` `advert_int` `virtual_ipaddress`必须与主节点一致

### 启动 keepalived 和 l4lb 后验证

+ lb 节点验证

``` bash
systemctl status l4lb 	# 检查进程状态
journalctl -u l4lb		# 检查进程日志是否有报错信息
systemctl status keepalived 	# 检查进程状态
journalctl -u keepalived	# 检查进程日志是否有报错信息
```
+ 在 keepalived 主节点

``` bash
ip a				# 检查 master的 VIP地址是否存在
```
### keepalived 主备切换演练

1. 尝试关闭 keepalived主节点上的 l4lb进程，然后在keepalived 备节点上查看 master的 VIP地址是否能够漂移过来，并依次检查上一步中的验证项。
1. 尝试直接关闭 keepalived 主节点系统，检查各验证项。

