## kubeasz 3.4.3

kubeasz 3.4.3 发布，小版本更新以及一些bugfix。

### 小版本更新

- k8s: v1.25.5
- containerd: v1.6.14
- calico: v3.23.5
- cilium: v1.12.4
- dashboard: v2.7.0
- pause: 3.9
- harbor: v2.1.5
- k8s-dns-node-cache: 1.22.13

### 其他

其他更新或者修复，主要以cherry-pick形式从master分支拉取更新。

- 忽略cloud-init相关文件不存在的错误 (#1206) by itswl
- 自定义 harbor 安装路径 (#1209) by itswl
- 调整集群内api服务地址和dns服务地址设置方式
- 修改 nodelocaldns 上游 dns 配置 (#1210) by itswl
- 增加检测 harbor 端口能否连接 (#1211) by itswl
- 修改 completion 自动补全 (#1213) by itswl
- 增加检测 ex-lb 的 kube-apiserver 是否正常 (#1215) by itswl
- 修复containerd无法拉取harbor的镜像
