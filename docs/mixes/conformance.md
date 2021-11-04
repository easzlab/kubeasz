## 关于K8S集群一致性认证

CNCF 一致性认证项目(https://github.com/cncf/k8s-conformance) 可以很方便帮助k8s搭建者和用户确认集群各项功能符合预期，既符合k8s设计标准。

## Conformance Test

按照测试文档，注意以下几点，通过所有的测试项也不是难事：

1.解决qiang的问题，可以临时去国外公有云创建集群，然后运行测试项目。

2.集群要保障资源，建议3个节点左右

3.网络组件选择calico，其他组件可能有bug导致特定测试项失败

4.kube-proxy暂时用iptables模式，使用ipvs再测试服务sessionAffinity时有bug，后续应该会修复


## kubeasz 技术上完全通过一致性测试

Cheers! 

使用kubeasz 3.0.0 版本，k8s v1.20.2（其他kubeasz版本应该也类似），开始测试时候在网络上走了一些弯路，后面还是很顺利的通过测试，测试结果：

``` bash
JUnit report was created: /tmp/results/junit_01.xml
{"msg":"Test Suite completed","total":311,"completed":311,"skipped":5356,"failed":0}

Ran 311 of 5667 Specs in 6179.487 seconds
SUCCESS! -- 311 Passed | 0 Failed | 0 Pending | 5356 Skipped
PASS

Ginkgo ran 1 suite in 1h43m0.59512776s
Test Suite Passed
```

具体的测试过程和结果请参考这里：https://github.com/cncf/k8s-conformance/pull/1326

PS：另外，我也花时间走流程正式申请成为官方认证的部署工具；目前来看作为免费的开源工具申请下来还是比较困难，估计是类似的发行版及部署工具太多了吧，中文项目估计也不被看好，有兴趣的或者有门路的朋友可以联系我，帮忙申请下来。

后续k8s主要版本发布或者kubeasz有大版本更新，我都会优先确保通过集群一致性认证。


## 附：测试流程

### Node Provisioning

Provision 2 nodes for your cluster (OS requirements: CentOS 7 or Ubuntu 1604/1804)

1 master node (4c16g)

1 worker node (4c16g)

for a High-Availability Kubernetes Cluster, read [more](https://github.com/easzlab/kubeasz/blob/master/docs/setup/00-planning_and_overall_intro.md)

### Install the cluster

(1) clone repo: kubeasz

```
git clone https://github.com/easzlab/kubeasz.git
mv ./kubeasz /etc
```

(2) Download the binaries and offline images

```
cd /etc/kubeasz
./ezdown -D -m standard
```

(3) install an all-in-one cluster

```
sed -i 's/^CLUSTER_NETWORK=.*$/CLUSTER_NETWORK="calico"/g' example/hosts.allinone
sed -i 's/^PROXY_MODE=.*$/PROXY_MODE="iptables"/g' example/hosts.allinone
./ezdown -S
docker exec -it kubeasz ezctl start-aio
```

(4) Add a worker node

```
ssh-copy-id ${worker_ip}
docker exec -it kubeasz ezctl add-node default ${worker_ip}
```

### Run Conformance Test
The standard tool for running these tests is Sonobuoy. Sonobuoy is regularly built and kept up to date to execute against all currently supported versions of kubernetes.

Download a [binary release](https://github.com/vmware-tanzu/sonobuoy/releases) of the CLI, or build it yourself by running:

```
go get -u -v github.com/vmware-tanzu/sonobuoy
```

Deploy a Sonobuoy pod to your cluster with:

```
sonobuoy run --mode=certified-conformance
```

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
