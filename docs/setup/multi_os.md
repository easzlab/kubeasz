# 操作系统说明

目前发现部分使用新内核的linux发行版，k8s 安装使用 cgroup v2版本时，有时候安装会失败，需要删除/清理集群后重新安装。已报告可能发生于 Alma Linux 9, Rocky Linux 9, Fedora 37；建议如下步骤处理：

- 1.确认系统使用的cgroup v2版本
```
stat -fc %T /sys/fs/cgroup/ 
cgroup2fs
```
- 2.初次安装时kubelet可能启动失败，日志报错类似：err="openat2 /sys/fs/cgroup/kubepods.slice/cpu.weight: no such file or directory"

- 3.建议删除集群然后重新安装，一般能够成功
```
# 删除集群
dk ezctl destroy xxxx

# 重启
reboot

# 启动后重新安装
dk ezctl setup xxxx all
```

## Debian

- Debian 11：默认可能没有安装iptables，使用kubeasz 安装前需要执行：

``` bash 
apt update

apt install iptables -y
```

## Alibaba

- Alibaba Linux 3.2104 LTS：安装前需要设置如下：

``` bash
# 修改使用dnf包管理
sed -i 's/package/dnf/g' /etc/kubeasz/roles/prepare/tasks/redhat.yml
```

## openSUSE

- openSUSE Leap 15.4：需要安装iptables

``` bash
zypper install iptables
ln -s /usr/sbin/iptables /sbin/iptables
```
