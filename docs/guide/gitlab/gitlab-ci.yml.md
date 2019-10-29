## 3.1 配置 gitlab-ci.yml

示例应用搭建 CI/CD 流水线的背景需求

- 应用测试环境部署在本地k8s平台，生产环境部署在阿里云上k8s平台
- 应用的多个feature分支可以并行测试
- 对于即将发布的release分支，本地提供封版测试环境，阿里云上提供UAT测试环境

以下示例配置为个人经验总结，仅供参考，可以根据自己的理解和项目需要不断优化完善；总体来说 gitlab-ci.yml 配置很丰富，基本上能够满足各种个性化的CI/CD流程需要。

``` bash
$ cat > .ci/gitlab-ci.yml << EOF
variables:                                                               ### 定义全局变量 http://gitlab.test.com/help/ci/variables/README.md
  PROJECT_NS: '$CI_PROJECT_NAMESPACE-$CI_JOB_STAGE'                      # 定义项目命名空间，对应k8s的namespace
  APP_NAME: '$CI_PROJECT_NAME-$CI_COMMIT_REF_SLUG'                       # 使用项目名和git提交信息作为应用名
  IMAGE_NAME: '$CI_PROJECT_NAMESPACE-$CI_PROJECT_NAME:$CI_PIPELINE_ID'   # 定义镜像名称 

stages:                                                                  ### 定义ci各阶段
  - beta-build                                                           # beta环境编译打包
  - beta-deploy                                                          # beta环境部署
  - beta-feature-delete                                                  # beta环境feature分支手动删除
  - prod-build                                                           # prod环境编译打包
  - prod-uat-deploy                                                      # prod-uat环境部署
  - prod-deploy                                                          # prod环境部署
  - prod-rollback                                                        # prod回滚

job_beta_build:
  stage: beta-build                                                      # beta环境编译打包
  tags:
    - build-shell                                                        # 定义带`build-shell`标签的runner可以运行该job
  only:                                                                  # 定义只在如下分支或者tag运行该job 
    - master
    - develop
    - /^feature.*$/
    - release
  #when: manual                                                          # 调试阶段可以先手动，后续可以注释掉以自动运行
  script:                                                                ### runner上运行的脚本
  - bash .ci/config.sh                                                   # 不同环境配置替换，后文详解 config.sh
  - mvn clean install -Dmaven.test.skip=true -U                          # mvn 编译，可以去runner 虚机上手动执行编译测试
  - mv example-web/target/*.jar dockerfiles/                             # 把mvn生成的xxx.jar移动到dockerfiles目录下
  - export IMAGE=`echo $IMAGE_NAME | sed 's/\//-/g'`                     # 转换镜像名，例：mygroup/java/example:172 >> mygroup-java-example:172
  - cd dockerfiles && docker build -t $BETA_HARBOR/example/$IMAGE .      # 创建 docker 镜像
  - docker login -u $BETA_HARBOR_USR -p $BETA_HARBOR_PWD $BETA_HARBOR    # 登录到内部镜像仓库 harbor，并推送
  - docker push $BETA_HARBOR/example/$IMAGE                                          
  - docker logout $BETA_HARBOR

job_push_beta:                                                           ### 推送到beta环境，可以推送不同分支 develop, feature-1, ...>
  stage: beta-deploy                                                     # 可以做到多分支同时测试，甚至最后的release分支也要在beta封版测试
  tags:
    - beta-shell                                                         # 定义带`beta-shell`标签的runner可以运行该job
  only:
    - master
    - develop
    - /^feature.*$/
    - release
  when: manual                                                           # 调试阶段可以先手动，后续可以注释掉以自动运行
  variables:
    BETA_EXP_Domain: '$CI_COMMIT_REF_SLUG.example.test.com'              # job内部变量，指定该应用在beta环境的 ingress 域名
  script:
  - export IMAGE=`echo $IMAGE_NAME | sed 's/\//-/g'`                     # 转换 $IMAGE_NAME 中可能的 / 字符
  - export PROJECT_NS=`echo $PROJECT_NS | sed 's/\//-/g'`                # 转换命名空间中可能有的 / 字符
  # 替换beta环境的参数配置
  - sed -i "s/PROJECT_NS/$PROJECT_NS/g" .ci/app.yaml                     ### app.yaml 即k8s的部署模板文件，详见后面 app.yaml.md 文档，注意这里的变量有的来自>
  - sed -i "s/APP_NAME/$APP_NAME/g" .ci/app.yaml                         # gitlab 系统变量, 有的是在项目 CI/CD 设置里面用户定义的变量
  - sed -i "s/APP_REP/$BETA_APP_REP/g" .ci/app.yaml
  - sed -i "s/AppDomain/$BETA_EXP_Domain/g" .ci/app.yaml
  - sed -i "s/ProjectImage/$BETA_HARBOR\/example\/$IMAGE/g" .ci/app.yaml
  - sed -i "s/DOCKER_KEY/$BETA_KEY/g" .ci/app.yaml                       # DOCKER_KEY 为k8s平台能从镜像仓库pull所需的认证信息，详见harbor文档
  #
  - mkdir -p /opt/kube/$PROJECT_NS/$APP_NAME                             # 在runner：beta-shell虚机本地创建应用配置目录，调试检查用
  - cp -f .ci/app.yaml /opt/kube/$PROJECT_NS/$APP_NAME
  - kubectl --kubeconfig=/etc/.beta/config apply -f .ci/app.yaml         # 部署应用（runner虚机上预先配置了kubectl权限执行测试k8s平台）

job_delete_beta:                                                         ### 多测试环境并行部署在beta k8s平台，feature分支测试完毕后删除代码分支，
  stage: beta-feature-delete                                             # 同时需要删除该分支在k8s平台上的部署，可以由开发人员自行执行该job删除
  tags:
    - beta-shell
  only:
    - /^feature.*$/
  when: manual
  script:
  - export PROJECT_NS=`echo $PROJECT_NS | sed 's/\//-/g'`
  - kubectl --kubeconfig=/etc/.beta/config delete deploy,svc,ing $APP_NAME -n $PROJECT_NS

job_prod_build:                                                          ### prod环境编译打包，这里prod环境我们使用阿里云上的K8S
  stage: prod-build                                                      # 阿里云k8s平台上运行的uat环境和正式环境都使用本次打包镜像
  tags:
    - build-shell
  only:                                                                  # 仅master和release分支可以执行该job
    - master
    - release
  #when: manual
  script:
  - bash .ci/config.sh                                                   # config.sh 会执行替换生产环境的变量
  - mvn clean install -Dmaven.test.skip=true -U                          # mvn 编译，可以去runner 虚机上手动执行编译测试
  - mv example-web/target/*.jar dockerfiles/                             # 把mvn生成的xxx.jar移动到dockerfiles目录下
  - export IMAGE=`echo $IMAGE_NAME | sed 's/\//-/g'`
  - cd dockerfiles && docker build -t $PROD_HARBOR/example/$IMAGE .
  - docker login -u $PROD_HARBOR_USR -p $PROD_HARBOR_PWD $PROD_HARBOR
  - docker push $PROD_HARBOR/example/$IMAGE
  - docker logout $PROD_HARBOR

job_push_prod_uat:                                                       ### 部署至阿里云uat环境
  stage: prod-uat-deploy
  tags:
    - prod-shell
  when: manual
  only:                                                                  # 仅master和release分支可以执行该job
    - master
    - release
  variables:
    PROD_EXP_Domain: 'example-uat.xxxx.com'                              # job内部变量，指定该应用在uat环境的 ingress 域名
  script:
  - export IMAGE=`echo $IMAGE_NAME | sed 's/\//-/g'`
  - export PROJECT_NS=`echo $PROJECT_NS | sed 's/\//-/g'`
  # 替换prod环境的参数配置
  - sed -i "s/PROJECT_NS/$PROJECT_NS/g" .ci/app.yaml
  - sed -i "s/APP_NAME/$CI_PROJECT_NAME/g" .ci/app.yaml
  - sed -i "s/APP_REP/1/g" .ci/app.yaml
  - sed -i "s/AppDomain/$PROD_EXP_Domain/g" .ci/app.yaml
  - sed -i "s/ProjectImage/$PROD_HARBOR\/example\/$IMAGE/g" .ci/app.yaml
  - sed -i "s/DOCKER_KEY/$PROD_KEY/g" .ci/app.yaml
  #
  - mkdir -p /opt/kube/$PROJECT_NS/$APP_NAME
  - cp -f .ci/app.yaml /opt/kube/$PROJECT_NS/$APP_NAME
  - kubectl --kubeconfig=/etc/.aliyun/config apply -f .ci/app.yaml

job_push_prod_release:                                                   ### 部署至阿里云正式环境
  stage: prod-deploy
  tags:
    - prod-shell
  when: manual
  only:                                                                  # 仅master和release分支可以执行该job
    - master
    - release
  variables:
    PROD_EXP_Domain: 'example.xxxx.com'                                  # 指定该应用在阿里云正式环境的 ingress 域名
  script:
  - export IMAGE=`echo $IMAGE_NAME | sed 's/\//-/g'`
  - export PROJECT_NS=`echo $PROJECT_NS | sed 's/\//-/g'`
  # 替换prod环境的参数配置
  - sed -i "s/PROJECT_NS/$PROJECT_NS/g" .ci/app.yaml
  - sed -i "s/APP_NAME/$CI_PROJECT_NAME/g" .ci/app.yaml
  - sed -i "s/APP_REP/$PROD_APP_REP/g" .ci/app.yaml
  - sed -i "s/AppDomain/$PROD_EXP_HOST/g" .ci/app.yaml
  - sed -i "s/ProjectImage/$PROD_HARBOR\/example\/$IMAGE/g" .ci/app.yaml
  - sed -i "s/DOCKER_KEY/$PROD_KEY/g" .ci/app.yaml
  #
  - mkdir -p /opt/kube/$PROJECT_NS/$APP_NAME
  - cp -f .ci/app.yaml /opt/kube/$PROJECT_NS/$APP_NAME
  - kubectl --kubeconfig=/etc/.aliyun/config apply -f .ci/app.yaml

1/3 rollback:                                                            ### 定义生产环境回退job  
  stage: prod-rollback
  tags:
    - prod-shell
  when: manual
  only:
    - master
    - /^release.*$/
  variables:
    PROJECT_NS: '$CI_PROJECT_NAMESPACE-prod-deploy'                      # 定义job内变量覆盖全局变量设置
  script:
  - kubectl --kubeconfig=/etc/.aliyun/config -n $PROJECT_NS rollout undo deployment $CI_PROJECT_NAME --to-revision=1

2/3 rollback:
  stage: prod-rollback
  tags:
    - prod-shell
  when: manual
  only:
    - master
    - /^release.*$/
  variables:
    PROJECT_NS: '$CI_PROJECT_NAMESPACE-prod-deploy'                      # 定义job内变量覆盖全局变量设置
  script:
  - kubectl --kubeconfig=/etc/.aliyun/config -n $PROJECT_NS rollout undo deployment $CI_PROJECT_NAME --to-revision=2

3/3 rollback:
  stage: prod-rollback
  tags:
    - prod-shell
  when: manual
  only:
    - master
    - /^release.*$/
  variables:
    PROJECT_NS: '$CI_PROJECT_NAMESPACE-prod-deploy'                      # 定义job内变量覆盖全局变量设置
  script:
  - kubectl --kubeconfig=/etc/.aliyun/config -n $PROJECT_NS rollout undo deployment $CI_PROJECT_NAME --to-revision=3
EOF
```

恭喜终于看完 gitlab-ci.yml 文件，怎么样，是不是一千个人可以写出一万个 CI/CD 流程 :)
