## 03-安装容器运行时（docker or containerd）

目前k8s官方推荐使用containerd，查阅[使用文档](containerd.md)

## 安装docker服务

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

### 配置daemon.json

roles/docker/templates/daemon.json.j2

``` bash
{
  "data-root": "{{ DOCKER_STORAGE_DIR }}",
  "exec-opts": ["native.cgroupdriver=cgroupfs"],
{% if ENABLE_MIRROR_REGISTRY %}
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "http://hub-mirror.c.163.com"
  ],
{% endif %}
{% if ENABLE_REMOTE_API %}
  "hosts": ["tcp://0.0.0.0:2376", "unix:///var/run/docker.sock"],
{% endif %}
  "insecure-registries": {{ INSECURE_REG }},
  "max-concurrent-downloads": 10,
  "live-restore": true,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "50m",
    "max-file": "1"
    },
  "storage-driver": "overlay2"
}
```
- data-root 配置容器数据目录，默认/var/lib/docker，在集群安装时要规划磁盘空间使用
- registry-mirrors 配置国内镜像仓库加速
- live-restore 可以重启docker daemon ，而不重启容器
- log-opts 容器日志相关参数，设置单个容器日志超过50M则进行回卷，回卷的副本数超过1个就进行清理

对于企业内部应用的docker镜像，想要在K8S平台运行的话，特别是结合开发`CI/CD` 流程，需要部署私有镜像仓库，参阅[harbor文档](../guide/harbor.md)。

### 清理 iptables

因为`calico`网络、`kube-proxy`等将大量使用 iptables规则，安装前清空所有`iptables`策略规则；常见发行版`Ubuntu`的 `ufw` 和 `CentOS`的 `firewalld`等基于`iptables`的防火墙最好直接卸载，避免不必要的冲突。

WARNNING: 如果有自定义的iptables规则也会被一并清除，如果一定要使用自定义规则，可以集群安装完成后在应用规则

``` bash
iptables -F && iptables -X \
        && iptables -F -t nat && iptables -X -t nat \
        && iptables -F -t raw && iptables -X -t raw \
        && iptables -F -t mangle && iptables -X -t mangle
```
+ calico 网络支持 `network-policy`，使用的`calico-kube-controllers` 会使用到`iptables` 所有的四个表 `filter` `nat` `raw` `mangle`，所以一并清理

### 可选-安装docker查询镜像 tag的小工具

docker官方没有提供在命令行直接查询某个镜像的tag信息的方式，可以使用一个工具脚本：

``` bash
$ docker-tag library/ubuntu
"14.04"
"16.04"
"17.04"
"latest"
"trusty"
"trusty-20171117"
"xenial"
...
``` 
+ 需要先apt安装轻量JSON处理程序 `jq`

### 验证

安装成功后验证如下：

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

[后一篇](04-install_kube_master.md)
