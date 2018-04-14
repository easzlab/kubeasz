# 更改高可用 `Master IP`

**WARNING:** 更改集群的 `Master VIP`操作有风险，不建议在生产环境直接操作，此文档实践一个修改的操作流程，帮助理解整个集群运行架构和 `kubeasz`的部署逻辑，请在测试环境操作练手。
**BUG:** 目前该操作只适用于集群网络选用`calico`，如果使用`flannel`操作变更后会出现POD地址分配错误的BUG。

首先分析大概操作思路：

- 修改`/etc/ansible/hosts`里面的配置项`MASTER_IP` `KUBE_APISERVER`
- 修改LB节点的keepalive的配置，重启keepalived服务
- 修改kubectl/kube-proxy的配置文件，使用新VIP地址更新api-server地址
- 重新生成master证书，hosts字段包含新VIP地址
- 修改kubelet的配置文件（kubelet的配置文件和证书是由bootstrap机制自动生成的）
  - 删除kubelet.kubeconfig
  - 删除集群所有node 节点
  - 所有节点重新bootstrap

## 变更前状态验证

``` bash
$ kubectl get cs,node,pod -o wide
NAME                 STATUS    MESSAGE             ERROR
controller-manager   Healthy   ok                  
scheduler            Healthy   ok                  
etcd-2               Healthy   {"health":"true"}   
etcd-0               Healthy   {"health":"true"}   
etcd-1               Healthy   {"health":"true"}   

NAME           STATUS                     ROLES     AGE       VERSION   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
192.168.1.41   Ready,SchedulingDisabled   <none>    2h        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.42   Ready,SchedulingDisabled   <none>    2h        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.43   Ready                      <none>    2h        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.44   Ready                      <none>    2h        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-98-generic   docker://18.3.0
192.168.1.45   Ready                      <none>    2h        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-98-generic   docker://18.3.0

NAME                     READY     STATUS    RESTARTS   AGE       IP               NODE
busy-5d6b6b5d4b-8wxkp    1/1       Running   0          17h       172.20.135.133   192.168.1.41
busy-5d6b6b5d4b-fcmkp    1/1       Running   0          17h       172.20.135.128   192.168.1.41
busy-5d6b6b5d4b-ptvd7    1/1       Running   0          17h       172.20.135.136   192.168.1.41
nginx-768979984b-ncqbp   1/1       Running   0          17h       172.20.135.137   192.168.1.41

# 查看待变更集群 Master VIP
$ kubectl cluster-info 
Kubernetes master is running at https://192.168.1.39:8443
```

## 变更操作

- `ansible playbook`可以使用tags来控制只允许部分任务执行，这里为简化操作没有细化，在deploy节点具体操作如下：

``` bash
# 1.修改/etc/ansible/hosts 配置项MASTER_IP，KUBE_APISERVER

# 2.删除集群所有node节点，等待重新bootstrap
$ kubectl get node |grep Ready|awk '{print $1}' |xargs kubectl delete node

# 3.重置keepalived 和修改kubectl/kube-proxy/bootstrap配置
$ ansible-playbook 01.prepare.yml

# 4.删除旧master证书
$ ansible kube-master -m file -a 'path=/etc/kubernetes/ssl/kubernetes.pem state=absent'

# 5.删除旧kubelet配置文件
$ ansible all -m file -a 'path=/etc/kubernetes/kubelet.kubeconfig state=absent'

# 6.重新配置启动master节点
$ ansible-playbook 04.kube-master.yml

# 7.重新配置启动node节点
$ ansible-playbook 05.kube-node.yml
```

## 变更后验证

``` bash
$ kubectl get cs,node,pod -o wide
NAME                 STATUS    MESSAGE             ERROR
scheduler            Healthy   ok                  
controller-manager   Healthy   ok                  
etcd-2               Healthy   {"health":"true"}   
etcd-1               Healthy   {"health":"true"}   
etcd-0               Healthy   {"health":"true"}   

NAME           STATUS                     ROLES     AGE       VERSION   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
192.168.1.41   Ready,SchedulingDisabled   <none>    4m        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.42   Ready,SchedulingDisabled   <none>    4m        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.43   Ready                      <none>    3m        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-97-generic   docker://18.3.0
192.168.1.44   Ready                      <none>    3m        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-98-generic   docker://18.3.0
192.168.1.45   Ready                      <none>    3m        v1.10.0   <none>        Ubuntu 16.04.3 LTS   4.4.0-98-generic   docker://18.3.0

NAME                     READY     STATUS    RESTARTS   AGE       IP               NODE
busy-5d6b6b5d4b-25hfr    1/1       Running   0          5m        172.20.237.64    192.168.1.43
busy-5d6b6b5d4b-cdzb5    1/1       Running   0          5m        172.20.145.192   192.168.1.44
busy-5d6b6b5d4b-m2rf7    1/1       Running   0          5m        172.20.26.131    192.168.1.45
nginx-768979984b-2ngww   1/1       Running   0          5m        172.20.145.193   192.168.1.44

# 查看集群master VIP已经变更 
$ kubectl cluster-info 
Kubernetes master is running at https://192.168.1.40:8443
```

## 小结

本示例操作演示了多主多节点k8s集群变更`Master VIP`的操作，有助于理解整个集群组件架构和`kubeasz`的安装逻辑，小结如下：

- 变更操作不影响集群已运行的业务POD，但是操作过程中业务会中断
- 已运行POD会重新调度到各node节点，如果业务POD量很大，短时间内会对集群造成压力
- 不建议在生成环境直接操作，本示例演示说明为主
