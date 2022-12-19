# Jenkins CI/CD

## 前言
本文档介绍如何快速通过K8s集群实现Jenkins 动态Slave CI/CD流程。

## 开始之前
在开始之前需要准备以下环境：
- k8s dns组件  
参考文档：[kubedns](kubedns.md)
- helm  
为了简化部署，通过helm来安装Jenkins，可参考文档：[helm](helm.md)
- 持久化存储  
这里使用**NFS**演示，参考文档：[cluster-storage](../setup/08-cluster-storage.md)。
如果k8s集群是部署在公有云，也可使用厂商的NAS等存储方案，项目中已集成支持阿里云NAS，其他的方案参考相关厂商文档

- Ingress Controller(nginx-ingress/traefik)  
默认是通过Ingress访问Jenkins，因此需要安装一种`Ingress Controller`。参考文档：[ingress](ingress.md)
- Gitlab 代码管理仓库  
用于提交代码后自动触发CI, 目前项目中还没有相关内容，可[参考官网](https://about.gitlab.com/installation/)进行安装。

## 安装Jenkins
执行以下命令快速安装：
```
helm install manifests/jenkins/ --name jenkins
```
如果通过/etc/kubeasz/roles/helm/helm.yml安装的helm，安装过程会出现如下错误

``` bash
E0703 08:40:22.376225   19888 portforward.go:331] an error occurred forwarding 41655 -> 44134: error forwarding port 44134 to pod 5098414beaaa07140a4ba3240690b1ce989ece01e5db33db65eec83bd64bdedf, uid : exit status 1: 2018/07/03 08:40:22 socat[19991] E write(5, 0x1aec120, 3424): Connection reset by peer
Error: transport is closing
```
请执行以下命令快速安装进行修复：
```
helm install --tls manifests/jenkins/ --name jenkins
```

由于初始化过程中，默认安装指定的插件，所以启动较慢，大概5-10分钟左右就可以启动完成了。  

部分默认配置说明：
**注**：以下配置都定义在`manifests/jenkins/values.yaml`文件中。
<table border="0">
    <tr>
        <td><b>字段</b></td>
        <td><b>说明</b></td>
        <td><b>默认值</b></td>
    </tr>
    <tr>
        <td>InstallPlugins</td>
        <td>初始化安装的插件</td>
        <td>
            <ul>
                <li>kubernetes:1.6.3</li>
                <li>workflow-aggregator:2.5</li>
                <li>workflow-job:2.21</li>
                <li>credentials-binding:1.16</li>
                <li>git:3.9.0</li>
                <li>gitlab:1.5.6</li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>HostName</td>
        <td>Ingress访问入口</td>
        <td>jenkins.local.com</td>
    </tr>
    <tr>
        <td>AdminPassword</td>
        <td>admin登录密码</td>
        <td>admin</td>
    </tr>
    <tr>
        <td>UpdateCenter</td>
        <td>插件下载镜像地址</td>
        <td>https://mirrors.tuna.tsinghua.edu.cn/jenkins</td>
    </tr>
    <tr>
        <td>StorageClass</td>
        <td>持久化存储SC</td>
        <td>nfs-dynamic-class</td>
    </tr>
</table>


## 配置Kubernetes plugin
登录Jenkins，点击左边导航`系统管理`——>`系统设置`，拖动到最下面可以看到`云——>Kubernetes`配置，默认配置有以下字段：  

- Name：配置名称，后面运行测试的时候会用到，用于区别多个Kubernetes配置，默认为：kubernetes
- Kubernetes URL：集群访问url，可通过`kubectl cluster-info`查看，如果集群有部署**DNS**插件, 也可以直接填服务名称(自动解析)，默认使用服务名称：https://kubernetes
- Jenkins URL：Jenkins访问地址，默认使用服务名称+端口号

在Jenkins初始化时，默认都已经配置好了，可以直接新建项目测试了。

## 简单测试
点击左边：新建任务——>流水线(Pipeline)
任务名称可以随便起，这里为：k8s-test
配置——>流水线，选择`Pipeline script`
以下为测试脚本内容：
```
podTemplate(label: 'jenkins-slave', cloud: 'kubernetes')
{
    node ('jenkins-slave') {
        stage('test') {
            echo "hello, world"
            sleep 60
        }
    }
}
```

- cloud：插件配置中的Name
- label：插件配置中的Images——>Kubernetes Pod Tempalte——>Labels
- node：与label一致即可

保存配置，点击立即构建，查看控制台输出，出现以下内容就表示运行成功了：
```
Agent default-lsths is provisioned from template Kubernetes Pod Template
Agent specification [Kubernetes Pod Template] (jenkins-slave): 
* [jnlp] jenkins/jnlp-slave:alpine(resourceRequestCpu: 200m, resourceRequestMemory: 256Mi, resourceLimitCpu: 200m, resourceLimitMemory: 256Mi)

Running on default-lsths in /home/jenkins/workspace/k8s-test
[Pipeline] {
[Pipeline] stage
[Pipeline] { (test)
[Pipeline] echo
hello, world
[Pipeline] sleep
Sleeping for 1 min 0 sec
[Pipeline] }
[Pipeline] // stage
[Pipeline] }
[Pipeline] // node
[Pipeline] }
[Pipeline] // podTemplate
[Pipeline] End of Pipeline
Finished: SUCCESS
```


## 配置自动触发CI

- 配置Gitlab项目  
在`Gitlab`中创建一个测试项目，将上面测试的脚本内容写入到一个`Jenkinsfile`文件中，然后上传到该测试项目根路径下。

- 配置Jenkins项目  
点击项目`配置`——>`构建触发器`——>勾选`Build when a change is pushed to GitLab. GitLab webhook URL:http://jenkins.local.com/project/k8s-test`——>保存配置

- 配置Webhook  
进入Gitlab测试项目的`Settings——>Integrations`，一般只需要填写`URL`即可，其他的可根据需求环境配置
默认Jenkins配置不允许匿名用户触发构建，因此还需要添加用户和token。  
URL的格式为：  
`http://[UserID]:[API Token]@jenkins.local.com/project/[ProjectName]`

Jenkins 用户ID Token查看：
点击右上角的`用户名——>设置——>API Token(点击Show API Token...)`

最终Webhook中的URL类似：
http://admin:a910b1492e39e9dd1ea48ea7f7638aaf@jenkins.local.com/project/k8s-test

后面只需要我们一提交代码到Git仓库，就会自动触发Jenkins进行构建了。

## 项目应用
这里我们以一个简单的Java项目为例，实战演示如何进行CI/CD。
基本环境配置上面已经说过了，这里就不多介绍。  
示例项目：https://github.com/lusyoe/springboot-k8s-example

结构说明：
- 镜像构建文件：`Dockerfile`
- k8s应用配置：`k8s-example.yaml`
- 项目源码：`src`
- Jenkins构建文件：`jenkins/Jenkinsfile`

构建流程说明：
- 通过Jenkins kubernetes插件，定义构建过程中所需的3个docker容器：maven、docker、kubectl (这3个容器都在一个pod中)
- 挂载docker.sock和kubeconfig文件
- 首先使用`maven`容器，检出代码，执行项目构建
- 使用`docker`容器，构建镜像，推送到镜像参考
- 使用`kubectl`容器，部署`k8s-example`应用(这里后面也可以使用helm)

访问：  
项目通过Ingress访问`k8s-example.com`，出现`hello, world`,就表示服务部署成功了。
