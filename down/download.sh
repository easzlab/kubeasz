#!/bin/bash
#主要组件版本如下
export K8S_VER=v1.7.7
export ETCD_VER=v3.2.8
export FLANNEL_VER=v0.9.0
export DOCKER_VER=17.09.0-ce
export CALICO_VER=v1.6.1

echo "\n----download k8s binary at:"
echo https://dl.k8s.io/${K8S_VER}/kubernetes-server-linux-amd64.tar.gz

echo "\n----download etcd binary at:"
echo https://github.com/coreos/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz
echo https://storage.googleapis.com/etcd/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz

echo "\n----download flannel binary at:"
echo https://github.com/coreos/flannel/releases/download/${FLANNEL_VER}/flannel-${FLANNEL_VER}-linux-amd64.tar.gz

echo "\n----download calicoctl binary at:"
echo https://github.com/projectcalico/calicoctl/releases/download/${CALICO_VER}/calicoctl

echo "\n----download docker binary at:"
echo https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VER}.tgz

echo "\n----download ca tools at:"
echo https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
echo https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
echo https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64

echo "\n注意1：因为网络原因不进行自动下载"
echo "请按照以上链接手动下载二进制包到down目录中，包含如下："
echo "-rw-r--r-- 1 root root   6595195 Mar 30  2016 cfssl-certinfo_linux-amd64"
echo "-rw-r--r-- 1 root root   2277873 Mar 30  2016 cfssljson_linux-amd64"
echo "-rw-r--r-- 1 root root  10376657 Mar 30  2016 cfssl_linux-amd64"
echo "-rwxr--r-- 1 root root  29699281 Aug 24 17:09 docker-17.06.1-ce.tgz*"
echo "-rwxr--r-- 1 root root  10176124 Aug 24 17:09 etcd-v3.2.6-linux-amd64.tar.gz*"
echo "-rwxr--r-- 1 root root   9090192 Aug 24 17:09 flannel-v0.8.0-linux-amd64.tar.gz*"
echo "-rwxr--r-- 1 root root 437406487 Aug 24 17:09 kubernetes-server-linux-amd64.tar.gz*"

echo "\n注意2：如果还没有手工下载tar包，请Ctrl-c结束此脚本\nsleep 60"
sleep 60

mkdir -p ../bin
### 准备证书工具程序
echo "\n准备证书工具程序..."
if [ -f "cfssl_linux-amd64" ]; then
  mv cfssl_linux-amd64 ../bin/cfssl
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
fi
if [ -f "cfssljson_linux-amd64" ]; then
  mv cfssljson_linux-amd64 ../bin/cfssljson
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
fi
if [ -f "cfssl-certinfo_linux-amd64" ]; then
  mv cfssl-certinfo_linux-amd64 ../bin/cfssl-certinfo
else
  echo 请先下载https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
fi

### 准备etcd程序
echo "\n准备etcd二进制程序..."
if [ -f "etcd-${ETCD_VER}-linux-amd64.tar.gz" ]; then
  echo "\nextracting etcd binaries..."
  tar zxf etcd-${ETCD_VER}-linux-amd64.tar.gz
  mv etcd-${ETCD_VER}-linux-amd64/etcd* ../bin
else
  echo 请先下载etcd-${ETCD_VER}-linux-amd64.tar.gz
fi

### 准备flannel程序
echo "\n准备flannel二进制程序..."
if [ -f "flannel-${FLANNEL_VER}-linux-amd64.tar.gz" ]; then
  echo "\nextracting flannel binaries..."
  tar zxf flannel-${FLANNEL_VER}-linux-amd64.tar.gz
  mv flanneld mk-docker-opts.sh ../bin
  rm README.md
else
  echo 请先下载flannel-${FLANNEL_VER}-linux-amd64.tar.gz
fi

### 准备kubernetes程序
echo "\n准备kubernetes二进制程序..."
if [ -f "kubernetes-server-linux-amd64.tar.gz" ]; then
  echo "\nextracting kubernetes binaries..."
  tar zxf kubernetes-server-linux-amd64.tar.gz
  mv kubernetes/server/bin/kube-apiserver ../bin
  mv kubernetes/server/bin/kube-controller-manager ../bin
  mv kubernetes/server/bin/kubectl ../bin
  mv kubernetes/server/bin/kubelet ../bin
  mv kubernetes/server/bin/kube-proxy ../bin
  mv kubernetes/server/bin/kube-scheduler ../bin
else
  echo 请先下载kubernetes-server-linux-amd64.tar.gz
fi

### 准备docker程序
echo "\n准备docker二进制程序..."
if [ -f "docker-${DOCKER_VER}.tgz" ]; then
  echo "\nextracting docker binaries..."
  tar zxf docker-${DOCKER_VER}.tgz
  mv docker/docker* ../bin
  if [ -f "docker/completion/bash/docker" ]; then
    mv -f docker/completion/bash/docker ../roles/kube-node/files/docker
  fi
else
  echo 请先下载docker-${DOCKER_VER}.tgz
fi
