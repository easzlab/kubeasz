# 容器化 GO 应用

Golang 作为服务器端新兴热门语言同时也是容器技术的主要编写语言备受关注；它简洁、有趣、并行、安全等特点让 GO 应用容器化相对省心；一般来说做下时间本地化、安装信任根证书，然后把编译生成的二进制拷贝进去即可。

## 一个演示 GO WEB 应用

[hellogo 代码](hellogo.go)

## Dockerfile

作为演示项目的Dockerfile比较简单，请看 [Dockerfile 文件](Dockerfile)

- 采用 docker 多阶段编译，使生成的目标镜像最小
- 使用 alpine 基础镜像
- 安装 tzdata 做时间本地化
- 安装信任根证书

一个真实复杂go项目的Dockerfile可能如这个例子：[复杂 Dockerfile](Dockerfile-more)

## 制作镜像

在 Dockerfile 文件所在目录，执行

```
docker build -t hellogo:v1.0 .
```

## 本地测试应用

- 1.单机运行 hellogo 容器应用 

```
docker run -d --name hello -p3000:3000 hellogo:v1.0
```

- 2.验证测试

``` bash
# 查看本地监听端口
$ ss -ntl|grep 3000
LISTEN   0         128                       *:3000                   *:*

# 查看应用状态
$ curl localhost:3000
Hello, Go! I'm instance 987 running version 1.2 at 13109-10-13 08:39:11

$ curl localhost:3000/health -i
HTTP/1.1 200 OK
Date: Sun, 13 Oct 2019 00:39:15 GMT
Content-Length: 0

$ curl localhost:3000/version
1.2
```

## 在 k8s 上运行演示应用

- 可以参考项目`github.com/easzlab/kubeasz` 快速搭建一个本地 k8s 测试环境

- 1.编写基于k8s的应用编排文件 [hellogo.yaml](hellogo.yaml)
  - 设置应用副本数`replicas: 3`
  - 预设新副本启动延迟5秒`minReadySeconds: 5`
  - 设置滚动更新策略
  - 设置资源使用限制，安装实际情况修改
  - 设置服务对外暴露方式 NodePort，根据实际情况修改端口，或者使用 ingress 方式

- 2.在 k8s 上运行应用

``` bash
# 运行
$ kubectl apply -f hellogo.yaml

# 验证
$ kubectl get pod
NAME                             READY   STATUS    RESTARTS   AGE
hellogo-deploy-854dcd85c-2zm9l   1/1     Running   0          12m
hellogo-deploy-854dcd85c-7nfk5   1/1     Running   0          12m
hellogo-deploy-854dcd85c-ns7fp   1/1     Running   0          12m

$kubectl get deploy
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
hellogo-deploy   3/3     3            3           13m

$kubectl get svc
NAME          TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
hellogo-svc   NodePort    10.68.194.109   <none>        80:30000/TCP   13m

# 使用curl测试应用三副本状态（用curl多次访问看到三个不同`instance id`）
$ curl http://192.168.111.3:30000
Hello, Go! I'm instance 629 running version 1.2 at 13109-10-13 09:06:25

$ curl http://192.168.111.3:30000
Hello, Go! I'm instance 722 running version 1.2 at 13109-10-13 09:06:27

$curl http://192.168.111.3:30000
Hello, Go! I'm instance 799 running version 1.2 at 13109-10-13 09:06:28
```
