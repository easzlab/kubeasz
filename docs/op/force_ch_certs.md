# 强制更新CA和所有证书

- WARNNING: 此命令使用需要小心谨慎，确保了解功能背景和可能的结果；执行后，它会重新创建集群CA证书以及由它颁发的所有其他证书；一般适合于集群admin.conf不小心泄露，为了避免集群被非法访问，重新创建CA，从而使已泄漏的admin.conf失效。

- 如果需要分发受限的kubeconfig，强烈建议使用[自定义权限和期限的kubeconfig](kcfg-adm.md)

## 使用帮助

确认需要强制更新后，在ansible 控制节点使用如下命令：(xxx 表示需要操作的集群名)

``` bash
docker exec -it kubeasz ezctl kca-renew xxx
# 或者使用 dk ezctl kca-renew xxx
```

上述命令执行后，按序进行以下的操作：详见`playbooks/96.update-certs.yml`

- 重新生成CA证书，以及各种kubeconfig
- 签发新etcd证书，并使用新证书重启etcd服务
- 签发新kube-apiserver 证书，并重启kube-apiserver/kube-controller-manager/kube-scheduler 服务
- 签发新kubelet 证书，并重启kubelet/kube-proxy 服务
- 重启网络组件pod
- 重启其他集群组件pod

- **特别注意：** 如果集群中运行的业务负载pod需要访问apiserver，需要重启这些pod

## 检查验证

更新完毕，注意检查集群组件日志和容器pod日志，确认集群处于正常状态

- 集群组件日志：使用journalctl -u xxxx.service -f 依次检查 etcd.service/kube-apiserver.service/kube-controller-manager.service/kube-scheduler.service/kubelet.service/kube-proxy.service
- 容器pod日志：使用 kubectl logs 方式检查容器日志
