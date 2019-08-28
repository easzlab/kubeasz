# Prometheus
随着`heapster`项目停止更新并慢慢被`metrics-server`取代，集群监控这项任务也将最终转移。`prometheus`的监控理念、数据结构设计其实相当精简，包括其非常灵活的查询语言；但是对于初学者来说，想要在k8s集群中实践搭建一套相对可用的部署却比较麻烦，由此还产生了不少专门的项目（如：[prometheus-operator](https://github.com/coreos/prometheus-operator)），本文介绍使用`helm chart`部署集群的prometheus监控。  
- `helm`已成为`CNCF`独立托管项目，预计会更加流行起来

## 前提

- 安装 helm：以本项目[安全安装helm](helm.md)为例
- 安装 [kube-dns](kubedns.md)

## 准备

安装目录概览 `ll /etc/ansible/manifests/prometheus`

``` bash
drwx------  3 root root  4096 Jun  3 22:42 grafana/
-rw-r-----  1 root root 67875 Jun  4 22:47 grafana-dashboards.yaml
-rw-r-----  1 root root   690 Jun  4 09:34 grafana-settings.yaml
-rw-r-----  1 root root  1105 May 30 16:54 prom-alertrules.yaml
-rw-r-----  1 root root   474 Jun  5 10:04 prom-alertsmanager.yaml
drwx------  3 root root  4096 Jun  2 21:39 prometheus/
-rw-r-----  1 root root   294 May 30 18:09 prom-settings.yaml
```
- 目录`prometheus/`和`grafana/`即官方的helm charts，可以使用`helm fetch --untar stable/prometheus` 和 `helm fetch --untar stable/grafana`下载，本安装不会修改任何官方charts里面的内容，这样方便以后跟踪charts版本的更新
- `prom-settings.yaml`：个性化prometheus安装参数，比如禁用PV，禁用pushgateway，设置nodePort等
- `prom-alertrules.yaml`：配置告警规则
- `prom-alertsmanager.yaml`：配置告警邮箱设置等
- `grafana-settings.yaml`：个性化grafana安装参数，比如用户名密码，datasources，dashboardProviders等
- `grafana-dashboards.yaml`：预设置dashboard

## 安装

``` bash
$ source ~/.bashrc
$ cd /etc/ansible/manifests/prometheus
# 安装 prometheus chart，如果你的helm安装没有启用tls证书，请忽略--tls参数
$ helm install --tls \
        --name monitor \
        --namespace monitoring \
        -f prom-settings.yaml \
        -f prom-alertsmanager.yaml \
        -f prom-alertrules.yaml \
        prometheus
# 安装 grafana chart
$ helm install --tls \
	--name grafana \
	--namespace monitoring \
	-f grafana-settings.yaml \
	-f grafana-dashboards.yaml \
	grafana
```

## 验证安装

``` bash 
# 查看相关pod和svc
$ kubectl get pod,svc -n monitoring 
NAME                                                     READY     STATUS    RESTARTS   AGE
grafana-54dc76d47d-2mk55                                 1/1       Running   0          1m
monitor-prometheus-alertmanager-6d9d9b5b96-w57bk         2/2       Running   0          2m
monitor-prometheus-kube-state-metrics-69f5d56f49-fh9z7   1/1       Running   0          2m
monitor-prometheus-node-exporter-55bwx                   1/1       Running   0          2m
monitor-prometheus-node-exporter-k8sb2                   1/1       Running   0          2m
monitor-prometheus-node-exporter-kxlr9                   1/1       Running   0          2m
monitor-prometheus-node-exporter-r5dx8                   1/1       Running   0          2m
monitor-prometheus-server-5ccfc77dff-8h9k6               2/2       Running   0          2m

NAME                                    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
grafana                                 NodePort    10.68.74.242   <none>        80:39002/TCP   1m
monitor-prometheus-alertmanager         NodePort    10.68.69.105   <none>        80:39001/TCP   2m
monitor-prometheus-kube-state-metrics   ClusterIP   None           <none>        80/TCP         2m
monitor-prometheus-node-exporter        ClusterIP   None           <none>        9100/TCP       2m
monitor-prometheus-server               NodePort    10.68.248.94   <none>        80:39000/TCP   2m
```

- 访问prometheus的web界面：`http://$NodeIP:39000`
- 访问alertmanager的web界面：`http://$NodeIP:39001`
- 访问grafana的web界面：`http://$NodeIP:39002` (默认用户密码 admin:admin，可在web界面修改)

## 管理操作

- 升级（修改配置）：修改配置请在`prom-settings.yaml` `prom-alertsmanager.yaml` 等文件中进行，保存后执行：  
``` bash
# 修改prometheus
$ helm upgrade --tls monitor -f prom-settings.yaml -f prom-alertsmanager.yaml -f prom-alertrules.yaml prometheus
# 修改grafana
$ helm upgrade --tls grafana -f grafana-settings.yaml -f grafana-dashboards.yaml grafana
```
- 回退：具体可以参考`helm help rollback`文档
``` bash
$ helm rollback --tls monitor [REVISION]
```
- 删除 
``` bash
$ helm del --tls monitor --purge
$ helm del --tls grafana --purge
```

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
