# calico 配置 BGP Route Reflectors

`Calico`作为`k8s`的一个流行网络插件，它依赖`BGP`路由协议实现集群节点上的`POD`路由互通；而路由互通的前提是节点间建立 BGP Peer 连接。BGP 路由反射器（Route Reflectors，简称 RR）可以简化集群BGP Peer的连接方式，它是解决BGP扩展性问题的有效方式；具体来说：

- 没有 RR 时，所有节点之间需要两两建立连接（IBGP全互联），节点数量增加将导致连接数剧增、资源占用剧增
- 引入 RR 后，其他 BGP 路由器只需要与它建立连接并交换路由信息，节点数量增加连接数只是线性增加，节省系统资源

calico-node 版本 v3.3 开始支持内建路由反射器，非常方便，因此使用 calico 作为网络插件可以支持大规模节点数的`K8S`集群。

- 建议集群节点数大于50时，应用BGP Route Reflectors 特性

## 前提条件

k8s 集群使用calico网络插件部署成功。本实验环境为按照kubeasz安装的2主2从集群，calico 版本 v3.19.4。

```
$ kubectl get node
NAME           STATUS                     ROLES    AGE    VERSION
192.168.1.1   Ready,SchedulingDisabled   master   178m   v1.13.1
192.168.1.2   Ready,SchedulingDisabled   master   178m   v1.13.1
192.168.1.3   Ready                      node     178m   v1.13.1
192.168.1.4   Ready                      node     178m   v1.13.1
$ kubectl get pod -n kube-system -o wide | grep calico
calico-kube-controllers-77487546bd-jqrlc   1/1     Running   0          179m   192.168.1.3   192.168.1.3   <none>           <none>
calico-node-67t5m                          2/2     Running   0          179m   192.168.1.1   192.168.1.1   <none>           <none>
calico-node-drmhq                          2/2     Running   0          179m   192.168.1.2   192.168.1.2   <none>           <none>
calico-node-rjtkv                          2/2     Running   0          179m   192.168.1.4   192.168.1.4   <none>           <none>
calico-node-xtspl                          2/2     Running   0          179m   192.168.1.3   192.168.1.3   <none>           <none>
```
查看当前集群中BGP连接情况：可以看到集群中4个节点两两建立了 BGP 连接

```
$ dk ansible -i /etc/kubeasz/clusters/xxx/hosts all -m shell -a '/opt/kube/bin/calicoctl node status'
192.168.1.3 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+--------------+-------------------+-------+----------+-------------+
| 192.168.1.1 | node-to-node mesh | up    | 03:08:20 | Established |
| 192.168.1.2 | node-to-node mesh | up    | 03:08:18 | Established |
| 192.168.1.4 | node-to-node mesh | up    | 03:08:19 | Established |
+--------------+-------------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.2 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+--------------+-------------------+-------+----------+-------------+
| 192.168.1.4 | node-to-node mesh | up    | 03:08:17 | Established |
| 192.168.1.3 | node-to-node mesh | up    | 03:08:18 | Established |
| 192.168.1.1 | node-to-node mesh | up    | 03:08:20 | Established |
+--------------+-------------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.1 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+--------------+-------------------+-------+----------+-------------+
| 192.168.1.2 | node-to-node mesh | up    | 03:08:21 | Established |
| 192.168.1.3 | node-to-node mesh | up    | 03:08:21 | Established |
| 192.168.1.4 | node-to-node mesh | up    | 03:08:21 | Established |
+--------------+-------------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.4 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+--------------+-------------------+-------+----------+-------------+
| 192.168.1.2 | node-to-node mesh | up    | 03:08:17 | Established |
| 192.168.1.3 | node-to-node mesh | up    | 03:08:19 | Established |
| 192.168.1.1 | node-to-node mesh | up    | 03:08:20 | Established |
+--------------+-------------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```

## kubeasz 自动安装启用 route reflector

- 修改`/etc/kubeasz/clusters/xxx/config.yml`文件，设置配置项`CALICO_RR_ENABLED: true` 
- 重新执行网络安装 `dk ezctl setup xxx 07`

执行完成，检查bgp连接验证即可。

### 附：手动安装route reflector 过程讲解

- 选择并配置 Route Reflector 节点

首先查看当前集群中的节点：

```
$ calicoctl get node -o wide
NAME     ASN       IPV4              IPV6
k8s401   (64512)   192.168.1.1/24
k8s402   (64512)   192.168.1.2/24
k8s403   (64512)   192.168.1.3/24
k8s404   (64512)   192.168.1.4/24
```

可以在集群中选择1个或多个节点作为 rr 节点，这里先选择节点：k8s401

``` bash
#配置routeReflectorClusterID
calicoctl patch node k8s401 -p '{"spec": {"bgp": {"routeReflectorClusterID": "244.0.0.1"}}}'

#配置node label
calicoctl patch node k8s401 -p '{"metadata": {"labels": {"route-reflector": "true"}}}'
```

- 配置 BGP node 与 Route Reflector 的连接建立规则

``` bash
$ cat << EOF | calicoctl create -f -
kind: BGPPeer
apiVersion: projectcalico.org/v3
metadata:
  name: peer-with-route-reflectors
spec:
  nodeSelector: all()
  peerSelector: route-reflector == 'true'
EOF
```

- 配置全局禁用全连接（BGP full mesh）

```
$ cat << EOF | calicoctl create -f -
apiVersion: projectcalico.org/v3
kind: BGPConfiguration
metadata:
  name: default
spec:
  logSeverityScreen: Info
  nodeToNodeMeshEnabled: false
  asNumber: 64512
EOF
```

- 验证增加 rr 之后的bgp 连接情况

``` 
$ dk ansible -i /etc/kubeasz/clusters/xxx/hosts all -m shell -a '/opt/kube/bin/calicoctl node status'
192.168.1.4 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-----------+-------+----------+-------------+
| PEER ADDRESS |   PEER TYPE   | STATE |  SINCE   |    INFO     |
+--------------+-----------+-------+----------+-------------+
| 192.168.1.1 | node specific | up    | 11:02:55 | Established |
+--------------+-----------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.3 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-----------+-------+----------+-------------+
| PEER ADDRESS | PEER TYPE | STATE |  SINCE   |    INFO     |
+--------------+-----------+-------+----------+-------------+
| 192.168.1.1 | node specific | up    | 11:02:55 | Established |
+--------------+-----------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.1 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+---------------+-------+----------+-------------+
| PEER ADDRESS |   PEER TYPE   | STATE |  SINCE   |    INFO     |
+--------------+---------------+-------+----------+-------------+
| 192.168.1.2 | node specific | up    | 11:02:55 | Established |
| 192.168.1.3 | node specific | up    | 11:02:55 | Established |
| 192.168.1.4 | node specific | up    | 11:02:55 | Established |
+--------------+---------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.

192.168.1.2 | SUCCESS | rc=0 >>
Calico process is running.

IPv4 BGP status
+--------------+-----------+-------+----------+-------------+
| PEER ADDRESS | PEER TYPE | STATE |  SINCE   |    INFO     |
+--------------+-----------+-------+----------+-------------+
| 192.168.1.1 | node specific | up    | 11:02:55 | Established |
+--------------+-----------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```
可以看到所有其他节点都与所选rr节点建立bgp连接。

- 再增加一个 rr 节点(略)

步骤同上，添加成功后可以看到所有其他节点都与两个rr节点建立bgp连接，两个rr节点之间也建立bgp连接。对于节点数较多的`K8S`集群建议配置2-3个 RR 节点。

## 参考文档

- 1.[Calico bgp 配置指南](https://projectcalico.docs.tigera.io/reference/resources/bgpconfig)
- 2.[BGP路由反射器基础](https://www.sohu.com/a/140033025_761420)
