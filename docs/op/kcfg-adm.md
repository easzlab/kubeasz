# 管理客户端kubeconfig 

默认 k8s集群安装成功后生成客户端kubeconfig，它拥有集群管理的所有权限（不要将这个admin权限、50年期限的kubeconfig流露出去）；而我们经常需要将限定权限、限定期限的kubeconfig 分发给普通用户；利用cfssl签发自定义用户证书和k8s灵活的rbac权限绑定机制，ezctl 工具封装了这个功能。

## 使用帮助

```
ezctl help kcfg-adm
Usage: ezctl kcfg-adm <cluster> <args>
available <args>:
    -A     to add a client kubeconfig with a newly created user
    -D     to delete a client kubeconfig with the existed user
    -L     to list all of the users
    -e     to set expiry of the user certs in hours (ex. 24h, 8h, 240h)
    -t     to set a user-type (admin or view)
    -u     to set a user-name prefix

examples: ./ezctl kcfg-adm test-k8s -L
          ./ezctl kcfg-adm default -A -e 240h -t admin -u jack
          ./ezctl kcfg-adm default -D -u jim-202101162141
```

- 可以设置过期时间
- 可以设置权限：管理员权限（admin）和只读权限（view）

## 使用举例

- 1.查看集群k8s-01当前自定义kubeconfig

```
ezctl kcfg-adm k8s-01 -L
2021-01-24 16:32:43 INFO list-kcfg k8s-01
2021-01-24 16:32:43 INFO list-kcfg in cluster:k8s-01

USER                           TYPE            EXPIRY(+8h if in Asia/Shanghai)
---------------------------------------------------------------------------------

2021-01-24 16:32:43 INFO list-kcfg k8s-01 success
```
初始情况下列表为空

- 2.增加集群k8s-01一个自定义用户kubeconfig，用户名user01，期限24h，只读权限

```
ezctl kcfg-adm k8s-01 -A -u user01 -e 24h -t view
2021-01-24 17:32:33 INFO add-kcfg k8s-01
2021-01-24 17:32:33 INFO add-kcfg in cluster:k8s-01 with user:user01-202101241732

PLAY [localhost] *****************************************************************************************************

...（此处省略输出） 

TASK [deploy : debug] ************************************************************************************************
ok: [localhost] => {
    "msg": "查看user01-202101241732自定义kubeconfig：/etc/kubeasz/clusters/k8s-01/ssl/users/user01-202101241732.kubeconfig"
}

PLAY RECAP ***********************************************************************************************************
localhost                  : ok=12   changed=10   unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

2021-01-24 17:32:41 INFO add-kcfg k8s-01 success
```
生成的kubeconfig位于 /etc/kubeasz/clusters/k8s-01/ssl/users/user01-202101241732.kubeconfig

- 3.再增加一个用户user02，期限240h，admin权限

```
ezctl kcfg-adm k8s-01 -A -u user02 -e 240h -t admin
2021-01-24 18:38:47 INFO add-kcfg k8s-01
2021-01-24 18:38:47 INFO add-kcfg in cluster:k8s-01 with user:user02-202101241838

PLAY [localhost] *****************************************************************************************************

...（此处省略输出）

TASK [deploy : debug] ************************************************************************************************
ok: [localhost] => {
    "msg": "查看user02-202101241838自定义kubeconfig：/etc/kubeasz/clusters/k8s-01/ssl/users/user02-202101241838.kubeconfig"
}

PLAY RECAP ***********************************************************************************************************
localhost                  : ok=12   changed=9    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

2021-01-24 18:38:55 INFO add-kcfg k8s-01 success
```

- 4.再次查看集群k8s-01当前自定义kubeconfig

```
ezctl kcfg-adm k8s-01 -L
2021-01-24 18:40:30 INFO list-kcfg k8s-01
2021-01-24 18:40:30 INFO list-kcfg in cluster:k8s-01

USER                           TYPE            EXPIRY(+8h if in Asia/Shanghai)
---------------------------------------------------------------------------------
user02-202101241838            cluster-admin   2021-02-03T10:34:00Z
user01-202101241732            view            2021-01-25T09:28:00Z

2021-01-24 18:40:31 INFO list-kcfg k8s-01 success
```

- 5.删除user01-202101241732 权限

``` bash
ezctl kcfg-adm k8s-01 -D -u user01-202101241732
2021-01-24 21:41:50 INFO del-kcfg k8s-01
2021-01-24 21:41:50 INFO del-kcfg in cluster:k8s-01 with user:user01-202101241732
clusterrolebinding.rbac.authorization.k8s.io "crb-user01-202101241732" deleted
2021-01-24 21:41:50 INFO del-kcfg k8s-01 success

ezctl kcfg-adm k8s-01 -L
2021-01-24 21:42:02 INFO list-kcfg k8s-01
2021-01-24 21:42:02 INFO list-kcfg in cluster:k8s-01

USER                           TYPE            EXPIRY(+8h if in Asia/Shanghai)
---------------------------------------------------------------------------------
user02-202101241838            cluster-admin   2021-02-03T10:34:00Z

2021-01-24 21:42:02 INFO list-kcfg k8s-01 success
```
