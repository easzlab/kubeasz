# 使用 traefik 配置 https ingress

本文档已过期，安装最新版本，请参考相关官方文档。

本文档基于 traefik 配置 https ingress 规则，请先阅读[配置基本 ingress](ingress.md)。与基本 ingress-controller 相比，需要额外配置 https tls 证书，主要步骤如下：

## 1.准备 tls 证书

可以使用Let's Encrypt签发的免费证书，这里为了测试方便使用自签证书 (tls.key/tls.crt)，注意CN 配置为 ingress 的域名：

``` bash
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=hello.test.com"
```

## 2.在 kube-system 命名空间创建 secret: traefik-cert，以便后面 traefik-controller 挂载该证书

``` bash
$ kubectl -n kube-system create secret tls traefik-cert --key=tls.key --cert=tls.crt
```

## 3.创建 traefik-controller，增加 traefik.toml 配置文件及https 端口暴露等，详见该 yaml 文件

``` bash
$ kubectl apply -f /etc/kubeasz/manifests/ingress/traefik/tls/traefik-controller.yaml
```

## 4.创建 https ingress 例子

``` bash
# 创建示例应用
$ kubectl run test-hello --image=nginx:alpine --port=80 --expose
# hello-tls-ingress 示例
apiVersion: networking.k8s.io/v1beta1 
kind: Ingress
metadata:
  name: hello-tls-ingress
  annotations:
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  - host: hello.test.com
    http:
      paths:
      - backend:
          serviceName: test-hello
          servicePort: 80
  tls:
  - secretName: traefik-cert
# 创建https ingress
$ kubectl apply -f /etc/kubeasz/manifests/ingress/traefik/tls/hello-tls.ing.yaml
# 注意根据hello示例，需要在default命名空间创建对应的secret: traefik-cert
$ kubectl create secret tls traefik-cert --key=tls.key --cert=tls.crt
```

## 5.验证 https 访问

验证 traefik-ingress svc

``` bash
$ kubectl get svc -n kube-system traefik-ingress-service 
NAME                      TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)                                     AGE
traefik-ingress-service   NodePort   10.68.250.253   <none>        80:23456/TCP,443:23457/TCP,8080:35941/TCP   66m
```

可以看到项目默认使用nodePort 23456暴露traefik 80端口，nodePort 23457暴露 traefik 443端口，因此在客户端 hosts 增加记录 `$Node_IP hello.test.com`之后，可以在浏览器验证访问如下：

``` bash
https://hello.test.com:23457
```

如果你已经配置了[转发 ingress nodePort](../op/loadballance_ingress_nodeport.md)，那么增加对应 hosts记录后，可以验证访问 `https://hello.test.com`

## 配置 dashboard ingress

前提1：k8s 集群的dashboard 已安装

```
$ kubectl get svc -n kube-system | grep dashboard
kubernetes-dashboard      NodePort    10.68.211.168   <none>        443:39308/TCP	3d11h
```
前提2：`/etc/kubeasz/manifests/ingress/traefik/tls/traefik-controller.yaml`的配置文件`traefik.toml`开启了`insecureSkipVerify = true`

配置 dashboard ingress：`kubectl apply -f /etc/kubeasz/manifests/ingress/traefik/tls/k8s-dashboard.ing.yaml` 内容如下：

```
apiVersion: networking.k8s.io/v1beta1 
kind: Ingress
metadata:
  name:  kubernetes-dashboard
  namespace: kube-system
  annotations:
    traefik.ingress.kubernetes.io/redirect-entry-point: https
spec:
  rules:
  - host: dashboard.test.com
    http:
      paths:
      - path: /
        backend:
          serviceName: kubernetes-dashboard
          servicePort: 443
```
- 注意annotations 配置了 http 跳转 https 功能
- 注意后端服务是443端口

## 参考

- [Add a TLS Certificate to the Ingress](https://docs.traefik.io/user-guide/kubernetes/#add-a-tls-certificate-to-the-ingress)
