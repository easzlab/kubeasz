# 如何删除单个节点

本文档所指删除的节点是指使用kubeasz项目安装的节点角色（可能是kube-master, kube-node, etcd, lb节点）

- 警告：此操作将清理单个node节点，包含k8s集群可能使用的数据，特别的：如果有pod使用了本地存储类型，请自行判断重要性

## 删除流程解释

- 1.待删除节点可能是kube-node节点，因此先执行`kubectl drain`，如果不是忽略执行报错
- 2.参照`99.clean.yml`脚本方式删除节点可能的服务和配置，忽略执行报错
- 3.待删除节点可能是kube-node节点，执行`kubectl delete node`, 如果不是忽略执行报错

## 删除操作

- 1.替换待删除节点变量，假设为192.168.1.1
``` bash
$ sed -i 's/NODE_TO_DEL/192.168.1.1/g' /etc/ansible/tools/clean_one_node.yml
```

- 2.执行删除
``` 
$ ansible-playbook /etc/ansible/tools/clean_one_node.yml
```

## Debug

如果出现清理失败，类似报错：`... Device or resource busy: '/var/run/docker/netns/xxxxxxxxxx'`，需要手动umount该目录后重新清理  

``` bash
$ umount /var/run/docker/netns/xxxxxxxxxx
$ ansible-playbook /etc/ansible/tools/clean_one_node.yml
```
