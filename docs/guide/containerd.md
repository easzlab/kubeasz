# containerd 容器运行时

作为 CNCF 毕业项目，containerd 致力于提供简洁、可靠、可扩展的容器运行时；它被设计用来集成到 kubernetes 等系统使用，而不是像 docker 那样独立使用。

- 安装指南 https://github.com/containerd/cri/blob/master/docs/installation.md
- 客户端 circtl 使用指南 https://github.com/containerd/cri/blob/master/docs/crictl.md
- man 文档 https://github.com/containerd/containerd/tree/master/docs/man

目前 containerd 官方文档还在整理中，但是作为集成在 kubernetes 集群里面使用，阅读以上的文档也就够了。

## kubeasz 集成安装 containerd

- 按照 example 例子，在 ansible hosts 设置全局变量 `CONTAINER_RUNTIME="containerd"`
- 执行 `ansible-playbook 90.setup.yml` 或 `easzctl setup` 即可
