# 替换k8s集群的网络插件

有时候我们在测试环境的k8s集群中希望试用多种网络插件（calico/flannel/kube-router），又不希望每测试一次就全部清除集群然后重建，那么你就需要这个文档。  
- WARNNING：重新安装k8s网络插件会短暂中断已有运行在k8s上的服务
  - 请在熟悉kubeasz的安装流程和k8s网络插件安装流程的基础上谨慎操作
  - 如果k8s集群已经运行庞大业务pod，重装网络插件时会引起所有pod的删除、重建，短时间内将给apiserver带来压力，可能引起master节点夯住

## 替换流程

kubeasz使用标准cni方式安装k8s集群的网络插件；cni负载创建容器网卡和IP分配（IPAM），不同的网络插件（calico,flannel等）创建容器网卡和IP分配方式不一样，所以在替换网络插件时候需要现有pod全部删除，然后自动按照新网络插件的方式重建pod网络；请参考[k8s网络插件章节](../06-安装网络组件.md)。

- 1.清除现有集群网络插件  
``` bash
ansible-playbook /etc/ansible/tools/clean_k8s_network.yml
```

对照脚本`clean_k8s_network.yml` 大致流程为：  
  - 根据实际运行情况，删除现有网络组件的daemonset pod
  - 如果现有组件是kube-router 需要进行一些额外清理和可能需要恢复默认kube-proxy服务
  - 清理cni网络配置和具体插件的运行、配置文件
  - 清理生成的容器网络组件（bridge,tunl等）
  - 如果现有组件是calico 需要额外清理bgp路由
  - 最后删除所有k8s上已运行的pod（会由controller负责重建）

- 2.修改ansible hosts文件指定新网络插件后，然后重新执行安装  
``` bash
ansible-playbook /etc/ansible/06.network.yml
```

## 验证新网络插件

参照[calico](../06.calico.md) [flannel](../06.flannel.md) [kube-router](../06.kube-router.md)

## 已知BUG

如果现有网络是kube-router, 按上述步骤完成替换成其他网络时，需要额外执行一次pod重建：  
``` bash
ansible-playbook /etc/ansible/tools/clean_k8s_network.yml -t reload_pods
```
