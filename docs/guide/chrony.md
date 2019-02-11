# chrony 时间同步

在安装k8s集群前需确保各节点时间同步；`chrony` 是一个优秀的 `NTP` 实现，性能比ntp好，且配置管理方便；它既可作时间服务器服务端，也可作客户端。

- `OpenStack` 社区也推荐使用 `chrony`实现各节点之间的时间同步

## 安装配置介绍

项目中选定一个节点(`deploy` )作为集群内部其他节点的时间同步源，而 deploy节点本身从公网源同步；当然如果整个集群都无法访问公网，那么请手动校准deploy 节点的时间后，仍旧可以作为内部集群的时间源服务器。

- 配置 chrony server, 在`/etc/chrony.conf` 配置以下几项，其他项默认值即可

``` bash
# 1. 配置时间源，国内可以增加阿里的时间源 ntp1.aliyun.com
server {{ ntp_server }} iburst

# 2. 配置允许同步的客户端网段
allow {{ local_network }}

# 3. 配置离线也能作为源服务器
local stratum 10
```

- 配置 chrony client

``` bash
# 1. 清除所有其他时间源，只配置一个本地 deploy节点作为源
server {{ groups.deploy[0] }} iburst

# 2. 其他所有项可以默认配置
```

## `kubeasz` 集成安装

- 修改 ansible hosts 文件，在 `deploy` 节点配置 `NTP_ENABLED=yes` (默认: no)
- [可选] 修改 roles/chrony/var/main.yml 中的变量定义，关于文件 roles/chrony/var/main.yml 的由来请看[这里](../setup/config_guide.md)

对于新集群或者新节点，`chrony` 的安装配置已经集成到 `90.setup.yml` `01.prepare.yml` `20.addnode.yml` `21.addmaster.yml` 等脚本中；对于已运行中的集群请执行如下命令进行安装：

`ansible-playbook /etc/ansible/roles/chrony/chrony.yml `

## 验证配置

- 在 chrony server (deploy 节点)检查时间源信息，默认配置为`ntp1.aliyun.com`的地址：

```
$ chronyc sources -v
210 Number of sources = 5

  .-- Source mode  '^' = server, '=' = peer, '#' = local clock.
 / .- Source state '*' = current synced, '+' = combined , '-' = not combined,
| /   '?' = unreachable, 'x' = time may be in error, '~' = time too variable.
||                                                 .- xxxx [ yyyy ] +/- zzzz
||      Reachability register (octal) -.           |  xxxx = adjusted offset,
||      Log2(Polling interval) --.      |          |  yyyy = measured offset,
||                                \     |          |  zzzz = estimated error.
||                                 |    |           \
MS Name/IP address         Stratum Poll Reach LastRx Last sample
===============================================================================
^* 120.25.115.20                 2   9   377    55   +147us[ +250us] +/-   15ms
^- 85.199.214.100                1  10   377   182    -25ms[  -24ms] +/-  128ms
^- makaki.miuku.net              2  10   367   307    +61ms[  +61ms] +/-  127ms
^- static-5-103-139-163.ip.f     1   9   167   572   +532us[ +336us] +/-  117ms
^- 119.28.183.184                2   7   377    33   -130us[ -130us] +/-   47ms
```

- 在 chrony server (deploy 节点)检查时间源同步状态

```
chronyc sourcestats -v
210 Number of sources = 5
                             .- Number of sample points in measurement set.
                            /    .- Number of residual runs with same sign.
                           |    /    .- Length of measurement set (time).
                           |   |    /      .- Est. clock freq error (ppm).
                           |   |   |      /           .- Est. error in freq.
                           |   |   |     |           /         .- Est. offset.
                           |   |   |     |          |          |   On the -.
                           |   |   |     |          |          |   samples. \
                           |   |   |     |          |          |             |
Name/IP Address            NP  NR  Span  Frequency  Freq Skew  Offset  Std Dev
==============================================================================
120.25.115.20              15  11   44m     +0.011      0.909  +4097ns   758us
85.199.214.100             22  13   49m     -3.588      5.097    -23ms  5709us
makaki.miuku.net           22  14   46m     +2.455      6.225    +64ms  4945us
static-5-103-139-163.ip.f  20  13   42m     -2.472     10.168  +3615us  6732us
119.28.183.184             16   9   19m    +10.378     25.190  +3469us  6803us
```

- 在 chrony client 检查，可以看到时间源只有一个（deploy 节点地址）

```
$ chronyc sources
210 Number of sources = 1
MS Name/IP address         Stratum Poll Reach LastRx Last sample
===============================================================================
^* 192.168.1.1                  3   6   377    15  +4085ns[  -25us] +/-   15ms
$ chronyc sourcestats
210 Number of sources = 1
Name/IP Address            NP  NR  Span  Frequency  Freq Skew  Offset  Std Dev
==============================================================================
192.168.1.1                5   4   323     -0.252      0.819  -3031ns    15us
```

## 验证时间同步状态完成

chrony 服务启动后，chrony server 会与配置的公网参考时间服务器进行同步；server 同步完成后，chrony client 会与 server 进行时间同步；一般来说整个集群达到时间同步需要几十分钟。可以用如下命令检查，初始时 **NTP synchronized: no**，同步完成后 **NTP synchronized: yes**

``` bash
$ ansible all -m shell -a 'timedatectl'
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
