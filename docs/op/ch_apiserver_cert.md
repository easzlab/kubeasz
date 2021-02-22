# 修改 APISERVER（MASTER）证书

`kubeasz` 创建集群后，APISERVER（MASTER）证书默认 CN 包含如下`域名`和`IP`：参见`roles/kube-master/templates/kubernetes-csr.json.j2`

```
  "hosts": [
    "127.0.0.1",
{% if groups['ex_lb']|length > 0 %}
    "{{ hostvars[groups['ex_lb'][0]]['EX_APISERVER_VIP'] }}",
{% endif %}
    "{{ inventory_hostname }}",
    "{{ CLUSTER_KUBERNETES_SVC_IP }}",
{% for host in MASTER_CERT_HOSTS %}
    "{{ host }}",
{% endfor %}
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local"
  ],
```

有的时候（比如apiserver地址通过边界防火墙的NAT转换成公网IP访问，或者需要添加公网域名访问）我们需要在 APISERVER（MASTER）证书中添加一些`域名`或者`IP`，可以方便操作如下：

## 1.修改配置文件`roles/kube-master/defaults/main.yml`

``` bash
# k8s 集群 master 节点证书配置，可以添加多个ip和域名（比如增加公网ip和域名）
MASTER_CERT_HOSTS:
  - "10.1.1.1"
  - "k8s.test.io"
  #- "61.182.11.41"
  #- "www.test.com"
```

## 2.执行新证书生成即可

``` bash
$ ansible-playbook 04.kube-master.yml -t change_cert
# 新证书生效需要重启kube-apiserver.service
$ ansible-playbook 04.kube-master.yml -t restart_master
```
