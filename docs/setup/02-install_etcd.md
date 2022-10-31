## 02-安装etcd集群

kuberntes 集群使用 etcd 存储所有数据，是最重要的组件之一，注意 etcd集群需要奇数个节点(1,3,5...)，本文档使用3个节点做集群。

请在另外窗口打开[roles/etcd/tasks/main.yml](../../roles/etcd/tasks/main.yml) 文件，对照看以下讲解内容。

### 创建etcd证书

注意：证书是在部署节点创建好之后推送到目标etcd节点上去的，以增加ca证书的安全性

创建ectd证书请求 [etcd-csr.json.j2](../../roles/etcd/templates/etcd-csr.json.j2)

``` bash
{
  "CN": "etcd",
  "hosts": [
{% for host in groups['etcd'] %}
    "{{ host }}",
{% endfor %}
    "127.0.0.1"
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
+ etcd使用对等证书，hosts 字段必须指定授权使用该证书的 etcd 节点 IP，这里枚举了所有ectd节点的地址

###  创建etcd 服务文件 [etcd.service.j2](../../roles/etcd/templates/etcd.service.j2)

``` bash
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target
Documentation=https://github.com/coreos

[Service]
Type=notify
WorkingDirectory={{ ETCD_DATA_DIR }}
ExecStart={{ bin_dir }}/etcd \
  --name=etcd-{{ inventory_hostname }} \
  --cert-file={{ ca_dir }}/etcd.pem \
  --key-file={{ ca_dir }}/etcd-key.pem \
  --peer-cert-file={{ ca_dir }}/etcd.pem \
  --peer-key-file={{ ca_dir }}/etcd-key.pem \
  --trusted-ca-file={{ ca_dir }}/ca.pem \
  --peer-trusted-ca-file={{ ca_dir }}/ca.pem \
  --initial-advertise-peer-urls=https://{{ inventory_hostname }}:2380 \
  --listen-peer-urls=https://{{ inventory_hostname }}:2380 \
  --listen-client-urls=https://{{ inventory_hostname }}:2379,http://127.0.0.1:2379 \
  --advertise-client-urls=https://{{ inventory_hostname }}:2379 \
  --initial-cluster-token=etcd-cluster-0 \
  --initial-cluster={{ ETCD_NODES }} \
  --initial-cluster-state={{ CLUSTER_STATE }} \
  --data-dir={{ ETCD_DATA_DIR }} \
  --wal-dir={{ ETCD_WAL_DIR }} \
  --snapshot-count=50000 \
  --auto-compaction-retention=1 \
  --auto-compaction-mode=periodic \
  --max-request-bytes=10485760 \
  --quota-backend-bytes=8589934592
Restart=always
RestartSec=15
LimitNOFILE=65536
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
```

+ 完整参数列表请使用 `etcd --help` 查询
+ 注意etcd 即需要服务器证书也需要客户端证书，为方便使用一个peer 证书代替两个证书
+ `--initial-cluster-state` 值为 `new` 时，`--name` 的参数值必须位于 `--initial-cluster` 列表中
+ `--snapshot-count` `--auto-compaction-retention` 一些性能优化参数，请查阅etcd项目文档
+ 设置`--data-dir` 和`--wal-dir` 使用不同磁盘目录，可以避免磁盘io竞争，提高性能，具体请参考etcd项目文档

### 验证etcd集群状态

+ systemctl status etcd 查看服务状态
+ journalctl -u etcd 查看运行日志
+ 在任一 etcd 集群节点上执行如下命令

``` bash
# 根据hosts中配置设置shell变量 $NODE_IPS
export NODE_IPS="192.168.1.1 192.168.1.2 192.168.1.3"
for ip in ${NODE_IPS}; do
  ETCDCTL_API=3 etcdctl \
  --endpoints=https://${ip}:2379  \
  --cacert=/etc/kubernetes/ssl/ca.pem \
  --cert=/etc/kubernetes/ssl/etcd.pem \
  --key=/etc/kubernetes/ssl/etcd-key.pem \
  endpoint health; done
```
预期结果：

``` text
https://192.168.1.1:2379 is healthy: successfully committed proposal: took = 2.210885ms
https://192.168.1.2:2379 is healthy: successfully committed proposal: took = 2.784043ms
https://192.168.1.3:2379 is healthy: successfully committed proposal: took = 3.275709ms
```
三台 etcd 的输出均为 healthy 时表示集群服务正常。

[后一篇](03-container_runtime.md)
