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
- [可选] 修改 roles/chrony/var/main.yml 中的变量定义，关于文件 roles/chrony/var/main.yml 的由来请看[这里](../config_guide.md)

对于新集群或者新节点，`chrony` 的安装配置已经集成到 `90.setup.yml` `01.prepare.yml` `20.addnode.yml` `21.addmaster.yml` 等脚本中；对于已运行中的集群请执行如下命令进行安装：

`ansible-playbook /etc/ansible/roles/chrony/chrony.yml `

