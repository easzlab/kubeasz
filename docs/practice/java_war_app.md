# JAVA WAR 应用迁移 K8S 实践

初步思路是这样：应用代码与应用配置分离，应用代码打包成 docker 镜像存于内部 harbor 仓库，应用配置使用 configmap 挂载，这样不同的环境只需要修改 configmap 即可部署。

- 使用 maven 把 java 应用代码打包成 xxx.war
- 基于 tomcat 镜像和 xxx.war 做成应用 docker 镜像
- 编写 k8s deployment 文件，在 pod 指定上述应用镜像，同时把应用配置做成 configmap 挂载到 pod 里

经过多次尝试部署发现问题：configmap配置是可以挂载上去，但是会把目录下其他的文件删掉，而且tomcat 目录 webapps/xxxxx/下其他目录也消失了。原来是因为 tomcat 容器完全启动完成后才会解压 war包，而 configmap 配置文件是一开始就挂载上去了，导致失败。

- 调整应用镜像打包过程：xxx.war 先解压后再进行应用镜像打包

## 应用 gitlab CI/CD 集成

- 在内部gitlab创建项目，上传应用java代码，同时在项目根目录下新加如下目录和文件，配置相应的 gitlab-runner 和 环境变量参数

``` bash
├── .app.yaml		# k8s deployment 部署模板文件 
├── config.yaml		# k8s configmap 配置模板文件
├── dockerfiles
│   └── Dockerfile	# Dockerfile 文件
├── .gitlab-ci.yml	# gitlab ci 配置文件
└── .ns.yaml		# k8s namespace 和 imagePullSecrets的配置文件
```
### gitlab-ci 文件摘要

``` bash
variables:
  PROJECT_NS: '$CI_PROJECT_NAMESPACE-$CI_JOB_STAGE'
  APP_NAME: '$CI_PROJECT_NAME-$CI_COMMIT_REF_SLUG'

stages:
  - package
  - beta

job_package:
  stage: package
  tags:
    - package-shell
  only:
    - master
    - /^feature-.*$/
  script:
  - mvn clean install -Dmaven.test.skip=true
  - unzip target/xxxx.war -d dockerfiles/project
  - cd dockerfiles && docker build -t harbor.test.lo/project/$CI_PROJECT_NAME:$CI_PIPELINE_ID .
  - docker login -u $HARBOR_USR -p $HARBOR_PWD harbor.test.lo
  - docker push harbor.test.lo/project/$CI_PROJECT_NAME:$CI_PIPELINE_ID
  - docker logout harbor.test.lo

job_push_beta:
  stage: beta
  tags:
    - beta-shell
  only:
    - master
    - /^feature-.*$/
  when: manual
  script:
  # 替换beta环境的参数配置
  - sed -i "s/PROJECT_NS/$PROJECT_NS/g" config.yaml .app.yaml .ns.yaml
  - sed -i "s/TemplateProject/$APP_NAME/g" config.yaml .app.yaml
  - sed -i "s/DB_HOST/$BETA_DB_HOST/g" config.yaml
  - sed -i "s/DB_PWD/$BETA_DB_PWD/g" config.yaml
  - sed -i "s/APP_REP/$BETA_APP_REP/g" .app.yaml
  - sed -i "s/ProjectImage/$CI_PROJECT_NAME:$CI_PIPELINE_ID/g" .app.yaml
  #
  - mkdir -p /opt/kube/$PROJECT_NS/$APP_NAME
  - cp -f .ns.yaml config.yaml .app.yaml /opt/kube/$PROJECT_NS/$APP_NAME
  - kubectl --kubeconfig=/etc/.beta/config apply -f .ns.yaml
  - kubectl --kubeconfig=/etc/.beta/config apply -f config.yaml
  - kubectl --kubeconfig=/etc/.beta/config apply -f .app.yaml

# 生产部署与beta环境类同，这里省略
```

### Dockerfile 编写

```
FROM tomcat:8.5.33-jre8-alpine

COPY . /usr/local/tomcat/webapps/

# 设置tomcat日志使用的时区
RUN sed -i 's/^JAVA_OPTS=.*webresources\"$/JAVA_OPTS=\"$JAVA_OPTS -Djava.protocol.handler.pkgs=org.apache.catalina.webresources -Duser.timezone=GMT+08\"/g' /usr/local/tomcat/bin/catalina.sh
```

### k8s deployment 配置举例

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: TemplateProject
  namespace: PROJECT_NS
spec:
  replicas: APP_REP
  template:
    metadata:
      labels:
        run: TemplateProject
    spec:
      containers:
      - name: TemplateProject
        image: harbor.test.lo/project/ProjectImage
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: db-config
          mountPath: "/usr/local/tomcat/webapps/project/xxxx/yyyy/config/datasource.properties"
          subPath: datasource.properties
      imagePullSecrets:
      - name: projectkey1
      volumes:
      - name: db-config
        configMap:
          name: TemplateProject-config
          defaultMode: 0640
          items:
          - path: datasource.properties
            key: datasource.properties

---
apiVersion: v1
kind: Service
metadata:
  labels:
    run: TemplateProject
  name: TemplateProject
  namespace: PROJECT_NS
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    run: TemplateProject
  sessionAffinity: None
```

### k8s configmap 配置举例

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: TemplateProject-config
  namespace: PROJECT_NS
data:
  datasource.properties: |
    dataSource.maxIdle = 5
    dataSource.maxActive = 41
    dataSource.driverClassName = com.mysql.jdbc.Driver
    dataSource.url = jdbc:mysql://DB_HOST:8066/project?useUnicode=true&characterEncoding=utf-8
    dataSource.username = username
    dataSource.password = DB_PWD
```
