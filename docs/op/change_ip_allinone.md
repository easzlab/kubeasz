# 修改AIO 部署的IP
前两天在项目[ISSUES #201](https://github.com/gjmzj/kubeasz/issues/201)看到有人提：`在虚拟机A装了allinone，并搭建一套开发环境，我想通过copy A出来一套B然后交给别人测试`，觉得这个场景蛮有用，就写了这个文档和对应的脚本，希望对各位有帮助，也可以熟悉kubeasz的安装逻辑。

首先，因为kubeasz创建的集群都是TLS双向认证的，所以修改host ip地址比想象中要复杂很多。具体步骤可以参考[脚本](../../tools/change_ip_aio.yml)中的注释内容。

- 本操作指南仅适用于测试交流

## 操作步骤
前提 ：一个运行正常的allinone部署在虚机，关机后复制给别人使用，新虚机开机后如果需要修改IP，请执行如下步骤：

- 0.拉取最新项目代码：`git pull origin master`
- 1.修改ansible hosts文件：`sed -i 's/$OLD_IP/$NEW_IP/g' /etc/ansible/hosts`
- 2.配置ssh免密码登陆：`ssh-copy-id $NEW_IP` 按提示完成
- 3.检查下修改是否成功，并且能够成功执行 `ansible all -m ping`
- 4.以上步骤完成后，执行 `ansible-playbook /etc/ansible/tools/change_ip_aio.yml`

执行成功即可，请自己验证原先集群中各应用是否正常。
