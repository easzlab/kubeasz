# 安装 gitlab [Deprecated]

gitlab 是深受企业用户喜爱的基于 git 的代码管理系统。安装 gitlab 最理想的方式是利用 gitlab charts 部署到 k8s 集群上，但此方式还未成熟，期待后续推出更成熟稳定版本；本文使用 Docker 方式安装 gitlab:

- 环境：Ubuntu 16.04，虚机内存/CPU/存储请根据实际使用情况配置，一般`4C/8G/200G`足够
- 安装 docker: 18.06.1-ce

## 准备启动脚本

``` bash
$ cat > gitlab-setup.sh << EOF
#!/bin/bash
# 注意：设置 gitlab_shell_ssh_port 是为了后续可以使用 SSH 方式访问你的项目
docker run --detach \\
    --hostname gitlab.test.com \\
    --env GITLAB_OMNIBUS_CONFIG="external_url 'http://gitlab.test.com/'; gitlab_rails['gitlab_shell_ssh_port'] = 6022;" \\
    --publish 443:443 --publish 80:80 --publish 6022:22 \\
    --name gitlab \\
    --restart always \\
    --volume /srv/gitlab/config:/etc/gitlab \\
    --volume /srv/gitlab/logs:/var/log/gitlab \\
    --volume /srv/gitlab/data:/var/opt/gitlab \\
    docker.mirrors.ustc.edu.cn/gitlab/gitlab-ce:11.2.2-ce.0
EOF
```
执行启动脚本：`sh gitlab-setup.sh` 执行成功后，等待数分钟可以看到

```
$ docker ps -a
CONTAINER ID        IMAGE                                                 COMMAND             CREATED             STATUS                   PORTS                                                            NAMES
4f9d5f97f494        docker.mirrors.ustc.edu.cn/gitlab/gitlab-ce:11.2.2-ce.0   "/assets/wrapper"   9 minutes ago       Up 9 minutes (healthy)   0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp, 0.0.0.0:6022->22/tcp   gitlab
```
## 配置 gitlab

```
$ docker exec -it gitlab vi /etc/gitlab/gitlab.rb
```
请阅读后修改（因为前面docker run 已经指定了必要参数，可以不修改，后续有需要再修改），修改保存以后需要重启容器

```
$ docker restart gitlab
```
## 首次访问 gitlab

使用域名`gitlab.test.com`或者该主机 IP 首次登录时会要求设置 root 用户的密码，完成后就可以用 root 和新设密码登录；然后按需创建 Group, User, Projects等，还有相关配置。

## 备份数据

无论是企业、组织、个人都十分重视代码资产，之前我们的 gitlab 安装是单机版的，虽然可以有硬盘 raid 等保护，还有是丢失 gitlab 数据和配置的风险，因此我们有必要再做一些备份操作。这里利用 crontab 定期执行 rsync 命令备份到其他服务器。

``` bash
# 创建备份脚本
cat > /root/gitlab-backup.sh << EOF
#!/bin/bash
# 请事先配置 gitlab 服务器到备份服务器的免密码 ssh 登录
rsync -av --delete /srv/gitlab/config '-e ssh -l root' 192.168.1.xx:/backup_gitlab/config
rsync -av --delete /srv/gitlab/data '-e ssh -l root' 192.168.1.xx:/backup_gitlab/data
EOF

# 创建并应用 crontab
cat > /etc/cron.d/gitlab-backup << EOF
## 每3个小时同步备份一次，具体根据需要修改
11 */3 * * * root bash /root/gitlab-backup.sh > /root/gitlab/sync.log 2>&1
EOF
```
如果 gitlab 服务器真的出现不可恢复的故障，丢失数据，那么至少保留有3小时前的备份，利用备份的文件，同样再用 docker 挂载 volume的方式运行，这样就可以恢复原 gitlab 服务运行。

## 升级 gitlab

因为前面使用了 docker 方式安装，因此 gitlab 升级很方便。

- 升级前停止/删除容器：`$ docker stop gitlab && docker rm gitlab`
- 如上节执行备份数据
- 修改 gitlab-setup.sh 指定新的版本，执行该脚本

## 参考

- 1.[Install GitLab with Docker](https://docs.gitlab.com/omnibus/docker/)
