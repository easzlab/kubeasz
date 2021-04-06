# chrony 时间同步

在安装k8s集群前需确保各节点时间同步；`chrony` 是一个优秀的 `NTP` 实现，性能比ntp好，且配置管理方便；它既可作时间服务器服务端，也可作客户端。

- `OpenStack` 社区也推荐使用 `chrony`实现各节点之间的时间同步

## 安装配置介绍

项目中选定一个节点(`groups.chrony[0]`)作为集群内部其他节点的时间同步源，而这个节点本身从公网源同步；当然如果整个集群都无法访问公网，那么请手动校准这个节点的时间后，仍旧可以作为内部集群的时间源服务器。

- 配置 chrony server，详见roles/chrony/templates/server.conf.j2 

- 配置 chrony client，详见roles/chrony/templates/client.conf.j2

## `kubeasz` 集成安装

- 修改 clusters/${cluster_name}/hosts 文件，在 `chrony`组中加入选中的节点ip
- [可选] 修改 clusters/${cluster_name}/config.yml 中的相关配置
-执行命令安装 `ezctl setup ${cluster_name} 01`

## 验证安装

- 检查chronyd服务状态 `systemctl status chronyd`
- 检查chronyd时间同步日志 `/var/log/chrony`

## 验证时间同步状态完成

chrony 服务启动后，chrony server 会与配置的公网参考时间服务器进行同步；server 同步完成后，chrony client 会与 server 进行时间同步；一般来说整个集群达到时间同步需要几十分钟。可以用如下命令检查，初始时 **NTP synchronized: no**，同步完成后 **NTP synchronized: yes**

``` bash
$ ansible -i clusters/${cluster_name}/hosts all -m shell -a 'timedatectl'
192.168.1.1 | SUCCESS | rc=0 >>
      Local time: Sat 2019-01-26 11:51:51 HKT
  Universal time: Sat 2019-01-26 03:51:51 UTC
        RTC time: Sat 2019-01-26 03:51:52
       Time zone: Asia/Hong_Kong (HKT, +0800)
 Network time on: yes
NTP synchronized: yes
 RTC in local TZ: no

192.168.1.4 | SUCCESS | rc=0 >>
      Local time: Sat 2019-01-26 11:51:51 HKT
  Universal time: Sat 2019-01-26 03:51:51 UTC
        RTC time: Sat 2019-01-26 03:51:52
       Time zone: Asia/Hong_Kong (HKT, +0800)
 Network time on: yes
NTP synchronized: yes
 RTC in local TZ: no

192.168.1.2 | SUCCESS | rc=0 >>
      Local time: Sat 2019-01-26 11:51:51 HKT
  Universal time: Sat 2019-01-26 03:51:51 UTC
        RTC time: Sat 2019-01-26 03:51:52
       Time zone: Asia/Hong_Kong (HKT, +0800)
 Network time on: yes
NTP synchronized: yes
 RTC in local TZ: no

192.168.1.3 | SUCCESS | rc=0 >>
      Local time: Sat 2019-01-26 11:51:51 HKT
  Universal time: Sat 2019-01-26 03:51:51 UTC
        RTC time: Sat 2019-01-26 03:51:52
       Time zone: Asia/Hong_Kong (HKT, +0800)
 Network time on: yes
NTP synchronized: yes
 RTC in local TZ: no
```
