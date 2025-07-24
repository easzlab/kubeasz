# nacos 部署

参考 https://nacos.io/docs/v2.4/manual/admin/deployment/deployment-overview/

## 前置准备

- 创建openebs 提供动态pvc存储
- 安装 mysql 数据库，初始化建表语句

## 验证Nacos集群配置文件信息

``` bash
for i in 0 1 2; do echo nacos-$i; kubectl exec nacos-$i -- cat conf/cluster.conf; done
nacos-0
Defaulted container "nacos" out of: nacos, peer-finder-plugin-install (init)
#2025-07-22T17:12:41.878
nacos-0.nacos-headless.default.svc.cluster.local:8848
nacos-1.nacos-headless.default.svc.cluster.local:8848
nacos-2.nacos-headless.default.svc.cluster.local:8848
nacos-1
Defaulted container "nacos" out of: nacos, peer-finder-plugin-install (init)
#2025-07-22T17:12:53.913
nacos-0.nacos-headless.default.svc.cluster.local:8848
nacos-1.nacos-headless.default.svc.cluster.local:8848
nacos-2.nacos-headless.default.svc.cluster.local:8848
nacos-2
Defaulted container "nacos" out of: nacos, peer-finder-plugin-install (init)
#2025-07-22T17:12:57.963
nacos-0.nacos-headless.default.svc.cluster.local:8848
nacos-1.nacos-headless.default.svc.cluster.local:8848
nacos-2.nacos-headless.default.svc.cluster.local:8848
```

## 访问nacos 控制台 

``` bash
http://${nodeIp}:${nodePort}/nacos
```

用户名：nacos 密码：Nacos1234!（首次登录时初始化设置）
