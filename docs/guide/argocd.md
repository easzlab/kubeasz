# argocd 安装

用 GitOps 方式把 Kubernetes 声明式配置“自动、可观测、可回滚”地同步到集群的控制器；它是 Kubernetes 世界里 GitOps 的事实标准。

## 初始安装
- 建议使用helm chart 方式基础安装；后续用声明式方式配置cluster、project、repository等

## 服务暴露
- 建议使用ingress方式
- 备用：kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'

## 密码登录
- 获取初始化密码 `argocd admin initial-password -n argocd`
- 登录 `argocd login {nodeIP}:{nodePort}`
- 更新密码 `argocd account update-password`
- 重置遗忘密码 

```
kubectl -n argocd patch secret argocd-secret -p '{"data": {"admin.password": null, "admin.passwordMtime": null}}'
kubectl -n argocd delete pods -l app.kubernetes.io/name=argocd-server
```

## SSO 登录
- 参考文档：https://help.aliyun.com/zh/ram/obtain-user-information-through-oidc
- 阿里云控制台-RAM访问控制-集成管理-OAuth应用：创建应用 https://ram.console.aliyun.com/applications/create
  - OAuth 协议版本：2.0
  - 应用类型：Web应用
  - 回调地址：填写 https://${argocd-server-domain}/api/dex/callback
  - OAuth 范围：openid(必选), aliuid(可选), profile(可选)

- OAuth应用创建后，准备以下参数
  - "应用 ID" --> dex.config: connectors oidc.config.clientID
  - 创建应用密码 --> dex.config: connectors oidc.config.clientSecret

- 配置argocd-cm

```
  dex.config: |
    connectors:
    - type: oidc
      id: aliyun
      name: aliyun
      config:
        issuer: https://oauth.aliyun.com
        clientID: "406************"
        clientSecret: E8G***************************************************b6
        scopes:
        - profile
        - openid
        - aliuid
        getUserInfo: true
        userIDKey: uid
        userNameKey: uid
        claimMapping:
          preferred_username: name
          email: uid
```

- 配置argocd-rbac-cm

```
data:
  policy.csv: |
    # 设置普通用户app-dev 只读权限
    p, role:app-dev, projects, get, *, allow
    p, role:app-dev, applications, get, *, allow
    p, role:app-dev, logs, get, *, allow
    p, role:app-dev, exec, create, */*, allow

    # 设置测试项目，所有权限
    p, role:app-dev, applications, *, test-project/*, allow
    
    # 阿里云子账号 ID：2***********84
    g, "2***********84", role:admin
    g, "2***********27", role:app-dev

  policy.default: role:''
  scopes: '[name]'
```

## 支持 application in any namespace

- 配置 argocd-cm

```
data:
  # 设置argocd 资源标记方式，使用annotation，禁用labelKey
  # application.instanceLabelKey: argocd.argoproj.io/instance
  application.resourceTrackingMethod: annotation
```

- 配置 argocd-cmd-params-cm

```
data:
  #application.namespaces: app-team-one, app-team-two
  application.namespaces: '*'
  applicationsetcontroller.allowed.scm.providers: '*'
  applicationsetcontroller.namespaces: '*'
```

然后重启 argocd-server 和 argocd-application-controller

## 其他设置

- argocd 部署应用 ingress 资源一直Progressing，参考：https://github.com/argoproj/argo-cd/issues/14607

```
# 修改argocd-cm configmap，重启argocd-application-controller
data:
  resource.customizations: |
    networking.k8s.io/Ingress:
      health.lua: |
        hs = {}
        hs.status = "Healthy"
        hs.message = "Skip health check for Ingress"
        return hs
```
