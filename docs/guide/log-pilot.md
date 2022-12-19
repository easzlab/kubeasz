# Log-Pilot Elasticsearch Kibana 日志解决方案

该方案是社区方案`EFK`的升级版，它支持两种搜集形式，对应容器标准输出日志和容器内的日志文件；个人使用了一把，在原有`EFK`经验的基础上非常简单、方便，值得推荐；更多的关于`log-pilot`的介绍详见链接：

- github 项目地址: https://github.com/AliyunContainerService/log-pilot
- 阿里云介绍文档: https://help.aliyun.com/document_detail/86552.html
- 介绍文档2: https://yq.aliyun.com/articles/674327

## 安装步骤

- 1.安装 ES 集群，同[EFK](efk.md)文档

- 2.安装 Kibana，同[EFK](efk.md)文档

- 3.安装 Log-Pilot

``` bash
kubectl apply -f /etc/kubeasz/manifests/efk/log-pilot/log-pilot-filebeat.yaml
```

- 4.创建示例应用，采集日志

``` bash
$ cat > tomcat.yaml << EOF
apiVersion: v1
kind: Pod
metadata:
  name: tomcat
spec:
  containers:
  - name: tomcat
    image: "tomcat:7.0"
    env:
    # 1、stdout为约定关键字，表示采集标准输出日志
    # 2、配置标准输出日志采集到ES的catalina索引下
    - name: aliyun_logs_catalina
      value: "stdout"
    # 1、配置采集容器内文件日志，支持通配符
    # 2、配置该日志采集到ES的access索引下
    - name: aliyun_logs_access
      value: "/usr/local/tomcat/logs/catalina.*.log"
    volumeMounts:
      - name: tomcat-log
        mountPath: /usr/local/tomcat/logs
  volumes:
    # 容器内文件日志路径需要配置emptyDir
    - name: tomcat-log
      emptyDir: {}
EOF

$ kubectl apply -f tomcat.yaml 
```

- 5.在 kibana 创建 Index Pattern，验证日志已搜集，如上示例应用，应创建如下 index pattern
  - catalina-*
  - access-*
