# gitlab CI/CD 基础

gitlab-ci 兼容 travis ci 格式，也是最流行的 CI 工具之一；本文讲解利用 gitlab, gitlab-runner, docker, harbor, kubernetes 等流行开源工具搭建一个自动化CI/CD流水线；举的例子以简单实用为原则，暂时没有选用 dind（docker in dockers）打包、gitlab Auto DevOps 等方式。一个最简单的流水线如下：

- 代码提交 --> 镜像构建 --> 部署测试 --> 部署生产

## 前提条件

- 正常运行的 gitlab: [安装 gitlab 文档](gitlab-install.md)
- 若干虚机运行 gitlab-runner: 运行自动化流水线 pipeline
- 正常运行的容器仓库：[安装 Harbor 文档](../harbor.md)
- 正常运行的 k8s 集群：可以是自建/公有云提供商
- 了解代码管理流程 gitflow 等

## 1.准备测试项目代码

假设你要开发一个 spring boot 项目；先登陆你的 gitlab 账号，创建项目，上传你的代码；项目根目录看起来如下：

```
-rw-r--r-- 1 root root    44 Jan  2 16:38 eclipse.bat
drwxr-xr-x 8 root root  4096 Jan  7 15:29 .git/
-rw-r--r-- 1 root root   276 Jan  7 08:44 .gitignore
drwxr-xr-x 3 root root  4096 Jan  7 08:44 example-api/
drwxr-xr-x 3 root root  4096 Jan  7 08:44 example-biz/
drwxr-xr-x 3 root root  4096 Jan  2 16:38 example-dal/
drwxr-xr-x 3 root root  4096 Jan  2 16:38 example-web/
-rw-r--r-- 1 root root    54 Jan  2 16:38 install.bat
-rw-r--r-- 1 root root 10419 Jan  2 16:38 pom.xml
```
传统做法是在本地配置好相关环境后使用 mvn 编译生成jar包，然后测试运行jar，这里我们要打包成 docker 镜像，并创建 CI/CD 流水线：如下示例，在项目根目录创建2个文件夹及里面文件

``` bash
dockerfiles        ### 新增文件夹用来 docker 镜像打包
└── Dockerfile     # 定义 docker 镜像
.ci                ### 新增文件夹用来存放 CI/CD 相关内容
├── app.yaml       # k8s 平台的应用部署文件
├── config.sh      # 配置替换脚本
└── gitlab-ci.yml  # gitlab-ci 的主配置文件
```

## 2.准备 docker 镜像描述文件 Dockerfile

我们把 Dockerfile 放在独立目录下，java spring boot 应用可以这样写：

``` bash
FROM openjdk:8-jdk-alpine
VOLUME /tmp
COPY *.jar app.jar         # 这里 *.jar 包就是后续在cicd pipeline 过程中 mvn 生成的jar包移动到此目录
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-jar","/app.jar"]
```

## 3.准备 CI/CD 相关脚本和文件

本地安装 gitlab 的同时也会安装帮助文档，非常有用，请阅读如下文档（假设你本地gitlab使用的域名`http://gitlab.test.com`）

- gitlab-ci 基本概念 http://gitlab.test.com/help/ci/README.md
- 变量 http://gitlab.test.com/help/ci/variables/README.md

目录`.ci`下面的三个文件`app.yaml`, `config.sh`, `gitlab-ci.yml`是互相关联的；整个流程中使用到的变量分为三种：

- 第一种是gitlab自身预定义变量（比如项目名: CI_PROJECT_NAME，流水线ID: CI_PIPELINE_ID）；无需更改；
- 第二种是在gitlab-ci.yml文件中定义的变量，一般是少量的自定义变量；按需少量改动；
- 第三种是用户可以在项目web界面配置的变量：“Settings”>"CI/CD">"Variables"，本示例项目用到该类型变量举例：

|变量|值|注解|
|:-|:-|:-|
|BETA_APP_REP|1|beta环境应用副本数|
|BETA_DB_HOST|1.1.1.1:3306|beta环境应用连接数据库主机|
|BETA_DB_PWD|xxxx|beta环境数据库连接密码|
|BETA_DB_USR|xxxx|beta环境数据库连接用户|
|BETA_REDIS_HOST|1.1.1.2|beta环境redis主机|
|BETA_REDIS_PORT|6379|beta环境redis端口|
|BETA_REDIS_PWD|xxxx|beta环境redis密码|
|BETA_HARBOR|1.1.1.3|beta环境镜像仓库地址|
|BETA_HARBOR_PWD|xxxx|beta环境镜像仓库密码|
|BETA_HARBOR_USR|xxxx|beta环境镜像仓库用户|
|PROD_APP_REP|2|prod环境应用副本数|
|PROD_DB_HOST|2.2.2.1:3306|prod环境应用连接数据库主机|
|PROD_DB_PWD|xxxx|prod环境数据库连接密码|
|PROD_DB_USR|xxxx|prod环境数据库连接用户|
|PROD_REDIS_HOST|2.2.2.2|prod环境redis主机|
|PROD_REDIS_PORT|6379|prod环境redis端口|
|PROD_REDIS_PWD|xxxx|prod环境redis密码|
|PROD_HARBOR|2.2.2.3|prod环境镜像仓库地址|
|PROD_HARBOR_PWD|xxxx|prod环境镜像仓库密码|
|PROD_HARBOR_USR|xxxx|prod环境镜像仓库用户|
|...|...|根据项目需要自行添加设置|

掌握了以上基础知识，可以开始以下三个任务：

- 3.1[配置 gitlab-ci.yml](gitlab-ci.yml.md), 整个CI/CD的主配置文件，定义所有的CI/CD阶段和每个阶段的任务
- 3.2[配置 config.sh](config.sh.md)，根据不同分支/环境替换不同的应用程序变量（对应上述第三种变量）
- 3.3[配置 app.yaml](app.yaml.md)，K8S应用部署简单模板，替换完成后可以部署到测试/生产的K8S平台上

## to be continued
