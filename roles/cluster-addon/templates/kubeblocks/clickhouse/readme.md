## 常用操作

1. 客户端连接

`clickhouse-client -h clickhouse-cluster-clickhouse.db.svc -u admin --port 9000 --password`

2. 查看集群状态

`clickhouse-client --format=Pretty --query="SELECT * FROM system.clusters"`

3. 查看库表磁盘占用

```
clickhouse-client --format=Pretty --query="
SELECT database, table, formatReadableSize(sum(bytes)) AS size
FROM system.parts
GROUP BY database, table
ORDER BY sum(bytes) DESC"
```

4. 备份与恢复

```
## 备份
#!/bin/bash
DATE=$(date +%Y%m%d)
BACKUP_DIR=/data/clickhouse/backups/$DATE
mkdir -p $BACKUP_DIR

clickhouse-client --query="BACKUP DATABASE production_db TO Disk('backup_disk', '$BACKUP_DIR/production_db')"
echo "Backup completed at $BACKUP_DIR"

## 恢复
RESTORE DATABASE production_db FROM Disk('backup_disk', '/backups/20240101/production_db')
```

5. 慢查询

```
# 检查慢查询
SELECT query_id, user, query_duration_ms, normalized_query_hash 
FROM system.query_log
WHERE event_date = today() AND query_duration_ms > 10000
ORDER BY query_duration_ms DESC
LIMIT 20

# 终止异常查询
KILL QUERY WHERE query_id = 'abcd-efgh-ijkl'
```
