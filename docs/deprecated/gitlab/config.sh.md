## 3.2 环境配置替换 config.sh 

首先应用开发人员需要整理在不同环境（测试环境/生产环境）的配置参数，并在源代码中约定好替换的名称（如db_host, db_usr）；然后用户必须在项目gitlab web界面（“Settings”>"CI/CD">"Variables"）配置变量；最后根据gitlab-ci.yml文件定义CI/CD执行的需要，编写如下简单变量替换shell脚本；该shell脚本分别在测试环境打包阶段（beta-build）和生产环境打包阶段（prod-build）阶段运行。

以下脚本仅作示例，实际应根据项目需要增加/修改需替换变量名称与对应源代码中的配置文件

``` bash
cat > .ci/config.sh << EOF
#!/bin/bash

#set -o verbose
#set -o xtrace

beta_config() {
sed -i \
        -e "s/db_host/$BETA_DB_HOST/g" \
        -e "s/db_usr/$BETA_DB_USR/g" \
        -e "s/db_pwd/$BETA_DB_PWD/g" \
    example-web/src/main/resources/config/datasource.properties        # 项目源码的配置文件
sed -i \
        -e "s/redis_host/$BETA_REDIS_HOST/g" \
        -e "s/redis_port/$BETA_REDIS_PORT/g" \
        -e "s/redis_pwd/$BETA_REDIS_PWD/g" \
    example-web/src/main/resources/config/redis.properties             # 项目源码的配置文件
}

prod_config() {
sed -i \
        -e "s/db_host/$PROD_DB_HOST/g" \
        -e "s/db_usr/$PROD_DB_USR/g" \
        -e "s/db_pwd/$PROD_DB_PWD/g" \
    example-web/src/main/resources/config/datasource.properties
sed -i \
        -e "s/redis_host/$PROD_REDIS_HOST/g" \
        -e "s/redis_port/$PROD_REDIS_PORT/g" \
        -e "s/redis_pwd/$PROD_REDIS_PWD/g" \
    example-web/src/main/resources/config/redis.properties
}

if [[ "$CI_JOB_STAGE" == "beta-build" ]];then
	beta_config
elif [[ "$CI_JOB_STAGE" == "prod-build" ]];then
	prod_config
else
	echo "error: undefined CI_JOB_STAGE!"
fi
EOF
```

