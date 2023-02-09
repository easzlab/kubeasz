## kubeasz 3.5.2

kubeasz 3.5.2 发布，解决3.5.1 版本中设置k8s_nodename的bug，以及其他一些fix。

### 支持设置k8s nodename

修复 ISSUE #1225，感谢 surel9

默认情况下kubeasz项目使用节点ip地址作为nodename，如果需要自定义设置支持两种方式：

- 1. 在clusters/xxxx/hosts 直接配置：比如

```
# work node(s), set unique 'k8s_nodename' for each node
[kube_node]
192.168.0.80 k8s_nodename=worker-01
192.168.0.79 k8s_nodename=worker-02
```

- 2. 在添加节点适合设置：比如

```
dk ezctl add-node xxxx 192.168.0.81 k8s_nodename=worker-03
```

特别注意：k8s_nodename 命名规范，只能由小写字母、数字、'-'、'.' 组成，并且开头和结尾必须为小写字母和数字
'k8s_nodename' must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'example.com')


### 其他

- 修复：prepare阶段如果安装系统包失败的错误不应该被忽略 
- 修复：清理节点时无法删除calico目录/var/run/calico
- 修复：deploy机器上调度的pod无法通信问题 issue #1224，感谢 bogeit

