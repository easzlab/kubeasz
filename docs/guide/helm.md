# Helm

`Helm`致力于成为k8s集群的应用包管理工具，希望像linux 系统的`RPM` `DPKG`那样成功；确实在k8s上部署复杂一点的应用很麻烦，需要管理很多yaml文件（configmap,controller,service,rbac,pv,pvc等等），而helm能够整齐管理这些文档：版本控制，参数化安装，方便的打包与分享等。  
- 建议积累一定k8s经验以后再去使用helm；对于初学者来说手工去配置那些yaml文件对于快速学习k8s的设计理念和运行原理非常有帮助，而不是直接去使用helm，面对又一层封装与复杂度。
- 本文基于helm 3（建议版本），helm 2 文档[请看这里](helm2.md)

## 安装 helm

在官方repo下载[release版本](https://github.com/helm/helm/releases)中自带的二进制文件即可（以Linux amd64为例）

```
wget https://get.helm.sh/helm-v3.2.1-linux-amd64.tar.gz
mv ./linux-amd64/helm /usr/bin
```

- 启用官方 charts 仓库

```
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
```

## 使用 helm 安装应用

helm3 安装命令与 helm2 稍有变化，个人习惯先下载对应charts到本地然后按照固定目录格式安装，以创建一个redis集群举例：

- 创建 redis-cluster 目录
``` bash
mkdir -p /opt/charts/redis-cluster
cd /opt/charts/redis-cluster
```

- 下载最新stalbe/redis-ha
```
helm repo update
helm pull stable/redis-ha
```

- 解压 charts，复制 values.yaml设置
```
tar zxvf redis-ha-*.tgz
cp redis-ha/values.yaml .
```

- 创建 start.sh 脚本记录启动命令
```
cat > start.sh << EOF
#!/bin/sh
set -x

ROOT=$(cd `dirname $0`; pwd)
cd $ROOT

helm install redis \
	--create-namespace \
	--namespace dependency \
	-f ./values.yaml \
	./redis-ha
EOF
```

- 查看当前目录结构如下
```
tree .
.
├── redis-ha		# redis-ha 原始charts目录
├── start.sh		# 启动命名脚本
└── values.yaml		# 个性化参数配置
```

- 修改当前目录的 values.yaml 为你的个性化配置
``` bash
#举例values.yaml 配置如下，没有启用PV
#cat values.yaml
image:
  repository: redis
  tag: 5.0.6-alpine

replicas: 2

## Redis specific configuration options
redis:
  port: 6379
  masterGroupName: "mymaster"       # must match ^[\\w-\\.]+$) and can be templated
  config:
    ## For all available options see http://download.redis.io/redis-stable/redis.conf
    min-replicas-to-write: 1
    min-replicas-max-lag: 5   # Value in seconds
    maxmemory: "4g"       # Max memory to use for each redis instance. Default is unlimited.
    maxmemory-policy: "allkeys-lru"  # Max memory policy to use for each redis instance. Default is volatile-lru.
    repl-diskless-sync: "yes"
    rdbcompression: "yes"
    rdbchecksum: "yes"

  resources:
    requests:
      memory: 200Mi
      cpu: 100m
    limits:
      memory: 4000Mi

## Sentinel specific configuration options
sentinel:
  port: 26379
  quorum: 1

  resources:
    requests:
      memory: 200Mi
      cpu: 100m
    limits:
      memory: 200Mi

hardAntiAffinity: true

## Configures redis with AUTH (requirepass & masterauth conf params)
auth: false

persistentVolume:
  enabled: false

hostPath:
  path: "/data/mcs-redis/{{ .Release.Name }}"
```

- 执行安装
```
bash ./start.sh
```

- 查看安装
```
helm ls -A
NAME 	NAMESPACE 	REVISION	UPDATED                                	STATUS  	CHART         	APP VERSION
redis	dependency	1       	2020-05-28 20:57:31.166002853 +0800 CST	deployed	redis-ha-4.4.4	5.0.6

# 查看k8s上资源
kubectl get pod,svc -n dependency
NAME                          READY   STATUS    RESTARTS   AGE
pod/redis-redis-ha-server-0   2/2     Running   0          119s
pod/redis-redis-ha-server-1   2/2     Running   0          104s

NAME                                TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)              AGE
service/redis-redis-ha              ClusterIP   None          <none>        6379/TCP,26379/TCP   119s
service/redis-redis-ha-announce-0   ClusterIP   10.68.41.65   <none>        6379/TCP,26379/TCP   119s
service/redis-redis-ha-announce-1   ClusterIP   10.68.64.49   <none>        6379/TCP,26379/TCP   119s
```

