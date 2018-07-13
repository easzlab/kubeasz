#!/bin/bash
#主要组件版本如下
export K8S_VER=v1.10.2
export ETCD_VER=v3.3.8
export DOCKER_VER=18.03.1-ce
export CNI_VER=v0.7.0
export DOCKER_COMPOSE=1.18.0
export HARBOR=v1.2.2

echo "\n建议直接下载本人打包好的所有必要二进制包k8s-***.all.tar.gz，然后解压到bin目录"
echo "\n建议不使用此脚本，如果你想升级组件或者实验，请通读该脚本，必要时适当修改后使用"
echo "\n注意1：请按照以下链接手动下载二进制包到down目录中"
echo "\n注意2：如果还没有手工下载tar包，请Ctrl-c结束此脚本"

echo "\n----download k8s binary at:"
echo https://dl.k8s.io/${K8S_VER}/kubernetes-server-linux-amd64.tar.gz

echo "\n----download etcd binary at:"
echo https://github.com/coreos/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz
echo https://storage.googleapis.com/etcd/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz

echo "\n----download docker binary at:"
echo https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VER}.tgz

echo "\n----download ca tools at:"
echo https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
echo https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
echo https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64

echo "\n----download docker-compose at:"
echo https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE}/docker-compose-Linux-x86_64

echo "\n----download harbor-offline-installer at:"
echo https://github.com/vmware/harbor/releases/download/${HARBOR}/harbor-offline-installer-${HARBOR}.tgz

echo "\n----download cni plugins at:"
echo https://github.com/containernetworking/plugins/releases

sleep 30

### 准备证书工具程序
echo "\n准备证书工具程序..."
if [ -f "cfssl_linux-amd64" ]; then
  mv -f cfssl_linux-amd64 ../bin/cfssl
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
fi
if [ -f "cfssljson_linux-amd64" ]; then
  mv -f cfssljson_linux-amd64 ../bin/cfssljson
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
fi
if [ -f "cfssl-certinfo_linux-amd64" ]; then
  mv -f cfssl-certinfo_linux-amd64 ../bin/cfssl-certinfo
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
fi

### 准备etcd程序
echo "\n准备etcd二进制程序..."
if [ -f "etcd-${ETCD_VER}-linux-amd64.tar.gz" ]; then
  echo "\nextracting etcd binaries..."
  tar zxf etcd-${ETCD_VER}-linux-amd64.tar.gz
  mv -f etcd-${ETCD_VER}-linux-amd64/etcd* ../bin
else
  echo 请先下载etcd-${ETCD_VER}-linux-amd64.tar.gz
fi

### 准备kubernetes程序
echo "\n准备kubernetes二进制程序..."
if [ -f "kubernetes-server-linux-amd64.tar.gz" ]; then
  echo "\nextracting kubernetes binaries..."
  tar zxf kubernetes-server-linux-amd64.tar.gz
  mv -f kubernetes/server/bin/kube-apiserver ../bin
  mv -f kubernetes/server/bin/kube-controller-manager ../bin
  mv -f kubernetes/server/bin/kubectl ../bin
  mv -f kubernetes/server/bin/kubelet ../bin
  mv -f kubernetes/server/bin/kube-proxy ../bin
  mv -f kubernetes/server/bin/kube-scheduler ../bin
else
  echo 请先下载kubernetes-server-linux-amd64.tar.gz
fi

### 准备docker程序
echo "\n准备docker二进制程序..."
if [ -f "docker-${DOCKER_VER}.tgz" ]; then
  echo "\nextracting docker binaries..."
  tar zxf docker-${DOCKER_VER}.tgz
  mv -f docker/docker* ../bin
  if [ -f "docker/completion/bash/docker" ]; then
    mv -f docker/completion/bash/docker ../roles/docker/files/docker
  fi
else
  echo 请先下载docker-${DOCKER_VER}.tgz
fi

### 准备cni plugins，仅安装flannel需要，安装calico由容器专门下载cni plugins 
echo "\n准备cni plugins，仅安装flannel需要，安装calico由容器专门下载cni plugins..."
if [ -f "cni-${CNI_VER}.tgz" ]; then
  echo "\nextracting cni plugins binaries..."
  tar zxf cni-${CNI_VER}.tgz
  mv -f bridge ../bin
  mv -f flannel ../bin
  mv -f host-local ../bin
  mv -f loopback ../bin
  mv -f portmap ../bin
else
  echo 请先下载cni-${CNI_VER}.tgz
fi
