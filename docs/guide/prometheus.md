# Prometheus
`prometheus`已经成为k8s集群上默认的监控解决方案，它的监控理念、数据结构设计其实相当精简，包括其非常灵活的查询语言；但是对于初学者来说，想要在k8s集群中实践搭建一套相对可用的部署却比较麻烦。本项目3.x采用的helm chart方式部署，使用的charts地址: https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack

## 安装

kubeasz 集成安装

- 1.修改 /etc/kubeasz/clusters/xxxx/config.yml 中配置项 prom_install: "yes"
- 2.下载镜像 /etc/kubeasz/ezdown -X
- 3.安装 /etc/kubeasz/ezctl setup xxxx 07

生成的charts自定义配置在/etc/kubeasz/clusters/xxxx/yml/prom-values.yaml

注1：如果需要修改配置，修改roles/cluster-addon/templates/prometheus/values.yaml.j2 后重新执行安装命令

注2：如果集群节点有增减，重新执行安装命令

注3：涉及到很多相关镜像下载比较慢，另外部分k8s.gcr.io的镜像已经替换成easzlab的mirror镜像地址

## 验证安装

``` bash 
# 查看相关pod和svc
$ kubectl get pod,svc -n monitor
NAME                                                         READY   STATUS    RESTARTS   AGE
pod/alertmanager-prometheus-kube-prometheus-alertmanager-0   2/2     Running   0          160m
pod/prometheus-grafana-69f88948bc-7hnbp                      3/3     Running   0          160m
pod/prometheus-kube-prometheus-operator-f8f4758cb-bm6gs      1/1     Running   0          160m
pod/prometheus-kube-state-metrics-74b8f49c6c-f9wgg           1/1     Running   0          160m
pod/prometheus-prometheus-kube-prometheus-prometheus-0       2/2     Running   0          160m
pod/prometheus-prometheus-node-exporter-6nfb4                1/1     Running   0          160m
pod/prometheus-prometheus-node-exporter-q4qq2                1/1     Running   0          160m

NAME                                              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
service/alertmanager-operated                     ClusterIP   None            <none>        9093/TCP,9094/TCP,9094/UDP   160m
service/prometheus-grafana                        NodePort    10.68.253.23    <none>        80:30903/TCP                 160m
service/prometheus-kube-prometheus-alertmanager   NodePort    10.68.125.191   <none>        9093:30902/TCP               160m
service/prometheus-kube-prometheus-operator       NodePort    10.68.161.218   <none>        443:30900/TCP                160m
service/prometheus-kube-prometheus-prometheus     NodePort    10.68.64.217    <none>        9090:30901/TCP               160m
service/prometheus-kube-state-metrics             ClusterIP   10.68.111.106   <none>        8080/TCP                     160m
service/prometheus-operated                       ClusterIP   None            <none>        9090/TCP                     160m
service/prometheus-prometheus-node-exporter       ClusterIP   10.68.252.83    <none>        9100/TCP                     160m
```

- 访问prometheus的web界面：`http://$NodeIP:30901`
- 访问alertmanager的web界面：`http://$NodeIP:30902`
- 访问grafana的web界面：`http://$NodeIP:30903` (默认用户密码 admin:Admin1234!)

## 其他操作

-- 以下内容没有更新测试

### [可选] 配置钉钉告警

- 创建钉钉群，获取群机器人 webhook 地址

使用钉钉创建群聊以后可以方便设置群机器人，【群设置】-【群机器人】-【添加】-【自定义】-【添加】，然后按提示操作即可，参考 https://open.dingtalk.com/document/group/custom-robot-access

上述配置好群机器人，获得这个机器人对应的Webhook地址，记录下来，后续配置钉钉告警插件要用，格式如下

```
https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxx
```

- 创建钉钉告警插件，参考:
  - https://github.com/timonwong/prometheus-webhook-dingtalk
  - http://theo.im/blog/2017/10/16/release-prometheus-alertmanager-webhook-for-dingtalk/

``` bash
# 编辑修改文件中 access_token=xxxxxx 为上一步你获得的机器人认证 token
$ vi /etc/kubeasz/roles/cluster-addon/templates/prometheus/dingtalk-webhook.yaml
# 运行插件
$ kubectl apply -f /etc/kubeasz/roles/cluster-addon/templates/prometheus/dingtalk-webhook.yaml
```

- 修改 alertsmanager 告警配置，重新运行安装命令/etc/kubeasz/ezctl setup xxxx 07，成功后如上节测试告警发送

``` bash
# 修改 alertsmanager 告警配置
$ vi /etc/kubeasz/roles/cluster-addon/templates/prometheus/values.yaml.j2 
# 增加 receiver dingtalk，然后在 route 配置使用 receiver: dingtalk
    receivers:
    - name: dingtalk
      webhook_configs:
      - send_resolved: false
        url: http://webhook-dingtalk.monitor.svc.cluster.local:8060/dingtalk/webhook1/send
# ...
```
