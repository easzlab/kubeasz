# gitlab CI/CD 基础

gitlab-ci 兼容 travis ci 格式，也是最流行的 CI 工具之一；本文讲解利用 gitlab, gitlab-runner, docker, harbor, kubernetes 等流行开源工具搭建一个自动化CI/CD流水线；举的例子以简单实用为原则，没有选用 dind（docker in dockers）打包、gitlab Auto DevOps 等先进但未必成熟的方式。一个最简单的流水线如下：

- 代码提交 --> 镜像构建 --> 部署测试 --> 部署生产

## 前提条件

- 正常运行的 gitlab: [安装 gitlab 文档](gitlab-install.md)
- 若干虚机运行 gitlab-runner: 运行自动化流水线 pipeline
- 正常运行的容器仓库：[安装 Harbor 文档](../harbor.md)
- 正常运行的 k8s 集群：可以是自建/公有云提供商

##
