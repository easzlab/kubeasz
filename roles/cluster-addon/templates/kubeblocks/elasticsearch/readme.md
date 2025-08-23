# 常用维护操作

1. 查看集群健康概况
`curl -XGET "http://127.0.0.1:9200/_cluster/health?pretty"`
关键字段：
  ● green: 所有主分片和副本分片均正常。
  ● yellow: 主分片正常，部分副本分片未分配。
  ● red: 存在未分配的主分片。
  ● number_of_nodes: 当前在线节点数。
  ● active_shards: 活跃分片数量。
  ● unassigned_shards: 未分配的分片数（若 >0 需排查原因）。

2. 查看所有节点信息
`curl -XGET "http://127.0.0.1:9200/_cat/nodes?v"`
关键列：
  ● heap.percent: JVM堆内存使用率（>75% 需关注）。
  ● cpu: CPU使用率。
  ● role: 节点角色（如 mdi 表示主+数据+ingest节点）。

3. 列出所有索引及状态
`curl -XGET "http://127.0.0.1:9200/_cat/indices?v"`
  ● 关注 health 状态（红/黄/绿）及 docs.count（文档数）。

4. 检查索引分片分布
`curl -XGET "http://127.0.0.1:9200/_cat/indices/<index-name>?v&h=index,shard,prirep,state,node"`
  ● 确认主分片（p）和副本分片（r）是否均衡分布在节点间。

5. 分片分配详情
`curl -XGET "http://127.0.0.1:9200/_cat/shards?v"`
● 检查是否有 UNASSIGNED 分片，及其对应的索引和原因。

6. 查看未分配分片原因
`curl -XGET "http://127.0.0.1:9200/_cluster/allocation/explain?pretty"`
● 输出会提示分片未分配的具体原因（如磁盘不足、节点离线）。
● 磁盘空间不足：清理旧数据或扩容磁盘。
● 节点离线：恢复节点或手动分配分片。
● 配置限制：调整 cluster.routing.allocation 相关设置。

7. 实时查看线程池状态
curl -XGET "http://127.0.0.1:9200/_cat/thread_pool?v"
● 关注 bulk、search 队列是否堆积（queue > 0）。
