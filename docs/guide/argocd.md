# argocd 使用小记

## 背景

使用argocd实现应用基于gitops的持续部署。

企业内部应用都以helm charts方式部署，charts托管在内部git仓库；具体应用配置（helm values）根据不同环境也托管在内部git仓库。所以可以简单理解部署方式如下：

**应用部署=应用chart+应用values**

## 安装

kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

## 登陆

kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'
argocd login ${nodeIp}:${nodePort}
argocd account update-password

## 添加新集群

- 1、配置多个集群的 CONTEXT
export KUBECONFIG=$HOME/.kube/config.1:$HOME/.kube/config.2
kubectl config get-contexts

- 2、添加新集群，根据上面get-contexts结果添加
argocd cluster add 2xxx104401xxxx7-cxxxxxxxxxxxxxa74afabbf --kubeconfig $HOME/.kube/kubeconfig.1 --name test

## 添加项目
kubectl apply -f project.yaml

```
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: myproject
  namespace: argocd
spec:
  clusterResourceWhitelist:
  - group: '*'
    kind: '*'
  description: '测试环境：myproject'
  destinations:
  - name: myproject
    namespace: '*'
    server: https://121.xx.xx.xx:6443
  namespaceResourceWhitelist:
  - group: '*'
    kind: '*'
  # 建议不要开启孤岛资源监控，很可能会引起大量非必要应用同步，造成cpu满载
  #orphanedResources:
  #  warn: false
  sourceRepos:
  - '*'
  sourceNamespaces:
  - '*'
```

## 添加git仓库

UI 界面添加即可

## 添加应用

- 使用git管理的charts仓库：git@172.16.1.1:git-charts.git
- 使用git管理的values仓库：git@172.16.1.1:git-values.git

### 添加单应用
```
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: test-app
  namespace: argocd
spec:
  syncPolicy:
    # 一般建议禁用自动应用同步
    #automated: {}
    syncOptions:
    - CreateNamespace=true
    - ServerSideApply=true
  project: myproject
  destination:
    server: https://121.xx.xx.xx:6443
    namespace: default
  sources:
  - repoURL: 'git@172.16.1.1:git-charts.git'
    targetRevision: master
    path: charts/test-app
    helm:
      valueFiles:
      - values.yaml
      - $values/myproject-test/global.yaml
      - $values/myproject-test/test-app.yaml
  - repoURL: 'git@172.16.1.1:git-values.git'
    targetRevision: master
    ref: values
```

### 添加批次应用
```
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: test-appset
  namespace: argocd
spec:
  generators:
  - git:
      repoURL: 'git@172.16.1.1:git-charts.git'
      revision: master
      directories:
      - path: charts/*
      - path: charts/extras
        exclude: true
  template:
    metadata:
      name: '{{path.basename}}'
    spec:
      project: myproject
      sources:
      - repoURL: 'git@172.16.1.1:git-charts.git'
        targetRevision: master
        path: charts/{{path.basename}}
        helm:
          valueFiles:
          - values.yaml
          - $values/myproject-test/global.yaml
          - $values/myproject-test/{{path.basename}}.yaml
      - repoURL: 'git@172.16.1.1:git-values.git'
        targetRevision: master
        ref: values
      destination:
        server: https://121.xx.xx.xx:6443
        namespace: default
      syncPolicy:
        #automated: {}
        syncOptions:
        - CreateNamespace=true
        - ServerSideApply=true
```

## 其他

- 允许argocd应用在任意命名空间创建
https://argo-cd.readthedocs.io/en/stable/operator-manual/app-any-namespace/
