# network-check

网络测试组件，根据cilium connectivity-check 脚本修改而来；利用cronjob 定期检测集群各节点、容器、serviceip、nodeport等之间的网络联通性；可以方便的判断当前集群网络是否正常。

目前检测如下：

``` bash
kubectl get cronjobs.batch -n network-test
NAME                                  SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
test01-pod-to-container               */5 * * * *   False     0        3m19s           6d3h
test02-pod-to-node-nodeport           */5 * * * *   False     0        3m19s           6d3h
test03-pod-to-multi-node-clusterip    */5 * * * *   False     1        6d3h            6d3h
test04-pod-to-multi-node-headless     */5 * * * *   False     1        6d3h            6d3h
test05-pod-to-multi-node-nodeport     */5 * * * *   False     1        6d3h            6d3h
test06-pod-to-external-1111           */5 * * * *   False     0        3m19s           6d3h
test07-pod-to-external-fqdn-baidu     */5 * * * *   False     0        3m19s           6d3h
test08-host-to-multi-node-clusterip   */5 * * * *   False     1        6d3h            6d3h
test09-host-to-multi-node-headless    */5 * * * *   False     1        6d3h            6d3h
```

+ 带`multi-node`的测试需要多节点集群才能运行，如果单节点集群，测试pod会处于`Pending`状态
+ 带`external`的测试需要节点能够访问互联网，否则测试会失败

## 启用网络检测

- 下载额外容器镜像 `./ezdown -X`

- 配置集群，在配置文件`/etc/kubeasz/clusters/xxx/config.yml` (xxx为集群名) 修改如下选项

```
# network-check 自动安装
network_check_enabled: true
network_check_schedule: "*/5 * * * *"  # 检测频率，默认5分钟执行一次
```

- 安装网络检测插件 `docker exec -it kubeasz ezctl setup xxx 07`

## 检查测试结果

大约等待5分钟左右，查看运行结果，如果pod 状态为`Completed` 表示检测正常通过。

```
kubectl get pod -n network-test
NAME                                                 READY   STATUS      RESTARTS   AGE
echo-server-58d7bb7f6-77ps6                          1/1     Running     0          6d4h
echo-server-host-cc87c966d-bk57t                     1/1     Running     0          6d4h
test01-pod-to-container-27606775-q6xlb               0/1     Completed   0          3m10s
test02-pod-to-node-nodeport-27606775-x2v5d           0/1     Completed   0          3m10s
test03-pod-to-multi-node-clusterip-27597895-cbq8d    0/1     Pending     0          6d4h
test04-pod-to-multi-node-headless-27597895-qzsgz     0/1     Pending     0          6d4h
test05-pod-to-multi-node-nodeport-27597895-kb5r7     0/1     Pending     0          6d4h
test06-pod-to-external-1111-27606775-p6v8s           0/1     Completed   0          3m10s
test07-pod-to-external-fqdn-baidu-27606775-qdfwd     0/1     Completed   0          3m10s
test08-host-to-multi-node-clusterip-27597895-qsgn9   0/1     Pending     0          6d4h
test09-host-to-multi-node-headless-27597895-hpkt5    0/1     Pending     0          6d4h
```

+ pod 状态为`Completed` 表示检测正常通过
+ pod 状态为`Pending` 表示该检测需要多节点的k8s集群才会运行

## 禁用网络检测

如果集群已经开启网络检测，检测结果符合预期，并且不想继续循环检测时，只要删除对应namespace即可

```
kubectl delete ns network-test
```
