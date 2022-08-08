# 安装 Kuboard

## Kuboard 介绍

Kuboard 是一款免费的 Kubernetes 管理工具，提供了丰富的功能：

* Kubernetes 多集群管理
* Kubernetes 基本管理功能
  * 节点管理
  * 名称空间管理
  * 存储类/存储卷管理
  * 控制器（Deployment/StatefulSet/DaemonSet/CronJob/Job/ReplicaSet）管理
  * Service/Ingress 管理
  * ConfigMap/Secret 管理
  * CustomerResourceDefinition 管理
* Kubernetes 问题诊断
  * Top Nodes / Top Pods
  * 事件列表及通知
  * 容器日志及终端
  * KuboardProxy (kubectl proxy 的在线版本)
  * PortForward (kubectl port-forward 的快捷版本)
  * 复制文件 （kubectl cp 的在线版本）
* 认证与授权
  * Github/GitLab 单点登录
  * KeyCloak 认证
  * LDAP 认证
  * 完整的 RBAC 权限管理
* Kuboard 特色功能
  * Kuboard 官方套件
    * Grafana+Prometheus 资源监控
    * Grafana+Loki+Promtail 日志聚合
  * Kuboard 自定义名称空间布局
  * Kuboard 中英文语言包

<p>
  <a aria-label="github" href="https://starchart.cc/eip-work/kuboard-press" target="_blank">
    <img src="https://badgen.net/github/stars/eip-work/kuboard-press?label=github stars"/>
  </a>
</p>

<a href="https://starchart.cc/eip-work/kuboard-press" target="_blank">
  <img src="https://starchart.cc/eip-work/kuboard-press.svg" alt="Kubernetes教程_Kuboard_Github_Star" style="height: 320px; width: 540px;">
</a>

点击这里可以查看 [Kuboard 的安装文档](https://kuboard.cn/install/v3/install.html)

## 在线演示

<div>
  在线演示环境中，您具备 <span style="color: red; font-weight: bold">只读</span> 权限，只能体验 Kuboard 的一部分功能。<br/>
</div>
<div style="padding: 10px; border: 1px solid #eee; border-radius: 10px; margin: 10px 0px; background-color: #fafafa;">
  <a href="http://demo.kuboard.cn" target="_blank">http://demo.kuboard.cn</a> <br/>
  <div style="width: 60px; display: inline-block; margin-top: 5px;">用&nbsp;户</div>
  demo <br/>
  <div style="width: 60px; display: inline-block;">密&nbsp;码</div>
  demo123
</div>

## 特点介绍

相较于 Kubernetes Dashboard 等其他 Kubernetes 管理界面，Kuboard 的主要特点有：

* 多种认证方式

  Kuboard 可以使用内建用户库、gitlab / github 单点登录或者 LDAP 用户库进行认证，避免管理员将 ServiceAccount 的 Token 分发给普通用户而造成的麻烦。使用内建用户库时，管理员可以配置用户的密码策略、密码过期时间等安全设置。

  ![Kuboard登录界面](https://kuboard.cn/images/intro.assets/image-20210405162940278.png)

* 多集群管理

  管理员可以将多个 Kubernetes 集群导入到 Kuboard 中，并且通过权限控制，将不同集群/名称空间的权限分配给指定的用户或用户组。

  ![Kuboard集群列表页](https://kuboard.cn/images/intro.assets/image-20210405164029151.png)

* 微服务分层展示

  在 Kuboard 的名称空间概要页中，以经典的微服务分层方式将工作负载划分到不同的分层，更加直观地展示微服务架构的结构，并且可以为每一个名称空间自定义名称空间布局。

  ![Kuboard-微服务分层](https://kuboard.cn/images/intro.assets/image-20210405164532452.png)

* 工作负载的直观展示

  Kuboard 中将 Deployment 的历史版本、所属的 Pod 列表、Pod 的关联事件、容器信息合理地组织在同一个页面中，可以帮助用户最快速的诊断问题和执行各种相关操作。

  ![Kuboard-工作负载详情](https://kuboard.cn/images/intro.assets/image-20210405180147614.png)

* 工作负载编辑

  Kuboard 提供了图形化的工作负载编辑界面，用户无需陷入繁琐的 YAML 文件细节中，即可轻松完成对容器的编排任务。支持的 Kubernetes 对象类型包括：Node、Namespace、Deployment、StatefulSet、DaemonSet、Secret、ConfigMap、Service、Ingress、StorageClass、PersistentVolumeClaim、LimitRange、ResourceQuota、ServiceAccount、Role、RoleBinding、ClusterRole、ClusterRoleBinding、CustomResourceDefinition、CustomResource 等各类常用 Kubernetes 对象，

  ![Kuboard-工作负载编辑](https://kuboard.cn/images/intro.assets/image-20210405180800712.png)

* 存储类型支持

  在 Kuboard 中，可以方便地对接 NFS、CephFS 等常用存储类型，并且支持对 CephFS 类型的存储卷声明执行扩容和快照操作。

  ![Kuboard-存储类](https://kuboard.cn/images/intro.assets/image-20210405181928653.png)

* 丰富的互操作性

  可以提供许多通常只在 `kubectl` 命令行界面中才提供的互操作手段，例如：

  * Top Nodes / Top Pods
  * 容器的日志、终端
  * 容器的文件浏览器（支持从容器中下载文件、上传文件到容器）
  * KuboardProxy（在浏览器中就可以提供 `kubectl proxy` 的功能）

  ![Kuboard-文件浏览器](https://kuboard.cn/images/intro.assets/image-20210405182805543.png)

* 套件扩展

  Kuboard 提供了必要的套件库，使得用户可以根据自己的需要扩展集群的管理能力。当前提供的套件有：

  * 资源层监控套件，基于 Prometheus / Grafana 提供 K8S 集群的监控能力，可以监控集群、节点、工作负载、容器组等各个级别对象的 CPU、内存、网络、磁盘等资源的使用情况；
  * 日志聚合套件，基于 Grafana / Loki / Promtail 实现日志聚合；
  * 存储卷浏览器，查看和操作存储卷中的内容；

  ![Kuboard-套件扩展](https://kuboard.cn/images/intro.assets/image-20210405183652378.png)



访问 Kuboard 网站 https://kuboard.cn 可以加入 Kuboard 社群，并获得帮助。
