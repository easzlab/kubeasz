## heapster
### 部署文件
1. [grafana](./grafana.yaml)
1. [heapster](./heapster.yaml)
1. [influxdb](./influxdb.yaml)

### tips-01
如果部署完heapster，检查状态均正常，但是dashboard不能展示 Pod、Nodes 的 CPU、内存等 metric 图形，请重启dashboard 容器
检查命令：
``` bash
$ kubectl get pods -n kube-system | grep -E 'heapster|monitoring'
heapster-3273315324-tmxbg               1/1       Running   0          11m
monitoring-grafana-2255110352-94lpn     1/1       Running   0          11m
monitoring-influxdb-884893134-3vb6n     1/1       Running   0          11m
```
检查Pods日志：
``` bash
$ kubectl logs heapster-3273315324-tmxbg -n kube-system
$ kubectl logs monitoring-grafana-2255110352-94lpn -n kube-system
$ kubectl logs monitoring-influxdb-884893134-3vb6n -n kube-system
```
