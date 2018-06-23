## 03-安装docker服务.md

``` bash
roles/docker/
├── files
│   ├── daemon.json
│   ├── docker
│   └── docker-tag
├── tasks
│   └── main.yml
└── templates
    └── docker.service.j2
```

请在另外窗口打开[roles/docker/tasks/main.yml](../roles/docker/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建docker的systemd unit文件 

``` bash
[Unit]
Description=Docker Application Container Engine
Documentation=http://docs.docker.io

[Service]
Environment="PATH={{ bin_dir }}:/bin:/sbin:/usr/bin:/usr/sbin"
ExecStart={{ bin_dir }}/dockerd
ExecStartPost=/sbin/iptables -I FORWARD -s 0.0.0.0/0 -j ACCEPT
ExecReload=/bin/kill -s HUP $MAINPID
Restart=on-failure
RestartSec=5
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
Delegate=yes
KillMode=process

[Install]
WantedBy=multi-user.target
```
+ dockerd 运行时会调用其它 docker 命令，如 docker-proxy，所以需要将 docker 命令所在的目录加到 PATH 环境变量中；
+ docker 从 1.13 版本开始，将`iptables` 的`filter` 表的`FORWARD` 链的默认策略设置为`DROP`，从而导致 ping 其它 Node 上的 Pod IP 失败，因此必须在 `filter` 表的`FORWARD` 链增加一条默认允许规则 `iptables -I FORWARD -s 0.0.0.0/0 -j ACCEPT`
+ 运行`dockerd --help` 查看所有可配置参数，确保默认开启 `--iptables` 和 `--ip-masq` 选项

### 配置国内镜像加速

从国内下载docker官方仓库镜像非常缓慢，所以对于k8s集群来说配置镜像加速非常重要，配置 `/etc/docker/daemon.json`

``` bash
{
  "registry-mirrors": ["https://registry.docker-cn.com"],
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
    }
}
```

这将在后续部署calico下载 calico/node镜像和kubedns/heapster/dashboard镜像时起到重要加速效果。

由于K8S的官方镜像存放在`gcr.io`仓库，因此这个镜像加速对K8S的官方镜像没有效果；好在`Docker Hub`上有很多K8S镜像的转存，而`Docker Hub`上的镜像可以加速。这里推荐两个K8S镜像的`Docker Hub`项目,几乎能找到所有K8S相关的镜像，而且更新及时，感谢维护者的辛勤付出！

+ [mirrorgooglecontainers](https://hub.docker.com/u/mirrorgooglecontainers/)
+ [anjia0532](https://hub.docker.com/u/anjia0532/), [项目github地址](https://github.com/anjia0532/gcr.io_mirror)

当然对于企业内部应用的docker镜像，想要在K8S平台运行的话，特别是结合开发`CI/CD` 流程，肯定是需要部署私有镜像仓库的，后续会简单提到 `Harbor`的部署。

另外，daemon.json配置中也配置了docker 容器日志相关参数，设置单个容器日志超过10M则进行回卷，回卷的副本数超过3个就进行清理。

### 清理 iptables

因为后续`calico`网络、`kube-proxy`等将大量使用 iptables规则，安装前清空所有`iptables`策略规则；常见发行版`Ubuntu`的 `ufw` 和 `CentOS`的 `firewalld`等基于`iptables`的防火墙最好直接卸载，避免不必要的冲突。

``` bash
iptables -F && iptables -X \
        && iptables -F -t nat && iptables -X -t nat \
        && iptables -F -t raw && iptables -X -t raw \
        && iptables -F -t mangle && iptables -X -t mangle
```
+ calico 网络支持 `network-policy`，使用的`calico-kube-controllers` 会使用到`iptables` 所有的四个表 `filter` `nat` `raw` `mangle`，所以一并清理

### 启动 docker 略

### 可选-安装docker查询镜像 tag的小工具

docker官方目前没有提供在命令行直接查询某个镜像的tag信息的方式，网上找来一个脚本工具，使用很方便。

``` bash
$ docker-tag library/ubuntu
"14.04"
"16.04"
"17.04"
"latest"
"trusty"
"trusty-20171117"
"xenial"
"xenial-20171114"
"zesty"
"zesty-20171114"
$ docker-tag mirrorgooglecontainers/kubernetes-dashboard-amd64
"v0.1.0"
"v1.0.0"
"v1.0.0-beta1"
"v1.0.1"
"v1.1.0-beta1"
"v1.1.0-beta2"
"v1.1.0-beta3"
"v1.7.0"
"v1.7.1"
"v1.8.0"
``` 
+ 需要先apt安装轻量JSON处理程序 `jq`
+ 然后下载脚本即可使用
+ 脚本很简单，就一行命令如下

``` bash
#!/bin/bash
curl -s -S "https://registry.hub.docker.com/v2/repositories/$@/tags/" | jq '."results"[]["name"]' |sort
```
+ 对于 CentOS7 安装 `jq` 稍微费力一点，需要启用 `EPEL` 源

``` bash
wget http://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
rpm -ivh epel-release-latest-7.noarch.rpm
yum install jq
```

### 验证

运行`ansible-playbook 03.docker.yml` 成功后可以验证

``` bash
systemctl status docker 	# 服务状态
journalctl -u docker 		# 运行日志
docker version
docker info
```
`iptables-save|grep FORWARD` 查看 iptables filter表 FORWARD链，最后要有一个 `-A FORWARD -j ACCEPT` 保底允许规则

``` bash
iptables-save|grep FORWARD
:FORWARD ACCEPT [0:0]
:FORWARD DROP [0:0]
-A FORWARD -j DOCKER-USER
-A FORWARD -j DOCKER-ISOLATION
-A FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A FORWARD -o docker0 -j DOCKER
-A FORWARD -i docker0 ! -o docker0 -j ACCEPT
-A FORWARD -i docker0 -o docker0 -j ACCEPT
-A FORWARD -j ACCEPT
```

[前一篇](02-安装etcd集群.md) -- [后一篇](04-安装kube-master节点.md)
