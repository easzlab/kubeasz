# metallb 网络负载均衡

`Metallb`是在自有硬件上（非公有云）实现 `Kubernetes Load-balancer`的工具，由`google`团队开源，值得推荐！项目[github主页](https://github.com/google/metallb)。

## metallb 简介

这里简单介绍下它的实现原理，具体可以参考[metallb官网](https://metallb.universe.tf/)，文档非常简洁、清晰。目前有如下的使用限制：

- `Kubernetes v1.9.0`版本以上，暂不支持`ipvs`模式
- 支持网络组件 (flannel/weave/romana), calico 部分支持
- `layer2`和`bgp`两种模式，其中`bgp`模式需要外部网络设备支持`bgp`协议

`metallb`主要实现了两个功能：地址分配和对外宣告

- 地址分配：需要向网络管理员申请一段ip地址，如果是layer2模式需要这段地址与node节点地址同个网段（同一个二层）；如果是bgp模式没有这个限制。
- 对外宣告：layer2模式使用arp协议，利用节点的mac额外宣告一个loadbalancer的ip（同mac多ip）；bgp模式下节点利用bgp协议与外部网络设备建立邻居，宣告loadbalancer的地址段给外部网络。

## kubeasz 集成安装metallb

因bgp模式需要外部路由器的支持，这里主要选用layer2模式（如需选择bgp模式，相应修改roles/cluster-addon/templates/metallb/bgp.yaml.j2）。

- 1.修改roles/cluster-addon/defaults/main.yml 配置文件相关

``` bash
# metallb 自动安装
metallb_install: "yes"
# 模式选择: 二层 "layer2" 或者三层 "bgp"
metallb_protocol: "layer2"
metallb_offline: "metallb_v0.7.3.tar"
metallb_vip_pool: "192.168.1.240/29"  # 选一段与node节点相同网段的地址
```

- 2.执行安装 `ansible-playbook 07.cluster-addon.yml`，其中controller 负责统一loadbalancer地址管理和服务监控，speaker 负责节点的loadbalancer地址的对外宣告（使用arp或者bgp网络协议），注意 **speaker是以DaemonSet 形式运行且只会调度到有node-role.kubernetes.io/metallb-speaker=true标签的节点**，所以你可以选择做speaker的节点（该节点网络性能要好），使用命令 `$ kubectl label nodes 192.168.1.43 node-role.kubernetes.io/metallb-speaker=true`

- 3.验证metallb相关 pod

``` bash
$ kubectl get node
NAME           STATUS                     ROLES                  AGE       VERSION
192.168.1.41   Ready,SchedulingDisabled   master                 4h        v1.11.3
192.168.1.42   Ready                      node                   4h        v1.11.3
192.168.1.43   Ready                      metallb-speaker,node   4h        v1.11.3
192.168.1.44   Ready                      metallb-speaker,node   4h        v1.11.3
$ kubectl get pod -n metallb-system 
NAME                        READY     STATUS    RESTARTS   AGE
controller-9c57dbd4-798nb   1/1       Running   0          4h
speaker-9rjmk               1/1       Running   0          4h
speaker-n79l4               1/1       Running   0          4h
```

- 3.创建测试应用验证 loadbalancer 地址分配

``` bash
# 创建测试应用
$ cat > test-nginx.yaml << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx3
spec:
  selector:
    matchLabels:
      app: nginx3
  template:
    metadata:
      labels:
        app: nginx3
    spec:
      containers:
      - name: nginx3
        image: nginx:1
        ports:
        - name: http
          containerPort: 80

---
apiVersion: v1
kind: Service
metadata:
  name: nginx3
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: nginx3
  type: LoadBalancer
EOF
$ kubectl apply -f test-nginx.yaml

# 查看生成的loadbalancer 地址，如下验证成功
$ kubectl get svc
NAME         TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)        AGE
kubernetes   ClusterIP      10.68.0.1      <none>          443/TCP        5h
nginx3       LoadBalancer   10.68.82.227   192.168.1.240   80:38702/TCP   1m
```

- 4.验证使用loadbalacer 来暴露ingress的服务地址，之前在[ingress文档](ingress.md)中我们是使用nodeport方式服务类型，现在我们可以方便的使用loadbalancer类型了，使用loadbalancer地址(192.168.1.241)方便的绑定你要的域名进行访问。

``` bash
# 修改traefik-ingress 使用 LoadBalancer服务
$ sed -i 's/NodePort$/LoadBalancer/g' /etc/ansible/manifests/ingress/traefik/traefik-ingress.yaml
# 创建traefik-ingress
$ kubectl apply -f /etc/ansible/manifests/ingress/traefik/traefik-ingress.yaml
# 验证
$ kubectl get svc --all-namespaces |grep traefik
kube-system   traefik-ingress-service   LoadBalancer   10.68.163.243   192.168.1.241   80:23456/TCP,8080:37088/TCP   1m
```
