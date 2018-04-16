#!/bin/bash

# 本脚本提供如下功能，作者：Joey Yang, https://github.com/Code2Life
# 1. 在Ubuntu/CentOS/Fedora/ArchLinux中自动化的安装python+ansible;
# 2. clone kubeasz项目代码, 并将需要的二进制文件下载解压到/etc/ansible/bin中;
# 另外, 相关的k8s二进制文件, 我同步到了个人在七牛上的CDN存储中（速度更快）, 方便大家下载: filecdn.code2life.top;
# 
# 使用方法：
# 1. 支持带参数的运行, 如: ./basic-env-setup.sh k8s.193.tar.gz 指定不同的kubernetes二进制文件, 无参数时默认最新的k8s.1100.tar.gz (k8s 1.10.0 + etcd 3.3.2).
# 2. 也可以在任何一台支持的linux设备运行：curl http://filecdn.code2life.top/kubeasz-basic-env-setup.sh | sh -s
# 已经亲测centos7/ubuntu16.04/debian9/fedora27都是可以的, 二进制包下载速度贼快.脚本运行完毕之后, 只需到/etc/ansible目录下配置好hosts, 复制完ssh的公钥即可通过ansible-playbook迅速搭建集群了.

set -e

# curl http://filecdn.code2life.top/kubeasz-basic-env-setup.sh | sh -s

# 默认1.10.0 版本的 Kubernetes
bin_resource_url='http://filecdn.code2life.top/k8s.1100.tar.gz'

# 如果参数指定k8s相关的bin以指定的为准, 例如: k8s.193.tar.gz
if [ "$1" ];then
  bin_resource_url="http://filecdn.code2life.top/"$1
fi

# 各Linux版本安装python/pip
# ---------------------------

# debian 默认的apt源在国内访问很慢, 可手动修改/etc/apt/source.list修改为其他源
# 以 debian 9 为例, source.list可修改为如下内容, ubuntu修改方法类似, 找到相应系统和版本的镜像源替换即可
# deb http://mirrors.163.com/debian/  stretch main non-free contrib
# deb http://mirrors.163.com/debian/  stretch-updates main non-free contrib
# deb http://mirrors.163.com/debian/  stretch-backports main non-free contrib
# deb http://mirrors.163.com/debian-security/  stretch/updates main non-free contrib
basic_ubuntu_debian() {
  echo "Setup Basic Environment for Ubuntu/Debian."
  apt-get update && apt-get upgrade -y && apt-get dist-upgrade -y
  apt-get install python2.7 git python-pip curl -y

  if [ ! -f /usr/bin/python ];then
    ln -s /usr/bin/python2.7 /usr/bin/python
  fi
}

# 红帽系Liunx可修改yum源加快下载速度, 修改/etc/yum.repos.d内文件即可
basic_centos() {
  echo "Setup Basic Environment for CentOS."
  yum install epel-release -y
  yum update -y
  yum erase firewalld firewalld-filesystem python-firewall -y
  yum install git python python-pip curl -y
}

basic_fedora() {
  echo "Setup Basic Environment for Fedora."
  yum update -y
  yum install git python python-pip curl -y
}

# archlinux 使用pacman进行包管理
basic_arch() {
  pacman -Syu --noconfirm
  pacman -S python git python-pip curl --noconfirm
}

# 使用pip安装ansible, 并下载k8s相关bin文件
setup_ansible_k8s() {
  echo "Download Ansible and Kubernetes binaries."
  pip install pip --upgrade -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
  pip install --no-cache-dir ansible -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com

  git clone https://github.com/gjmzj/kubeasz.git
  mv kubeasz /etc/ansible

  # Download from CDN & Move bin files
  curl -o k8s_download.tar.gz "$bin_resource_url"
  tar zxvf k8s_download.tar.gz
  mv -f bin/* /etc/ansible/bin
  rm -rf bin
  echo "Finish setup. Please config your hosts and run 'ansible-playbook' command at /etc/ansible."
}
# ---------------------------

# 判断Linux发行版, 执行不同基础环境设置方法
# ---------------------------
lsb_dist=''
command_exists() {
    command -v "$@" > /dev/null 2>&1
}
if command_exists lsb_release; then
    lsb_dist="$(lsb_release -si)"
    lsb_version="$(lsb_release -rs)"
fi
if [ -z "$lsb_dist" ] && [ -r /etc/lsb-release ]; then
    lsb_dist="$(. /etc/lsb-release && echo "$DISTRIB_ID")"
    lsb_version="$(. /etc/lsb-release && echo "$DISTRIB_RELEASE")"
fi
if [ -z "$lsb_dist" ] && [ -r /etc/debian_version ]; then
    lsb_dist='debian'
fi
if [ -z "$lsb_dist" ] && [ -r /etc/fedora-release ]; then
    lsb_dist='fedora'
fi
if [ -z "$lsb_dist" ] && [ -r /etc/os-release ]; then
    lsb_dist="$(. /etc/os-release && echo "$ID")"
fi
if [ -z "$lsb_dist" ] && [ -r /etc/centos-release ]; then
    lsb_dist="$(cat /etc/*-release | head -n1 | cut -d " " -f1)"
fi
if [ -z "$lsb_dist" ] && [ -r /etc/redhat-release ]; then
    lsb_dist="$(cat /etc/*-release | head -n1 | cut -d " " -f1)"
fi
lsb_dist="$(echo $lsb_dist | cut -d " " -f1)"
lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
# ---------------------------

# ---------------------------
setup_env(){
    case "$lsb_dist" in
        centos)
        basic_centos
        setup_ansible_k8s
        exit 0
    ;;
        fedora)
        basic_fedora
        setup_ansible_k8s
        exit 0
    ;;
        ubuntu)
        basic_ubuntu_debian
        setup_ansible_k8s
        exit 0
    ;;
        debian)
        basic_ubuntu_debian
        setup_ansible_k8s
        exit 0
    ;;
        arch)
        basic_arch
        setup_ansible_k8s
        exit 0
    ;;
        suse)
        echo 'Not implementation yet.'
        exit 1
    esac
    echo "Error: Unsupported OS, please set ansible environment manually."
    exit 1
}
setup_env
# ---------------------------
