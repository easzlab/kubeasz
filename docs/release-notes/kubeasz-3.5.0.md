## kubeasz 3.5.0 (Winter Solstice)

天时人事日相催，冬至阳生春又来。kubeasz 3.5.0 发布，支持k8s v1.26版本，组件更新以及一些bugfix。

### 版本更新

- k8s: v1.26.0
- calico: v3.23.5
- cilium: v1.12.4
- dashboard: v2.7.0
- pause: 3.9
- harbor: v2.1.5
- k8s-dns-node-cache: 1.22.13

### 调整项目分支更新规则

k8s大版本对应kubeasz特定的大版本号，详见README.md 中版本对照表，当前积极更新的分支如下：

- v3.4：对应k8s v1.25 版本，继续保持更新，会使用cherry-pick方式合并后续版本中的重要commit
- v3.5：对应k8s v1.26 版本，持续保持更新
- master：默认保持与最新分支办法同步，当前与v3.5同步

### 其他

- 忽略cloud-init相关文件不存在的错误 (#1206) by itswl
- 自定义 harbor 安装路径 (#1209) by itswl
- 调整集群内api服务地址和dns服务地址设置方式
- 修改 nodelocaldns 上游 dns 配置 (#1210) by itswl
- 增加检测 harbor 端口能否连接 (#1211) by itswl
- 修改 completion 自动补全 (#1213) by itswl
- 增加检测 ex-lb 的 kube-apiserver 是否正常 (#1215) by itswl
- 修复containerd无法拉取harbor的镜像
- 
