# 删除节点

本文档所指删除的节点是指使用kubeasz项目安装的节点角色（可能是kube-master, kube-node, etcd, lb节点）

- 警告：此操作将清理单个node节点，包含k8s集群可能使用的数据，特别的：如果有pod使用了本地存储类型，请自行判断重要性

## 删除流程解释

- 0.判断待删除节点不是 etcd/master 组的唯一节点，否则不允许删除
- 1.待删除节点可能是kube-node节点，因此先执行`kubectl drain`，如果不是忽略执行报错
- 2.参照`99.clean.yml`脚本方式删除节点可能的服务和配置，忽略执行报错
- 3.待删除节点可能是kube-node节点，执行`kubectl delete node`, 如果不是忽略执行报错
- 4.修改ansible hosts，移除删除节点

## 删除操作

可以使用以下三种方式删除节点（i.e. 192.168.1.1）

``` bash
# 1.推荐使用 easzctl 工具
$ easzctl clean-node 192.168.1.1

# 2.ansible-playbook 带参数执行如下
$ ansible-playbook /etc/ansible/tools/clean_one_node.yml -e NODE_TO_DEL=192.168.1.1

# 3.ansible-playbook 不带参数执行，然后根据提示输入/确认
$ ansible-playbook /etc/ansible/tools/clean_one_node.yml
```

## 验证

- 验证删除节点上是否相关服务均已停止
- 验证 ansible hosts 文件中已删除节点

## Debug

如果出现清理失败，类似报错：`... Device or resource busy: '/var/run/docker/netns/xxxxxxxxxx'`，需要手动umount该目录后重新清理  

``` bash
$ umount /var/run/docker/netns/xxxxxxxxxx
$ ansible-playbook /etc/ansible/tools/clean_one_node.yml
```
