## 02-安装etcd集群.md

``` bash
roles/etcd
├── tasks
│   └── main.yml
└── templates
    ├── etcd-csr.json.j2
    └── etcd.service.j2
```
kuberntes 系统使用 etcd 存储所有数据，是最重要的组件之一，注意 etcd集群只能有奇数个节点(1,3,5...)，本文档使用3个节点做集群。

请在另外窗口打开[roles/etcd/tasks/main.yml](../roles/etcd/tasks/main.yml) 文件，对照看以下讲解内容。

### 下载etcd/etcdctl 二进制文件、创建证书目录

### 创建etcd证书请求 [etcd-csr.json.j2](../roles/etcd/templates/etcd-csr.json.j2)

首先判断下是否etcd 证书已经存在，如果已经存在就跳过证书生成步骤

``` bash
{
  "CN": "etcd",
  "hosts": [
    "127.0.0.1",
    "{{ inventory_hostname }}"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "HangZhou",
      "L": "XS",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
```
+ etcd使用对等证书，hosts 字段必须指定授权使用该证书的 etcd 节点 IP

### 创建证书和私钥

``` bash
cd /etc/etcd/ssl && {{ bin_dir }}/cfssl gencert \
        -ca={{ ca_dir }}/ca.pem \
        -ca-key={{ ca_dir }}/ca-key.pem \
        -config={{ ca_dir }}/ca-config.json \
        -profile=kubernetes etcd-csr.json | {{ bin_dir }}/cfssljson -bare etcd
```

###  创建etcd 服务文件 [etcd.service.j2](../roles/etcd/templates/etcd.service.j2)

先创建工作目录 /var/lib/etcd/

``` bash
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target
Documentation=https://github.com/coreos

[Service]
Type=notify
WorkingDirectory=/var/lib/etcd/
ExecStart={{ bin_dir }}/etcd \
  --name={{ NODE_NAME }} \
  --cert-file=/etc/etcd/ssl/etcd.pem \
  --key-file=/etc/etcd/ssl/etcd-key.pem \
  --peer-cert-file=/etc/etcd/ssl/etcd.pem \
  --peer-key-file=/etc/etcd/ssl/etcd-key.pem \
  --trusted-ca-file={{ ca_dir }}/ca.pem \
  --peer-trusted-ca-file={{ ca_dir }}/ca.pem \
  --initial-advertise-peer-urls=https://{{ inventory_hostname }}:2380 \
  --listen-peer-urls=https://{{ inventory_hostname }}:2380 \
  --listen-client-urls=https://{{ inventory_hostname }}:2379,http://127.0.0.1:2379 \
  --advertise-client-urls=https://{{ inventory_hostname }}:2379 \
  --initial-cluster-token=etcd-cluster-0 \
  --initial-cluster={{ ETCD_NODES }} \
  --initial-cluster-state=new \
  --data-dir=/var/lib/etcd
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
+ 完整参数列表请使用 `etcd --help` 查询
+ 注意etcd 即需要服务器证书也需要客户端证书，这里为方便使用一个peer 证书代替两个证书，更多证书相关请阅读 [01-创建CA证书和环境配置.md](01-创建CA证书和环境配置.md)
+ 注意{{ }} 中的参数与ansible hosts文件中设置对应
+ `--initial-cluster-state` 值为 `new` 时，`--name` 的参数值必须位于 `--initial-cluster` 列表中；

### 启动etcd服务

``` bash
systemctl daemon-reload && systemctl enable etcd && systemctl start etcd
```

### 验证etcd集群状态

+ systemctl status etcd 查看服务状态
+ journalctl -u etcd 查看运行日志
+ 在任一 etcd 集群节点上执行如下命令

``` bash
# 根据hosts中配置设置shell变量 $NODE_IPS
export NODE_IPS="192.168.1.1 192.168.1.2 192.168.1.3"
$ for ip in ${NODE_IPS}; do
  ETCDCTL_API=3 etcdctl \
  --endpoints=https://${ip}:2379  \
  --cacert=/etc/kubernetes/ssl/ca.pem \
  --cert=/etc/etcd/ssl/etcd.pem \
  --key=/etc/etcd/ssl/etcd-key.pem \
  endpoint health; done
```
预期结果：

``` text
https://192.168.1.1:2379 is healthy: successfully committed proposal: took = 2.210885ms
https://192.168.1.2:2379 is healthy: successfully committed proposal: took = 2.784043ms
https://192.168.1.3:2379 is healthy: successfully committed proposal: took = 3.275709ms
```
三台 etcd 的输出均为 healthy 时表示集群服务正常。


[前一篇](01-创建CA证书和环境配置.md) -- [后一篇](03-安装docker服务.md)
