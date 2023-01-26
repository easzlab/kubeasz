## kubeasz 3.5.1

kubeasz 3.5.1 发布，k8s v1.26小版本更新，组件更新以及一些bugfix。

### 版本更新

- k8s: v1.26.1
- calico: v3.24.5
- chrony: 4.3
- containerd: v1.6.14
- docker: 20.10.22
- keepalived: 2.2.7
- nginx: 1.22.1
- harbor: v2.6.3

### 支持设置k8s nodename

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

- 更新etcd文档 (#1218) by itswl
- fix: start/stop scripts for ex-lb
- fix: 'ezctl'-ignore /usr/bin/python link existed warning
- 更新高可用架构图及核心安装文档
