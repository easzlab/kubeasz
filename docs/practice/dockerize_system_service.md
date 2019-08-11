# 容器化系统服务

## 容器化 haproxy

本例使用 [docker hub 官方](https://github.com/docker-library/haproxy) 维护的 haproxy 镜像；haproxy 配置举例如下

```
global
        log stdout format raw local1 notice
        nbproc 1

defaults
        log     global
        timeout connect 5s
        timeout client  10m
        timeout server  10m

listen apiservers
        bind 0.0.0.0:6443
        mode tcp
        option tcplog
        option dontlognull
        option dontlog-normal
        balance roundrobin 
        server 192.168.1.1 192.168.1.1:6443 check inter 10s fall 2 rise 2 weight 1
        server 192.168.1.2 192.168.1.2:6443 check inter 10s fall 2 rise 2 weight 1
```

在 systemd 系统上编写服务文件如下 /etc/systemd/system/haproxy.service

```
[Unit]
Description=haproxy
Documentation=https://github.com/docker-library/haproxy
After=docker.service
Requires=docker.service

[Service]
User=root
ExecStart=/bin/docker run \
  --name haproxy \
  --publish 6443:6443 \
  --volume /etc/haproxy/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg \
  docker.io/library/haproxy:1.9.8-alpine
ExecStop=/bin/docker rm -f haproxy
ExecReload=/bin/docker kill -s HUP haproxy
Restart=always
RestartSec=10
Delegate=yes
LimitNOFILE=50000
LimitNPROC=50000

[Install]
WantedBy=multi-user.target
```

## 容器化 chrony

- chrony 服务器端配置（假设chrony服务器端192.168.1.1）

```
$ cat /etc/chrony.conf
# Use public servers from the pool.ntp.org project.
server ntp1.aliyun.com iburst
server ntp2.aliyun.com iburst
pool pool.ntp.org iburst

# Ignor source level
stratumweight 0

# Record the rate at which the system clock gains/losses time.
driftfile /var/lib/chrony/drift

# Allow the system clock to be stepped in the first five updates
# if its offset is larger than 1 second.
makestep 1 5

# Enable kernel synchronization of the real-time clock (RTC).
rtcsync

# Allow NTP client access from local network.
allow 0.0.0.0/0

# Serve time even if not synchronized to a time source.
local stratum 10

# Select which information is logged.
#log measurements statistics tracking

#
noclientlog
```
- chrony 客户端配置

```
$ cat /etc/chrony.conf
# Use local chrony server.
server 192.168.1.1 iburst

# Record the rate at which the system clock gains/losses time.
driftfile /var/lib/chrony/drift

# Allow the system clock to be stepped in the first five updates
# if its offset is larger than 1 second.
makestep 1 5

# Enable kernel synchronization of the real-time clock (RTC).
rtcsync

# Select which information is logged.
#log measurements statistics tracking
```

- 在 systemd 系统上编写服务文件如下 /etc/systemd/system/chrony.service

```
[Unit]
Description=chrony
Documentation=https://github.com/kubeasz/dockerfiles/chrony
After=docker.service
Requires=docker.service

[Service]
User=root
ExecStart=/opt/kube/bin/docker run \
  --cap-add SYS_TIME \
  --name chrony \
  --network host \
  --volume /etc/chrony.conf:/etc/chrony/chrony.conf \
  --volume /var/lib/chrony:/var/lib/chrony \
  easzlab/chrony:0.1.0
ExecStartPost=/sbin/iptables -t raw -A PREROUTING -p udp -m udp --dport 123 -j NOTRACK
ExecStartPost=/sbin/iptables -t raw -A OUTPUT -p udp -m udp --sport 123 -j NOTRACK
ExecStop=/opt/kube/bin/docker rm -f chrony
Restart=always
RestartSec=10
Delegate=yes

[Install]
WantedBy=multi-user.target
```
