## kubeasz 3.3.2

kubeasz 3.3.2 发布，小版本更新以及一些bugfix。

### 版本更新
 - k8s: v1.24.9
 - containerd: v1.6.14
 - calico: v3.23.5
 - cilium: v1.12.4
 - dashboard: v2.7.0
 - pause: 3.9
 - harbor: v2.1.5
 - k8s-dns-node-cache: 1.22.13

### 其他

其他更新或者修复，主要以cherry-pick形式从master分支拉取更新。

 - update kubeasz-ext-bin 1.6.4, containerd v1.6.14
 - 更新etcd文档 (#1218)
 - update harbor v2.1.5
 - fix: containerd无法拉取harbor的镜像
 - 检测  ex-lb 的 kube-apiserver 是否正常 (#1215)
 - update cilium 1.12.4
 - update docs:README.md
 - update component
 - minor changes
 - 修改 completion 自动补全 (#1213)
 - 检测 harbor 端口能否连接 (#1211)
 - 修改 nodelocaldns 上游 dns 配置 (#1210)
 - docs update
 - 调整集群内api服务地址和dns服务地址设置方式
 - 自定义 harbor 安装路径 (#1209)
 - 忽略文件不存在 (#1206)
 - update components
 - 加载环境变量, 避免 iptables 命令不存在 (#1203)
 - shell加载环境变量 (#1202)
 - fix: installing offline system packages
 - fix to support recreating CA and certs
 - adjust scripts to support recreating CA and certs
 - modify to run 'cluster-addon' setup on ansible host
 - fix kube-apiserver 访问 kubelet的权限
 - #1183 add 96.update-certs.yml
 - 修改默认配置的 Kubernetes CA 证书 (#1197)
 - New dev (#1193)
 - update container image mirror site #1192
 - 修改ectd 备份命令和备份路径均在 ansible 节点
 -        modified:   roles/calico/templates/calico-v3.23.yaml.j2
 - Create calico-v3.19.yaml.j2
 - 更新 roles/kube-master/main.yml 修改证书时复制新证书到Master节点
 - modify stale issues checking
 - add a bot to manage issues
 - fix issue:#1159
 - fix docs
 - replace 'uname -p(non-portable)' to 'uname -m'
 - update README
 - update docs
 - remove the need of admin kubeconfig on master nodes
 - remove kubectl admin kubeconfig to improve security
 - fix: curl dns resolving problem in a rare case
 - add multi-platform support for json-mock:v1.3.0
 - add support for multi-platform part2
 - feat: add support for multi-platform in "ezdown"
 - fix display images in docs:cilium-example.md
 - to clean some pics
 - fix logo url
 - fix: ca-config.json by libinglong
 - Update main.yml
 - Update ezctl
 - Update 01-CA_and_prerequisite.md
 - update cilium 1.12.2
 - fix flannel v0.19.2
 - update flannel v0.19.2
 - update nodelocaldns 1.22.8
 - fix: cilium 1.12.1 setup
 - update cilium 1.12.1
 - update kube-prometheus-stack-39.11.0
 - update k8s binary & calico version
 - Update config.yml
 - 修复calico ipip隧道模式说明错误，并完善可选参数说明以及使用场景
 - fix custom PATH settings
