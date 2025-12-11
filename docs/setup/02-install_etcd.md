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
  etcdctl \
  --endpoints=https://${ip}:2379  \
  --cacert=/etc/kubernetes/ssl/ca.pem \
  --cert=/etc/kubernetes/ssl/etcd.pem \
  --key=/etc/kubernetes/ssl/etcd-key.pem \
  endpoint health; done

# 预期结果
https://192.168.1.1:2379 is healthy: successfully committed proposal: took = 2.210885ms
https://192.168.1.2:2379 is healthy: successfully committed proposal: took = 2.784043ms
https://192.168.1.3:2379 is healthy: successfully committed proposal: took = 3.275709ms

for ip in ${NODE_IPS}; do
  etcdctl \
  --endpoints=https://${ip}:2379  \
  --cacert=/etc/kubernetes/ssl/ca.pem \
  --cert=/etc/kubernetes/ssl/etcd.pem \
  --key=/etc/kubernetes/ssl/etcd-key.pem \
  --write-out=table endpoint status; done

# 预期结果
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
|          ENDPOINT          |        ID        | VERSION | STORAGE VERSION | DB SIZE | IN USE | PERCENTAGE NOT IN USE | QUOTA  | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS | DOWNGRADE TARGET VERSION | DOWNGRADE ENABLED |
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
| https://192.168.1.1:2379   | 5f64925bd78a482c |   3.6.4 |           3.6.0 |   38 MB |  28 MB |                   28% | 8.6 GB |      true |      false |       269 |    6582307 |            6582307 |        |                          |             false |
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
+----------------------------+-----------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
|          ENDPOINT          |       ID        | VERSION | STORAGE VERSION | DB SIZE | IN USE | PERCENTAGE NOT IN USE | QUOTA  | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS | DOWNGRADE TARGET VERSION | DOWNGRADE ENABLED |
+----------------------------+-----------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
| https://192.168.1.2:2379   | 18e1b1602639adb |   3.6.4 |           3.6.0 |   37 MB |  28 MB |                   25% | 8.6 GB |     false |      false |       269 |    6582307 |            6582307 |        |                          |             false |
+----------------------------+-----------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
|          ENDPOINT          |        ID        | VERSION | STORAGE VERSION | DB SIZE | IN USE | PERCENTAGE NOT IN USE | QUOTA  | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS | DOWNGRADE TARGET VERSION | DOWNGRADE ENABLED |
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
| https://192.168.1.3:2379   | 3d375f7546465b4e |   3.6.4 |           3.6.0 |   37 MB |  28 MB |                   26% | 8.6 GB |     false |      false |       269 |    6582308 |            6582308 |        |                          |             false |
+----------------------------+------------------+---------+-----------------+---------+--------+-----------------------+--------+-----------+------------+-----------+------------+--------------------+--------+--------------------------+-------------------+
```

- 所有节点可达：etcdctl endpoint health 对所有三个节点都返回 healthy。
- 有且仅有一个领导者：etcdctl endpoint status 显示一个节点 is leader: true，另外两个节点 is leader: false。
- Raft 任期一致：所有三个节点的 raft term 值完全相同。
- Raft 索引同步：所有节点的 raft index 值相差不大（跟随者与领导者的差距在可接受范围内）。
- 无活跃告警：etcdctl alarm list 返回空。
- 节点间网络稳定：没有频繁的领导者切换（通过监控 etcd_server_leader_changes_seen_total 指标）。
- 磁盘空间充足：没有 NOSPACE 告警，且磁盘使用率在安全阈值内（例如低于80%）。

### 磁盘性能

快速的磁盘是 etcd 部署性能和稳定性的最关键因素。

磁盘速度慢会增加 etcd 请求延迟，并可能损害集群稳定性。由于 etcd 的共识协议依赖于将元数据持久地存储到日志中，因此大多数 etcd 集群成员必须将每个请求写入磁盘。此外，etcd 还会逐步将其状态检查点写入磁盘，以便截断此日志。如果这些写入耗时过长，心跳可能会超时并触发选举，从而损害集群的稳定性。通常，要判断磁盘速度是否足以满足 etcd 的要求，可以使用fio等基准测试工具。

etcd 对磁盘写入延迟非常敏感。通常需要 50 的顺序 IOPS（例如，7200 RPM 磁盘）。对于负载较重的集群，建议使用 500 的顺序 IOPS（例如，典型的本地 SSD 或高性能虚拟化块设备）。请注意，大多数云提供商发布的是并发 IOPS，而不是顺序 IOPS；发布的并发 IOPS 可能比顺序 IOPS 高出 10 倍。要测量实际的顺序 IOPS，我们建议使用磁盘基准测试工具，例如diskbench或fio。

``` bash
# 测试示例
mkdir test-data
fio --rw=write --ioengine=sync --fdatasync=1 --directory=test-data --size=2200m --bs=2300 --name=mytest
```



[后一篇](03-container_runtime.md)
