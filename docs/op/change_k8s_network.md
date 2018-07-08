# 替换k8s集群的网络插件

有时候我们在测试环境的k8s集群中希望试用多种网络插件（calico/flannel/kube-router），又不希望每测试一次就全部清除集群然后重建，那么可能这个文档适合你。  
- WARNNING：重新安装k8s网络插件会短暂中断已有运行在k8s上的服务
  - 请在熟悉kubeasz的安装流程和k8s网络插件安装流程的基础上谨慎操作
  - 如果k8s集群已经运行庞大业务pod，重装网络插件时会引起所有pod的删除、重建，短时间内将给apiserver带来压力，可能引起master节点夯住

## 替换流程

kubeasz使用标准cni方式安装k8s集群的网络插件；cni负载创建容器网卡和IP分配（IPAM），不同的网络插件（calico,flannel等）创建容器网卡和IP分配方式不一样，所以在替换网络插件时候需要现有pod全部删除，然后自动按照新网络插件的方式重建pod网络；请参考[k8s网络插件章节](../06-安装网络组件.md)。

### 替换操作

替换网络插件操作很简单，只要两步：  
- 1.修改ansible hosts文件指定新网络插件
- 2.执行替换脚本 `ansible-playbook /etc/ansible/tools/change_k8s_network.yml`

对照脚本`change_k8s_network.yml` 讲解下大致流程为：  
a.根据实际运行情况，删除现有网络组件的daemonset pod  
b.如果现有组件是kube-router 需要进行一些额外清理  
c.暂停node相关服务，后面才可以进一步清理iptables等  
d.执行旧网络插件相关清理  
e.重新开启node相关服务  
f.安装新网络插件  
g.删除所有运行pod，然后等待自动重建  

## 验证新网络插件

参照[calico](../06.calico.md) [flannel](../06.flannel.md) [kube-router](../06.kube-router.md)

