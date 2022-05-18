# Prometheus
随着`heapster`项目停止更新并慢慢被`metrics-server`取代，集群监控这项任务也将最终转移。`prometheus`的监控理念、数据结构设计其实相当精简，包括其非常灵活的查询语言；但是对于初学者来说，想要在k8s集群中实践搭建一套相对可用的部署却比较麻烦，由此还产生了不少专门的项目（如：[prometheus-operator](https://github.com/coreos/prometheus-operator)），本文介绍使用`helm chart`部署集群的prometheus监控。  
- `helm`已成为`CNCF`独立托管项目，预计会更加流行起来

## 前提

- 安装 helm
- 安装 [kube-dns](kubedns.md)

## 安装

项目3.x采用的部署charts: https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack

kubeasz 集成安装

- 1.修改 clusters/xxxx/config.yml 中配置项 prom_install: "yes"
- 2.安装 ezctl setup xxxx 07

注：涉及到镜像需从 quay.io 下载，国内比较慢，可以使用项目中的工具脚本 tools/imgutils

--- 以下内容暂未更新

## 验证安装

``` bash 
# 查看相关pod和svc
$ kubectl get pod,svc -n monitor
NAME                                                         READY   STATUS    RESTARTS   AGE
pod/alertmanager-prometheus-kube-prometheus-alertmanager-0   2/2     Running   0          3m11s
pod/prometheus-grafana-6d6d47996f-7xlpt                      2/2     Running   0          3m14s
pod/prometheus-kube-prometheus-operator-5f6774b747-bpktd     1/1     Running   0          3m14s
pod/prometheus-kube-state-metrics-95d956569-dhlkx            1/1     Running   0          3m14s
pod/prometheus-prometheus-kube-prometheus-prometheus-0       2/2     Running   1          3m11s
pod/prometheus-prometheus-node-exporter-d9m7j                1/1     Running   0          3m14s

NAME                                              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
service/alertmanager-operated                     ClusterIP   None            <none>        9093/TCP,9094/TCP,9094/UDP   3m12s
service/prometheus-grafana                        NodePort    10.68.31.225    <none>        80:30903/TCP                 3m14s
service/prometheus-kube-prometheus-alertmanager   NodePort    10.68.212.136   <none>        9093:30902/TCP               3m14s
service/prometheus-kube-prometheus-operator       NodePort    10.68.226.171   <none>        443:30900/TCP                3m14s
service/prometheus-kube-prometheus-prometheus     NodePort    10.68.100.42    <none>        9090:30901/TCP               3m14s
service/prometheus-kube-state-metrics             ClusterIP   10.68.80.70     <none>        8080/TCP                     3m14s
service/prometheus-operated                       ClusterIP   None            <none>        9090/TCP                     3m12s
service/prometheus-prometheus-node-exporter       ClusterIP   10.68.64.56     <none>        9100/TCP                     3m14s
```

- 访问prometheus的web界面：`http://$NodeIP:30901`
- 访问alertmanager的web界面：`http://$NodeIP:30902`
- 访问grafana的web界面：`http://$NodeIP:30903` (默认用户密码 admin:Admin1234!)

## 管理操作

## 验证告警

- 修改`prom-alertsmanager.yaml`文件中邮件告警为有效的配置内容，并使用 helm upgrade更新安装
- 手动临时关闭 master 节点的 kubelet 服务，等待几分钟看是否有告警邮件发送
 
``` bash
# 在 master 节点运行
$ systemctl stop kubelet
```

## [可选] 配置钉钉告警

- 创建钉钉群，获取群机器人 webhook 地址

使用钉钉创建群聊以后可以方便设置群机器人，【群设置】-【群机器人】-【添加】-【自定义】-【添加】，然后按提示操作即可，参考 https://open-doc.dingtalk.com/docs/doc.htm?spm=a219a.7629140.0.0.666d4a97eCG7XA&treeId=257&articleId=105735&docType=1

上述配置好群机器人，获得这个机器人对应的Webhook地址，记录下来，后续配置钉钉告警插件要用，格式如下

```
https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxx
```

- 创建钉钉告警插件，参考 http://theo.im/blog/2017/10/16/release-prometheus-alertmanager-webhook-for-dingtalk/

``` bash
# 编辑修改文件中 access_token=xxxxxx 为上一步你获得的机器人认证 token
$ vi /etc/ansible/manifests/prometheus/dingtalk-webhook.yaml
# 运行插件
$ kubectl apply -f /etc/ansible/manifests/prometheus/dingtalk-webhook.yaml
```

- 修改 alertsmanager 告警配置后，更新 helm prometheus 部署，成功后如上节测试告警发送

``` bash
# 修改 alertsmanager 告警配置
$ cd /etc/ansible/manifests/prometheus
$ vi prom-alertsmanager.yaml
# 增加 receiver dingtalk，然后在 route 配置使用 receiver: dingtalk
    receivers:
    - name: dingtalk
      webhook_configs:
      - send_resolved: false
        url: http://webhook-dingtalk.monitoring.svc.cluster.local:8060/dingtalk/webhook1/send
# ...
# 更新 helm prometheus 部署
$ helm upgrade --tls monitor -f prom-settings.yaml -f prom-alertsmanager.yaml -f prom-alertrules.yaml prometheus
```

## 下一步

- 继续了解prometheus查询语言和配置文件
- 继续了解prometheus告警规则，编写适合业务应用的告警规则
- 继续了解grafana的dashboard编写，本项目参考了部分[feisky的模板](https://grafana.com/orgs/feisky/dashboards)  
如果对以上部分有心得总结，欢迎分享贡献在项目中。
