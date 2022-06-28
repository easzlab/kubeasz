## 06-安装kube-ovn网络组件.md

(以下文档暂未更新，以插件官网文档为准)

由灵雀云开源的网络组件 kube-ovn，将已被 openstack 社区采用的成熟网络虚拟化技术 ovs/ovn 引入 kubernetes 平台；为 kubernetes 网络打开了新的大门，令人耳目一新；强烈推荐大家试用该网络组件，反馈建议以帮助项目早日走向成熟。

- 介绍 https://blog.csdn.net/alauda_andy/article/details/88886128
- 项目地址 https://github.com/alauda/kube-ovn

### 特性介绍

kube-ovn 提供了针对企业应用场景下容器网络实用功能，并为实现更高级的网络管理控制提供了可能性；现有主要功能:

- 1.Namespace 和子网的绑定，以及子网间的访问控制;
- 2.静态IP分配;
- 3.动态QoS;
- 4.分布式和集中式网关;
- 5.内嵌 LoadBalancer;
- 6.Pod IP对外直接暴露
- 7.流量镜像
- 8.IPv6

### kubeasz 集成安装 kube-ovn

kube-ovn 的安装十分简单，详见项目的安装文档；基于 kubeasz，以下两步将安装一个集成了 kube-ovn 网络的 k8s 集群；

- 在 ansible hosts 中设置变量 `CLUSTER_NETWORK="kube-ovn"`
- 执行安装 `ansible-playbook 90.setup.yml` 或者 `ezctl setup`

kubeasz 项目为`kube-ovn`网络生成的 ansible role 如下：

``` bash
roles/kube-ovn
├── defaults
│   └── main.yml		# kube-ovn 相关配置文件
├── tasks
│   └── main.yml		# 安装执行文件
└── templates
    ├── crd.yaml.j2	        # crd 模板
    ├── kube-ovn.yaml.j2	# kube-ovn yaml 模板
    └── ovn.yaml.j2		    # ovn yaml 模板
    
```

安装成功后，可以验证所有 k8s 集群功能正常，查看集群的 pod 网络如下：

```
$ kubectl get pod --all-namespaces -o wide
NAMESPACE     NAME                                    READY   STATUS    RESTARTS   AGE   IP             NODE           NOMINATED NODE   READINESS GATES
kube-ovn      kube-ovn-cni-5php2                      1/1     Running   2          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      kube-ovn-cni-7dwmx                      1/1     Running   2          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-ovn      kube-ovn-cni-lhlvl                      1/1     Running   2          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      kube-ovn-controller-57955db7b4-6x6hd    1/1     Running   0          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      kube-ovn-controller-57955db7b4-chvz4    1/1     Running   0          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-ovn      ovn-central-bb8747d77-tr5nz             1/1     Running   0          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      ovs-ovn-2qhhr                           1/1     Running   0          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      ovs-ovn-np8rn                           1/1     Running   0          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      ovs-ovn-pkjw4                           1/1     Running   0          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-system   coredns-55f46dd959-76qb5                1/1     Running   0          35h   10.16.0.12     192.168.1.42   <none>           <none>
kube-system   coredns-55f46dd959-wn8kw                1/1     Running   0          35h   10.16.0.11     192.168.1.43   <none>           <none>
kube-system   heapster-fdb7596d6-xmmrx                1/1     Running   0          35h   10.16.0.15     192.168.1.42   <none>           <none>
kube-system   kubernetes-dashboard-68ddcc97fc-dwzbf   1/1     Running   0          35h   10.16.0.14     192.168.1.42   <none>           <none>
kube-system   metrics-server-6c898b5b8b-zvct2         1/1     Running   0          35h   10.16.0.13     192.168.1.43   <none>           <none>
```

直观上 kube-ovn 与传统 k8s 网络（flannel/calico等）比较最大的不同是 pod 子网的分配：

- 传统网络插件下，集群中 pod 一般是不同 node 节点分配不同的子网；然后通过 overlay 等技术打通不同 node 节点的 pod 子网；
- kube-ovn 中 pod 网络根据其所在的 namespace 而定； namespace 在创建时可以根据 annotation 来配置它的子网/网关等参数；默认使用 10.16.0.0/16 的子网；

### 测试 namespace 子网分配

新建一个 subnet 并绑定 namespace 测试分配一个新的 pod 子网

```
# 创建一个 namespace: test-ns
$ cat > test-ns.yaml << EOF
apiVersion: v1
kind: Namespace
metadata:
  annotations:
  name: test-ns
EOF
$ kubectl apply -f test-ns.yaml

# 创建一个 subnet: test-subnet 并绑定 namespace test-ns
$ cat > test-subnet.yaml << EOF
apiVersion: kubeovn.io/v1
kind: Subnet
metadata:
  name: test-subnet
spec:
  protocol: IPv4
  default: false
  namespaces:
  - test-ns
  cidrBlock: 10.17.0.0/24
  gateway: 10.17.0.1
  excludeIps:
  - 10.17.0.1..10.17.0.10
EOF
$ kubectl apply -f test-subnet.yaml

# 在 test-ns 中创建 nginx 部署
$ kubectl run -n test-ns nginx --image=nginx --replicas=2 --port=80 --expose

# 在 default 中创建 busy 客户端
$ kubectl run busy --image=busybox sleep 360000
```

创建成功后，查看 pod 地址的分配，可以看到确实 test-ns 中 pod 使用新的子网，而 default 中 pod 使用了默认子网，并验证 pod 之间的联通性（默认可通）

```
$ kubectl get pod --all-namespaces -o wide
NAMESPACE     NAME                                    READY   STATUS    RESTARTS   AGE   IP             NODE           NOMINATED NODE   READINESS GATES
default       busy-6c55ccddc5-qrm5j                   1/1     Running   0          31h   10.16.0.16     192.168.1.43   <none>           <none>
kube-ovn      kube-ovn-cni-5php2                      1/1     Running   2          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      kube-ovn-cni-7dwmx                      1/1     Running   2          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-ovn      kube-ovn-cni-lhlvl                      1/1     Running   2          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      kube-ovn-controller-57955db7b4-6x6hd    1/1     Running   0          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      kube-ovn-controller-57955db7b4-chvz4    1/1     Running   0          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-ovn      ovn-central-bb8747d77-tr5nz             1/1     Running   0          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      ovs-ovn-2qhhr                           1/1     Running   0          35h   192.168.1.41   192.168.1.41   <none>           <none>
kube-ovn      ovs-ovn-np8rn                           1/1     Running   0          35h   192.168.1.43   192.168.1.43   <none>           <none>
kube-ovn      ovs-ovn-pkjw4                           1/1     Running   0          35h   192.168.1.42   192.168.1.42   <none>           <none>
kube-system   coredns-55f46dd959-76qb5                1/1     Running   0          35h   10.16.0.12     192.168.1.42   <none>           <none>
kube-system   coredns-55f46dd959-wn8kw                1/1     Running   0          35h   10.16.0.11     192.168.1.43   <none>           <none>
kube-system   heapster-fdb7596d6-xmmrx                1/1     Running   0          35h   10.16.0.15     192.168.1.42   <none>           <none>
kube-system   kubernetes-dashboard-68ddcc97fc-dwzbf   1/1     Running   0          35h   10.16.0.14     192.168.1.42   <none>           <none>
kube-system   metrics-server-6c898b5b8b-zvct2         1/1     Running   0          35h   10.16.0.13     192.168.1.43   <none>           <none>
test-ns       nginx-755464dd6c-s6flj                  1/1     Running   0          31h   10.17.0.12     192.168.1.42   <none>           <none>
test-ns       nginx-755464dd6c-zct56                  1/1     Running   0          31h   10.17.0.11     192.168.1.43   <none>           <none>
```

- 更多的测试（pod网络QOS限速，namespace网络隔离等）请参考 kube-ovn 项目说明文档

### 延伸阅读

- [kube-ovn 官方文档](https://github.com/alauda/kube-ovn/tree/master/docs)
- [从 Bridge 到 OVS，探索虚拟交换机](https://www.cnblogs.com/bakari/p/8097439.html)
