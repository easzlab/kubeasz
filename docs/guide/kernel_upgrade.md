# Linux Kernel 升级

k8s,docker,cilium等很多功能、特性需要较新的linux内核支持，所以有必要在集群部署前对内核进行升级；CentOS7 和 Ubuntu16.04可以很方便的完成内核升级。

## CentOS7

红帽企业版 Linux 仓库网站 https://www.elrepo.org，主要提供各种硬件驱动（显卡、网卡、声卡等）和内核升级相关资源；兼容 CentOS7 内核升级。如下按照网站提示载入elrepo公钥及最新elrepo版本，然后按步骤升级内核（以安装长期支持版本 kernel-lt 为例）

``` bash
#安装所需软件包
yum install -y perl wget

#下载所需内核版本的 RPM 包，更多版本可以从中寻找（http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/）
wget http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/kernel-lt-5.4.278-1.el7.elrepo.x86_64.rpm
wget http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/kernel-lt-devel-5.4.278-1.el7.elrepo.x86_64.rpm
wget http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/kernel-lt-headers-5.4.278-1.el7.elrepo.x86_64.rpm
wget http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/kernel-lt-tools-5.4.278-1.el7.elrepo.x86_64.rpm
wget http://mirrors.coreix.net/elrepo-archive-archive/kernel/el7/x86_64/RPMS/kernel-lt-tools-libs-5.4.278-1.el7.elrepo.x86_64.rpm

# 卸载旧版工具（安装kernel-lt-tools会和旧版本的kernel-tools导致冲突，需要卸载旧版本的）
yum remove kernel-tools kernel-tools-libs -y

#安装下载的 RPM 包
rpm -ivh kernel-lt-tools-libs-5.4.278-1.el7.elrepo.x86_64.rpm
rpm -ivh kernel-lt-tools-5.4.278-1.el7.elrepo.x86_64.rpm 
rpm -ivh kernel-lt-5.4.278-1.el7.elrepo.x86_64.rpm
rpm -ivh kernel-lt-devel-5.4.278-1.el7.elrepo.x86_64.rpm 

#验证安装，可以看到新版本的和旧版本的
rpm -qa | grep kernel
kernel-lt-5.4.278-1.el7.elrepo.x86_64
kernel-lt-tools-libs-5.4.278-1.el7.elrepo.x86_64
kernel-3.10.0-1160.71.1.el7.x86_64
kernel-lt-devel-5.4.278-1.el7.elrepo.x86_64
kernel-lt-tools-5.4.278-1.el7.elrepo.x86_64

#查看默认启动顺序
awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
0 : CentOS Linux (5.4.278-1.el7.elrepo.x86_64) 7 (Core)
1 : CentOS Linux (3.10.0-1160.71.1.el7.x86_64) 7 (Core)
2 : CentOS Linux (0-rescue-0b208d4cc51848998d32430e022d3040) 7 (Core)
#设置默认启动内核顺序
grub2-set-default 0  
#重启
reboot
#重启后进行检查是否成功切换到新内核
uname -r
5.4.278-1.el7.elrepo.x86_64
```

## Ubuntu16.04

``` bash
打开 http://kernel.ubuntu.com/~kernel-ppa/mainline/ 并选择列表中选择你需要的版本（以4.16.3为例）。
接下来，根据你的系统架构下载 如下.deb 文件：
Build for amd64 succeeded (see BUILD.LOG.amd64):
  linux-headers-4.16.3-041603_4.16.3-041603.201804190730_all.deb
  linux-headers-4.16.3-041603-generic_4.16.3-041603.201804190730_amd64.deb
  linux-image-4.16.3-041603-generic_4.16.3-041603.201804190730_amd64.deb
#安装后重启即可
$ sudo dpkg -i *.deb
```
