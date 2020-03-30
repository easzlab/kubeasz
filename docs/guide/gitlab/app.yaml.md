## 3.3 K8S 应用部署模板 app.yaml

以下示例配置仅做参考，描述一个简单 java spring boot项目的 k8s 部署文件模板；在实际部署前，CI/CD流程中会对变量做替换。详见 [gitlab-ci.yml文件](gitlab-ci.yml.md)。

``` bash
cat > .ci/app.yaml << EOF
---
apiVersion: v1
kind: Namespace
metadata:
  name: PROJECT_NS
---
apiVersion: v1
kind: Secret
metadata:
  name: harborkey1
  namespace: PROJECT_NS
data:
    #待替换的变量DOCKER_KEY，参考 docs/guide/harbor.md#k8s%E4%B8%AD%E4%BD%BF%E7%94%A8harbor
    .dockerconfigjson: DOCKER_KEY
type: kubernetes.io/dockerconfigjson

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: APP_NAME
  namespace: PROJECT_NS
spec:
  replicas: APP_REP
  template:
    metadata:
      labels:
        run: APP_NAME
    spec:
      containers:
      - name: APP_NAME
        image: ProjectImage
        env:
          # 设置java的时区
          - name: TZ
            value: "Asia/Shanghai"
        resources:
          limits:
            cpu: 500m
            memory: 1600Mi
          requests:
            cpu: 200m
            memory: 800Mi
        ports:
        - containerPort: 8080
      imagePullSecrets:
      - name: harborkey1

---
apiVersion: v1
kind: Service
metadata:
  labels:
    run: APP_NAME
  name: APP_NAME
  namespace: PROJECT_NS
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    run: APP_NAME
  sessionAffinity: None

---
apiVersion: networking.k8s.io/v1beta1 
kind: Ingress
metadata:
  name: APP_NAME-ingress
  namespace: PROJECT_NS
spec:
  rules:
  - host: AppDomain
    http:
      paths:
      - path: /AppPath
        backend:
          serviceName: APP_NAME
          servicePort: 80
EOF
```

