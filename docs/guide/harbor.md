# harbor 镜像仓库

Habor是由VMWare中国团队开源的企业级容器镜像仓库。特性包括：友好的用户界面，基于角色的访问控制，水平扩展，同步复制，AD/LDAP集成以及审计日志等。本文档仅说明单机安装harbor 服务。

- 目录
  - 安装步骤
  - 安装讲解
  - 配置docker/containerd信任harbor证书
  - 在k8s集群使用harbor
  - 管理维护

### 安装步骤

1. 下载离线安装包，成功后在/etc/kubeasz/down/目录下有离线包harbor-offline-installer-$HARBOR_VER.tgz

```
ezdown -D
ezdown -R
```

2. 利用ezctl [文档](../setup/ezctl.md) 创建一个新的集群，已有集群修改同样的文件

```
#clusters/xxx/hosts 中修改如下，配置harbor组下机器，设置NEW_INSTALL=true
...
# 'NEW_INSTALL': 'true' to install a harbor server; 'false' to integrate with existed one
[harbor]
192.168.1.8 NEW_INSTALL=true
...

#clusters/xxx/config.yml 中修改如下，按需修改HARBOR_DOMAIN/HARBOR_TLS_PORT 等配置项
############################
# role:harbor
############################
# harbor version，完整版本号
HARBOR_VER: "v2.1.3"
HARBOR_DOMAIN: "harbor.yourdomain.com"
HARBOR_TLS_PORT: 8443

# if set 'false', you need to put certs named harbor.pem and harbor-key.pem in directory 'down'
HARBOR_SELF_SIGNED_CERT: true

# install component
HARBOR_WITH_NOTARY: false
HARBOR_WITH_TRIVY: false
HARBOR_WITH_CHARTMUSEUM: true
```

3. 配置完成后，执行 `./ezctl setup xxx harbor`，完成harbor安装和docker 客户端配置

- 安装验证

1. 在harbor节点使用`docker ps -a` 查看harbor容器组件运行情况
2. 浏览器访问地址（忽略证书报错） `https://${HARBOR_DOMAIN}:${HARBOR_TLS_PORT}`，管理员账号是 admin ，密码见harbor.yml文件 harbor_admin_password 对应值（默认密码 Harbor12345 已被随机生成的16位随机密码替换，不然存在安全隐患)

### 安装讲解

根据`playbooks/11.harbor.yml`文件，harbor节点需要以下步骤：

- role `os-harden` 系统安全加固（可选）
- role `chrony` 时间同步服务（可选）
- role `prepare` 基础系统环境准备
- role `docker` 安装docker
- role `harbor` 安装harbor
- 注意：`kube_node`节点在harbor部署完之后，需要配置harbor的证书（详见下节配置docker/containerd信任harbor证书），并可以在hosts里面添加harbor的域名解析，如果你的环境中有dns服务器，可以跳过hosts文件设置

1. 下载docker-compose可执行文件到$PATH目录
1. 自注册变量result判断是否已经安装harbor，避免重复安装问题
1. 解压harbor离线安装包到指定目录
1. 导入harbor所需 docker images
1. 创建harbor证书和私钥(复用集群的CA证书)
1. 修改harbor.yml配置文件
1. 启动harbor安装脚本

### 配置docker/containerd信任harbor证书

因为我们创建的harbor仓库使用了自签证书，所以当docker/containerd客户端拉取自建harbor仓库镜像前必须配置信任harbor证书，否则出现如下错误：

```
# docker
$ docker pull harbor.test.lo/pub/hello:v0.1.4
Error response from daemon: Get https://harbor.test.lo/v1/_ping: x509: certificate signed by unknown authority

# containerd
$ crictl pull harbor.test.lo/pub/hello:v0.1.4
FATA[0000] pulling image failed: rpc error: code = Unknown desc = failed to resolve image "harbor.test.lo/pub/hello:v0.1.4": no available registry endpoint: failed to do request: Head https://harbor.test.lo/v2/pub/hello/manifests/v0.1.4: x509: certificate signed by unknown authority
```

项目脚本`11.harbor.yml`中已经自动为k8s集群的每个node节点配置 docker/containerd 信任自建 harbor 证书；如果你无法运行此脚本，可以参考下述手工配置

#### docker配置信任harbor证书

在集群每个 node 节点进行如下配置

- 创建目录 /etc/docker/certs.d/harbor.yourdomain.com:8443/  (以默认配置举例)
- 复制 harbor 安装时的 CA 证书到上述目录，并改名 ca.crt 即可

#### containerd配置信任harbor证书

在集群每个 node 节点进行如下配置（假设ca.pem为自建harbor的CA证书）

- ubuntu 1604:
  - cp ca.pem /usr/share/ca-certificates/harbor-ca.crt
  - echo harbor-ca.crt >> /etc/ca-certificates.conf
  - update-ca-certificates

- CentOS 7:
  - cp ca.pem /etc/pki/ca-trust/source/anchors/harbor-ca.crt
  - update-ca-trust

上述配置完成后，重启 containerd 即可 `systemctl restart containerd`

### 在k8s集群使用harbor

admin用户web登录后可以方便的创建项目，并指定项目属性(公开或者私有)；然后创建用户，并在项目`成员`选项中选择用户和权限；

#### 镜像上传

在node上使用harbor私有镜像仓库首先需要在指定目录配置harbor的CA证书，详见 `11.harbor.yml`文件。

使用docker客户端登录`harbor.test.com`，然后把镜像tag成 `harbor.test.com/$项目名/$镜像名:$TAG` 之后，即可使用docker push 上传

``` bash
docker login harbor.test.com
Username: 
Password:
Login Succeeded
docker tag busybox:latest harbor.test.com/library/busybox:latest
docker push harbor.test.com/library/busybox:latest
The push refers to a repository [harbor.test.com/library/busybox]
0271b8eebde3: Pushed 
latest: digest: sha256:91ef6c1c52b166be02645b8efee30d1ee65362024f7da41c404681561734c465 size: 527
```
#### k8s中使用harbor

1. 如果镜像保存在harbor中的公开项目中，那么只需要在yaml文件中简单指定harbor私有镜像即可，例如

``` bash
apiVersion: v1
kind: Pod
metadata:
  name: test-busybox
spec:
  containers:
  - name: test-busybox
    image: harbor.test.com/xxx/busybox:latest
    imagePullPolicy: Always
```

2. 如果镜像保存在harbor中的私有项目中，那么yaml文件中使用该私有项目的镜像需要指定`imagePullSecrets`，例如

``` bash
apiVersion: v1
kind: Pod
metadata:
  name: test-busybox
spec:
  containers:
  - name: test-busybox
    image: harbor.test.com/xxx/busybox:latest
    imagePullPolicy: Always
  imagePullSecrets:
  - name: harborkey1
```
其中 `harborKey1`可以用以下两种方式生成：

+ 1.使用 `kubectl create secret docker-registry harborkey1 --docker-server=harbor.test.com --docker-username=admin --docker-password=Harbor12345 --docker-email=team@test.com`
+ 2.使用yaml配置文件生成 

``` bash
//harborkey1.yaml
apiVersion: v1
kind: Secret
metadata:
  name: harborkey1
  namespace: default
data:
    .dockerconfigjson: {base64 -w 0 ~/.docker/config.json}
type: kubernetes.io/dockerconfigjson
```
前面docker login会在~/.docker下面创建一个config.json文件保存鉴权串，这里secret yaml的.dockerconfigjson后面的数据就是那个json文件的base64编码输出（-w 0让base64输出在单行上，避免折行）

### 管理维护

+ 日志目录 `/var/log/harbor`
+ 数据目录 `/var/data` ，其中最主要是 `/var/data/database` 和 `/var/data/registry` 目录，如果你要彻底重新安装harbor，删除这两个目录即可

先进入harbor安装目录 `cd /var/data/harbor`，常规操作如下：

1. 暂停harbor `docker-compose stop` : docker容器stop，并不删除容器
2. 恢复harbor `docker-compose start` : 恢复docker容器运行
3. 停止harbor `docker-compose down -v` : 停止并删除docker容器
4. 启动harbor `docker-compose up -d` : 启动所有docker容器

修改harbor的运行配置，需要如下步骤：

``` bash
# 停止 harbor
 docker-compose down -v
# 修改配置
 vim harbor.cfg
# 执行./prepare已更新配置到docker-compose.yml文件
 ./prepare
# 启动 harbor
 docker-compose up -d
```
