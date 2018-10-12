# 个性化集群参数配置

简单来说，对于刚接触项目者，如"快速指南"说明，只需要：

- **1** 个配置：`/etc/ansible/hosts`
- **1** 键安装：`ansible-playbook /etc/ansilbe/90.setup.yml`

具体来讲 `kubeasz`创建集群主要在以下两个地方进行配置：

- ansible hosts 文件（模板在examples目录）：集群主要节点定义和主要参数配置、全局变量
- roles/xxx/vars/main.yml 文件：其他参数配置或者部分组件附加参数

## ansible hosts

项目在[快速指南](quickStart.md)或者[集群规划与安装概览](00-planning_and_overall_installing.md)已经介绍过，主要包括集群节点定义和集群范围的主要参数配置；目前提供三种集群部署模板。

- 尽量保持配置简单灵活
- 尽量保持配置项稳定

## roles/xxx/vars/main.yml

主要包括集群某个具体组件的个性化配置，具体组件的配置项可能会不断增加；项目初始时该配置与默认配置(`roles/xxx/defaults/main.yml`)一致；因 ansilbe 变量优先级关系，后续如果对 roles/xxx/vars/main.yml变量修改，那么它将覆盖默认配置。

- 需要初始化时使用 `ansible-playbook /etc/ansilbe/tools/init_vars.yml` 生成
- 确保在不做任何配置更改情况下可以使用默认值创建集群
- 被.gitignore忽略，修改后项目目录能够保持干净(`git status | clean`)

