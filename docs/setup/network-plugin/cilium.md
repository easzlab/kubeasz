# 06-安装cilium网络组件

`cilium` 是一个革新的网络与安全组件；基于 linux 内核新技术--`BPF`，它可以透明、零侵入地实现服务间安全策略与可视化，主要优势如下：

- 支持L3/L4, L7(如：HTTP/gRPC/Kafka)的安全策略
- 支持基于安全ID而不是地址+端口的传统防火墙策略
- 支持基于Overlay或Native Routing的扁平多节点pod网络
  - Overlay VXLAN 方式类似于 flannel 的VXLAN后端
- 高性能负载均衡，支持DSR
- 支持事件、策略跟踪和监控集成

cilium 项目当前文档比较完整，建议仔细阅读下[官网文档]()

## kubeasz 集成安装 cilium

kubeasz 3.3.1 更新重写了cilium 安装流程，使用helm charts 方式，配置文件在 roles/cilium/templates/values.yaml.j2，请阅读原charts中values.yaml 文件后自定义修改。

- 相关镜像已经离线打包并推送到本地镜像仓库，通过 `ezdown -X` 命令下载cilium等额外镜像

### 0.升级内核并重启

- Linux kernel >= 4.9.17，请阅读文档[升级内核](guide/kernel_upgrade.md)
- etcd >= 3.1.0 or consul >= 0.6.4

### 1.选择cilium网络后安装

- 参考[快速指南](quickStart.md)，设置`/etc/kubeasz/clusters/xxx/hosts`文件中变量 `CLUSTER_NETWORK="cilium"` 
- 执行集群安装 `dk ezctl setup xxx all`

注意默认安装后集成了cilium_connectivity_check 和 cilium_hubble，可以在`/etc/kubeasz/clusters/xxx/config.yml`配置关闭

- cilium_connectivity_check：检查集群cilium网络是否工作正常，非常实用
- cilium_hubble：很酷很实用的监控、策略追踪排查工具

Cilium CLI 和 Hubble CLI 二进制已经默认包含在kubeasz-ext-bin 1.2.0版本中 https://github.com/kubeasz/dockerfiles/blob/master/kubeasz-ext-bin/Dockerfile

### 2.验证

一键安装完成后如下，注意cilium_connectivity_check 中带`multi-node`的检查任务需要多节点集群才能完成

```
kubectl get pod -A
NAMESPACE     NAME                                                    READY   STATUS    RESTARTS   AGE
cilium-test   echo-a-5dd478f5d8-74xg5                                 1/1     Running   0          3m10s
cilium-test   echo-b-78c79f6cdd-t9vk6                                 1/1     Running   0          3m10s
cilium-test   echo-b-host-75c44b897-c8f5m                             1/1     Running   0          3m10s
cilium-test   host-to-b-multi-node-clusterip-7895fd494c-92cb2         1/1     Running   0          2m59s
cilium-test   host-to-b-multi-node-headless-74bbc877b5-ffxxx          1/1     Running   0          2m59s
cilium-test   pod-to-a-allowed-cnp-598fc5c547-b885q                   1/1     Running   0          2m59s
cilium-test   pod-to-a-b8b456c99-r6272                                1/1     Running   0          2m59s
cilium-test   pod-to-a-denied-cnp-c78c44f5c-7xhkw                     1/1     Running   0          2m59s
cilium-test   pod-to-b-intra-node-nodeport-6ccdb55779-j8gnd           1/1     Running   0          2m59s
cilium-test   pod-to-b-multi-node-clusterip-55d8448b5c-5b4nj          1/1     Running   0          2m59s
cilium-test   pod-to-b-multi-node-headless-5fbf655bb9-pszpr           1/1     Running   0          2m59s
cilium-test   pod-to-b-multi-node-nodeport-65f5b95569-qglb7           1/1     Running   0          2m59s
cilium-test   pod-to-external-1111-64496c754c-bvqlt                   1/1     Running   0          2m59s
cilium-test   pod-to-external-fqdn-allow-baidu-cnp-6f96597855-c84zs   1/1     Running   0          2m59s
kube-system   cilium-7trcs                                            1/1     Running   0          3m42s
kube-system   cilium-hvclp                                            1/1     Running   0          3m42s
kube-system   cilium-operator-8566689975-vcxpp                        1/1     Running   0          3m42s
kube-system   cilium-pw2sv                                            1/1     Running   0          3m42s
kube-system   cilium-qppnc                                            1/1     Running   0          3m42s
kube-system   coredns-84b58f6b4-m8x7s                                 1/1     Running   0          3m20s
kube-system   dashboard-metrics-scraper-864d79d497-92l2w              1/1     Running   0          3m14s
kube-system   hubble-relay-655dc744d7-8d9n7                           1/1     Running   0          3m42s
kube-system   hubble-ui-54599d7967-lfkvk                              2/2     Running   0          3m42s
kube-system   kubernetes-dashboard-5fc74cf5c6-pqdvc                   1/1     Running   0          3m14s
kube-system   metrics-server-69797698d4-2jbg8                         1/1     Running   0          3m17s
kube-system   node-local-dns-5n8gc                                    1/1     Running   0          3m19s
kube-system   node-local-dns-5pm2p                                    1/1     Running   0          3m19s
kube-system   node-local-dns-9x229                                    1/1     Running   0          3m19s
kube-system   node-local-dns-jz8lj                                    1/1     Running   0          3m19s
```

检查 cilium 节点状态

```
cilium status
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         OK
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

DaemonSet         cilium             Desired: 4, Ready: 4/4, Available: 4/4
Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-relay       Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-ui          Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium             Running: 4
                  cilium-operator    Running: 1
                  hubble-relay       Running: 1
                  hubble-ui          Running: 1
Cluster Pods:     17/17 managed by Cilium
Image versions    hubble-relay       easzlab.io.local:5000/cilium/hubble-relay:v1.11.6: 1
                  hubble-ui          easzlab.io.local:5000/cilium/hubble-ui:v0.9.0: 1
                  hubble-ui          easzlab.io.local:5000/cilium/hubble-ui-backend:v0.9.0: 1
                  cilium             easzlab.io.local:5000/cilium/cilium:v1.11.6: 4
                  cilium-operator    easzlab.io.local:5000/cilium/operator-generic:v1.11.6: 1
```

## cilium network policy

cilium network policy 提供了比k8s network policy更丰富的网络安全策略功能，有兴趣的请阅读官网文档，以下是一个有趣的小例子：

- [星战死星登陆系统](cilium-example.md)
