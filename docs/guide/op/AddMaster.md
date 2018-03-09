## 增加 kube-master 节点

注意：目前仅支持按照本项目`多主模式`(hosts.m-masters.example)部署的`k8s`集群增加`master`节点

新增`kube-master`节点大致流程为：
- 节点预处理 prepare
- 重新配置LB节点的 haproxy服务
- 安装 master 节点服务

### 操作步骤

按照本项目说明，首先确保deploy节点能够ssh免密码登陆新增节点，然后在**deploy**节点执行三步：

- 修改ansible hosts 文件，在 [kube-master] 组添加新增的节点；在[lb] 组添加新增master 节点，举例如下：

``` bash
[kube-master]
192.168.1.1 NODE_IP="192.168.1.1"
192.168.1.2 NODE_IP="192.168.1.2"
192.168.1.5 NODE_IP="192.168.1.5"   # 新增 master节点

[lb]
192.168.1.1 LB_IF="ens3" LB_ROLE=backup
192.168.1.4 LB_IF="ens3" LB_ROLE=master
[lb:vars]
master1="192.168.1.1:6443"
master2="192.168.1.2:6443"
master3="192.168.1.5:6443"   # 新增 master节点
```
- 修改roles/lb/templates/haproxy.cfg.j2 文件，增加新增的master节点，举例如下：

``` bash
listen kube-master
        bind 0.0.0.0:{{ MASTER_PORT }}
        mode tcp
        option tcplog
        balance source
        server s1 {{ master1 }}  check inter 10000 fall 2 rise 2 weight 1
        server s2 {{ master2 }}  check inter 10000 fall 2 rise 2 weight 1
        server s3 {{ master3 }}  check inter 10000 fall 2 rise 2 weight 1
```

- 执行安装脚本

``` bash
$ cd /etc/ansible && ansible-playbook 20.addmaster.yml
```

### 验证

``` bash
# 在新节点master 服务状态
$ systemctl status kube-apiserver 
$ systemctl status kube-controller-manager
$ systemctl status kube-scheduler

# 查看新master的服务日志
$ journalctl -u kube-apiserver -f

```
