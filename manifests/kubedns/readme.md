### 说明

+ 本目录为k8s集群的插件 kube-dns的配置目录，初始时本目录为空
+ 因kubedns.yaml文件中参数(CLUSTER_DNS_SVC_IP, CLUSTER_DNS_DOMAIN)根据hosts文件设置而定，需要使用ansible template模块替换参数后生成
+ 运行 `ansible-playbook 01.prepare.yml`后该目录下生成kubedns.yaml 文件
+ kubedns.yaml [模板文件](../../roles/deploy/template/kubedns.yaml.j2)
