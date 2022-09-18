## 开始使用 cilium

以下为简要翻译 `cilium doc`上的一个应用示例[原文](https://docs.cilium.io/en/stable/gettingstarted/http/)，部署在单节点k8s 环境的实践。

### 部署示例应用

官方文档用几个`pod/svc` 抽象一个有趣的应用场景（星战迷）：星战中帝国方建造了被称为“终极武器”的“死星”，它是一个卫星大小的战斗空间站，它的核心是使用凯伯晶体（Kyber Crystal）的超级激光炮，剧中它的首秀就以完全火力摧毁了“杰达圣城”（Jedha）。下面将用运行于 k8s上的 pod/svc/cilium 等模拟“死星“的一个“飞船登陆”系统安全策略设计。

- deploy/deathstar：作为控制整个“死星”的飞船登陆管理系统，它暴露一个SVC，提供HTTP REST 接口给飞船请求登陆使用；
- pod/tiefighter：作为“帝国”方的常规战斗飞船，它会调用上述 HTTP 接口，请求登陆“死星”；
- pod/xwing：作为“盟军”方的飞行舰，它也尝试调用 HTTP 接口，请求登陆“死星”；

<img alt="cilium_http_gsg" width="400" height="300" src="https://docs.cilium.io/en/stable/_images/cilium_http_gsg.png">

根据文件[http-sw-app.yaml](../../../roles/cilium/files/star_war_example/http-sw-app.yaml) 创建 `$ kubectl create -f http-sw-app.yaml` 后，验证如下：

``` bash
$ kubectl get pods,svc
NAME                             READY     STATUS    RESTARTS   AGE
pod/deathstar-5fc7c7795d-djf2q   1/1       Running   0          4h
pod/deathstar-5fc7c7795d-hrgst   1/1       Running   0          4h
pod/tiefighter                   1/1       Running   0          4h
pod/xwing                        1/1       Running   0          4h

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/deathstar    ClusterIP   10.68.242.130   <none>        80/TCP    4h
service/kubernetes   ClusterIP   10.68.0.1       <none>        443/TCP   5h
```
每个 POD 在 `cilium` 中都表示为 `Endpoint`，初始每个 `Endpoint` 的”进出安全策略“状态均为 `Disabled`，如下：(已省略部分无关 POD 信息)

``` bash
$ kubectl exec -n kube-system cilium-6t5vx -- cilium endpoint list
ENDPOINT   POLICY (ingress)   POLICY (egress)   IDENTITY   LABELS (source:key[=value])                                    IPv6                  IPv4           STATUS   
           ENFORCEMENT        ENFORCEMENT                                                                                                                      
643        Disabled           Disabled          31371      k8s:class=deathstar                                            f00d::ac14:0:0:283    172.20.0.246   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
1011       Disabled           Disabled          31371      k8s:class=deathstar                                            f00d::ac14:0:0:3f3    172.20.0.63    ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
32030      Disabled           Disabled          5350       k8s:class=tiefighter                                           f00d::ac14:0:0:7d1e   172.20.0.201   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
45943      Disabled           Disabled          14309      k8s:class=xwing                                                f00d::ac14:0:0:b377   172.20.0.189   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=alliance                                                                                            
52035      Disabled           Disabled          4          reserved:health                                                f00d::ac14:0:0:cb43   172.20.0.92    ready   
```

### 检查初始状态

当然“死星”应该只允许“帝国”的飞船着陆，因为没有应用任何策略，所以初始状态下“帝国”和“联盟”的飞船都可以登陆，如下测试：

``` bash
$ kubectl exec xwing -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
Ship landed # 成功着陆
$ kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
Ship landed # 成功着陆
```

### 应用 L3/L4 策略

现在我们应用策略，仅让带有标签 `org=empire`的飞船登陆“死星”；那么带有标签 `org=alliance`的“联盟”飞船将禁止登陆；这个就是我们熟悉的传统L3/L4 防火墙策略，并跟踪连接（会话）状态；

<img alt="cilium_http_l3_l4_gsg" width="400" height="300" src="https://docs.cilium.io/en/stable/_images/cilium_http_l3_l4_gsg.png">

根据文件[sw_l3_l4_policy.yaml](../../../roles/cilium/files/star_war_example/sw_l3_l4_policy.yaml) 创建 `$ kubectl apply -f sw_l3_l4_policy.yaml` 后，验证如下：

``` bash
$ kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
Ship landed # 成功着陆

$ kubectl exec xwing -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
# 失败超时
```

### 查看安全策略

再次执行 `cilium endpoint list`，可以看到标签带`deathstar`的 POD 已经应用了 `Ingress`方向的策略：

``` bash
# kubectl exec -n kube-system cilium-6t5vx -- cilium endpoint list
ENDPOINT   POLICY (ingress)   POLICY (egress)   IDENTITY   LABELS (source:key[=value])                                    IPv6                  IPv4           STATUS   
           ENFORCEMENT        ENFORCEMENT                                                                                                                      
643        Enabled            Disabled          31371      k8s:class=deathstar                                            f00d::ac14:0:0:283    172.20.0.246   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
1011       Enabled            Disabled          31371      k8s:class=deathstar                                            f00d::ac14:0:0:3f3    172.20.0.63    ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
32030      Disabled           Disabled          5350       k8s:class=tiefighter                                           f00d::ac14:0:0:7d1e   172.20.0.201   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=empire                                                                                              
45943      Disabled           Disabled          14309      k8s:class=xwing                                                f00d::ac14:0:0:b377   172.20.0.189   ready   
                                                           k8s:io.cilium.k8s.policy.serviceaccount=default                                                             
                                                           k8s:io.kubernetes.pod.namespace=default                                                                     
                                                           k8s:org=alliance                                                                                            
52035      Disabled           Disabled          4          reserved:health                                                f00d::ac14:0:0:cb43   172.20.0.92    ready   
```

查看具体策略内容 `kubectl describe cnp rule1`

### L7 安全策略

上述的策略可以进行简单的安全防护了，但是“死星”的这个系统还有很多复杂的功能；比如它还提供了一个内部维护接口，如果被不合理调用将带来严重灾难性后果，也许“联盟”勇士劫持了一架“帝国”飞船正在进行这个任务（虽然我们内心希望他能够成功摧毁“死星”）。不幸的是“死星”系统设计者考虑到这个风险，它有办法严格限制每架飞船能够请求的权限。

没有限制飞船请求权限时，如下运行：

``` bash
$ kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
Panic: deathstar exploded

goroutine 1 [running]:
main.HandleGarbage(0x2080c3f50, 0x2, 0x4, 0x425c0, 0x5, 0xa)
        /code/src/github.com/empire/deathstar/
        temp/main.go:9 +0x64
main.main()
        /code/src/github.com/empire/deathstar/
        temp/main.go:5 +0x85
```

<img alt="cilium_http_l3_l4_l7_gsg" width="400" height="300" src="https://docs.cilium.io/en/stable/_images/cilium_http_l3_l4_l7_gsg.png">

限制L7 的安全策略，根据文件[sw_l3_l4_l7_policy.yaml](../../../roles/cilium/files/star_war_example/sw_l3_l4_l7_policy.yaml) 创建 `$ kubectl apply -f sw_l3_l4_l7_policy.yaml` 后，验证如下：

``` bash
$ kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
Ship landed
$ kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
Access denied
```

我们同样可以使用 `kubectl desribe cnp`检查更新的策略，或者使用 `cilium` 命令行：

``` bash
$ kubectl exec -n kube-system cilium-6t5vx -- cilium policy get
[
  {
    "endpointSelector": {
      "matchLabels": {
        "any:class": "deathstar",
        "any:org": "empire",
        "k8s:io.kubernetes.pod.namespace": "default"
      }
    },
    "ingress": [
      {
        "fromEndpoints": [
          {
            "matchLabels": {
              "any:org": "empire",
              "k8s:io.kubernetes.pod.namespace": "default"
            }
          }
        ],
        "toPorts": [
          {
            "ports": [
              {
                "port": "80",
                "protocol": "TCP"
              }
            ],
            "rules": {
              "http": [
                {
                  "path": "/v1/request-landing",
                  "method": "POST"
                }
              ]
            }
          }
        ]
      }
    ],
    "labels": [
      {
        "key": "io.cilium.k8s.policy.name",
        "value": "rule1",
        "source": "k8s"
      },
      {
        "key": "io.cilium.k8s.policy.namespace",
        "value": "default",
        "source": "k8s"
      }
    ]
  }
]
Revision: 267
```
我们看到 `cilium` 可以实现 `7层 HTTP `协议的请求方法（GET/PUT/POST等）、路径（/v1/request-landing）等等安全策略；另外，它还可以防护其他应用（如：Kafka, gRPC, Elasticsearch），可以去官网文档示例学习！

## 参考资料

- [cilium github](https://github.com/cilium/cilium)
- [cilium doc](http://docs.cilium.io)
