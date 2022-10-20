# 关于K8S集群一致性认证

CNCF 一致性认证项目(https://github.com/cncf/k8s-conformance) 可以很方便帮助k8s搭建者和用户确认集群各项功能符合预期，既符合k8s设计标准。

# kubeasz 通过一致性测试

Cheers! 

自kubeasz 3.0.0 版本，k8s v1.20.2开始，正式通过cncf一致性认证，成为cncf 官方认证安装工具；后续k8s主要版本发布或者kubeasz有大版本更新，会优先确保通过集群一致性认证。

- v1.25 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.25/kubeasz)
- v1.24 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.24/kubeasz)
- v1.23 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.23/kubeasz)
- v1.22 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.22/kubeasz)
- v1.21 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.21/kubeasz)
- v1.20 [已认证](https://github.com/cncf/k8s-conformance/tree/master/v1.20/kubeasz)


## Conformance Test

按照测试文档，注意以下几点：

1.解决qiang的问题，可以临时去国外公有云创建集群，然后运行测试项目。

2.集群要保障资源，建议3个节点

3.网络组件选择calico，其他组件可能有bug导致特定测试项失败


# 附：测试流程

## Node Provisioning

Provision 3 nodes for your cluster (OS: Ubuntu 20.04)

1 master node (4c16g)

2 worker node (4c16g)

for a High-Availability Kubernetes Cluster, read [more](https://github.com/easzlab/kubeasz/blob/master/docs/setup/00-planning_and_overall_intro.md)

## Install the cluster

(1) Download 'kubeasz' code, the binaries and offline images

```
export release=3.2.0
curl -C- -fLO --retry 3 https://github.com/easzlab/kubeasz/releases/download/${release}/ezdown
chmod +x ./ezdown
./ezdown -D -m standard
```

(2) install an all-in-one cluster

```
./ezdown -S
docker exec -it kubeasz ezctl start-aio
```

(3) Add two worker nodes

```
ssh-copy-id ${worker1_ip}
ssh ${worker1_ip} ln -s /usr/bin/python3 /usr/bin/python
docker exec -it kubeasz ezctl add-node default ${worker1_ip}
ssh-copy-id ${worker2_ip}
ssh ${worker2_ip} ln -s /usr/bin/python3 /usr/bin/python
docker exec -it kubeasz ezctl add-node default ${worker2_ip}
```

## Run Conformance Test

The standard tool for running these tests is
[Sonobuoy](https://github.com/heptio/sonobuoy).  Sonobuoy is
regularly built and kept up to date to execute against all
currently supported versions of kubernetes.

Download a [binary release](https://github.com/heptio/sonobuoy/releases) of the CLI, or build it yourself by running:

```
$ go get -u -v github.com/heptio/sonobuoy
```

Deploy a Sonobuoy pod to your cluster with:

```
$ sonobuoy run --mode=certified-conformance
```

**NOTE:** You can run the command synchronously by adding the flag `--wait` but be aware that running the Conformance tests can take an hour or more.

View actively running pods:

```
$ sonobuoy status
```

To inspect the logs:

```
$ sonobuoy logs
```

Once `sonobuoy status` shows the run as `completed`, copy the output directory from the main Sonobuoy pod to a local directory:

```
$ outfile=$(sonobuoy retrieve)
```

This copies a single `.tar.gz` snapshot from the Sonobuoy pod into your local
`.` directory. Extract the contents into `./results` with:

```
mkdir ./results; tar xzf $outfile -C ./results
```

**NOTE:** The two files required for submission are located in the tarball under **plugins/e2e/results/{e2e.log,junit.xml}**.

To clean up Kubernetes objects created by Sonobuoy, run:

```
sonobuoy delete
```
