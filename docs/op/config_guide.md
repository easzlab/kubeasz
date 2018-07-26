# 个性化集群参数配置

`kubeasz`创建集群主要在以下两个地方进行配置：

- ansible hosts 文件（模板在examples目录）：集群主要节点定义和主要参数配置
- roles/xxx/vars/main.yml 文件：其他参数配置或者部分组件参数配置

这些文件都在.gitignore忽略范围，因此修改后项目目录能够保持`git status | clean`

## ansible hosts

项目尽量保持`ansible hosts`简单、灵活，在[快速指南](../quickStart.md)或者[集群规划与安装概览](../00-集群规划和基础参数设定.md)已经介绍过，主要包括集群节点定义和集群范围的主要参数配置；目前提供三种集群部署模板。

尽量保持配置项稳定。

## roles/xxx/vars/main.yml

主要包括集群某个具体组件的个性化配置，具体组件的配置项可能会不断增加；项目初始时该配置与 roles/xxx/defaults/main.yml 一致，确保在不做任何配置情况下可以使用默认值创建集群；因 ansilbe 变量优先级关系，后续如果对 roles/xxx/vars/main.yml变量修改，那么它将覆盖默认配置。

