## 升级注意事项

### v1.8 >>> v1.9

+ 1.下载最新项目代码 `cd /etc/ansible && git pull origin master`
+ 2.下载新的二进制 `k8s.190.tar.gz` 解压并覆盖 `/etc/ansible/bin/` 目录下文件
+ 3.更新集群 `cd /etc/ansible && ansible-playbook 90.setup.yml`

注1：升级过程会短暂中断集群中已经运行的应用；如果你想要零中断升级，可以在熟悉项目安装原理基础上自行尝试，或者关注后续项目[使用指南]中的文档更新

注2：k8s集群v1.8升级v1.9.0，目前测试不用修改任何服务参数，只要替换二进制文件；
