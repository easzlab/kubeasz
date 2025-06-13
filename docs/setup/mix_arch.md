# 混合架构集群部署

混合架构集群本文是指集群中既有linux amd64架构机器，也有linux arm64架构机器；这里只记录一个简单的操作说明，实际操作注意风险。

## 部署思路

1. 先选定一台amd64架构的机器做“amd64部署机”，使用它先部署amd64架构的集群
2. 选一台arm64架构的机器做“arm64部署机”，复制amd64部署机的/etc/kubeasz目录文件（除去目录中的bin、down子目录），然后重新下载arm64架构的二进制和镜像，然后添加arm64节点到原有集群即可

## 操作步骤

1. 假设已经正常部署了amd64架构的三节点集群
2. 在“amd64部署机” 目录 /etc/kubeasz 中移除子目录 bin 和 down，然后把整体/etc/kubeasz 目录复制到“arm64部署机”

```
# 登录amd64部署机
cd /etc/kubeasz; mv bin down /tmp/; scp -r /etc/kubeasz root@{_ip_arm64}:/etc/
# 复制完成后找回 bin 和 down 子目录
mv /tmp/bin /etc/kubeasz/; mv /tmp/down /etc/kubeasz/
```
3. 登录“arm64部署机”，执行下载，其他准备工作

```
cd /etc/kubeasz
# 下载基础部分
./ezdown -D
# 下载额外部分（如有）
./ezdown -X ...
# 运行部署容器
./ezdown -S
# 配置机器ssh免密码登录，集群所有节点都免密，包括待新增arm64节点
ssh-copy-id xx.xx.xx.xx
ssh-copy-id ...
# 复制kubeconfig
mkdir /root/.kube/; cp clusters/default/kubectl.kubeconfig /root/.kube/config
```
4. 添加arm64新节点到集群

```
source ~/.bashrc
# 添加新节点 x.x.x.x
dk ezctl add-node default x.x.x.x
```
5. 验证

```
$ kubectl get node -owide
NAME           STATUS   ROLES    AGE    VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION       CONTAINER-RUNTIME
k8s-x.x.x-19   Ready    master   5d8h   v1.33.1   x.x.x.19   <none>        Ubuntu 20.04.4 LTS   5.4.0-122-generic    containerd://2.1.1
k8s-x.x.x-90   Ready    node     5d8h   v1.33.1   x.x.x.90   <none>        Ubuntu 22.04.5 LTS   5.15.0-134-generic   containerd://2.1.1
k8s-x.x.x-91   Ready    node     5d8h   v1.33.1   x.x.x.91   <none>        Ubuntu 22.04.5 LTS   5.15.0-134-generic   containerd://2.1.1
k8s-x.x.x-93   Ready    node     79s    v1.33.1   x.x.x.93   <none>        Ubuntu 22.04.5 LTS   5.15.0-140-generic   containerd://2.1.1

$ kubectl describe node|grep beta.kubernetes.io/arch
Labels:             beta.kubernetes.io/arch=amd64
Labels:             beta.kubernetes.io/arch=amd64
Labels:             beta.kubernetes.io/arch=amd64
Labels:             beta.kubernetes.io/arch=arm64
```

## 小结
通过以上步骤，成功实现了在amd64集群中添加arm64节点；充分展示kubeasz 项目部署集群的灵活性和可配置性；部署过程中ansible执行的过程性输出内容，以近乎白盒的方式展示每一个细节；假如出错有详细的说明，帮助定位，并且随时可以修改执行脚本，安装的幂等性保证随时可以重新安装以修复错误。`Hack it, and have fun!`

