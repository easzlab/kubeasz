## kubeasz-0.2.1 发布说明

CHANGELOG:
如果服务器能够使用内部yum源/apt源，但是无法访问公网情况下，请下载离线docker镜像完成集群安装：从百度云盘把`basic_images_kubeasz_x.y.tar.gz` 下载解压到项目`down`目录即可
- 组件更新：
  - 更新 coredns版本1.1.3  
- 功能更新：
  - 集成网络插件（可选）使用离线docker镜像安装 
  - 集成其他插件（可选）使用离线docker镜像安装
  - 增加切换集群网络插件的脚本
- 文档更新：
  - [快速指南](https://github.com/gjmzj/kubeasz/blob/master/docs/quickStart.md)
  - [安装规划](https://github.com/gjmzj/kubeasz/blob/master/docs/00-%E9%9B%86%E7%BE%A4%E8%A7%84%E5%88%92%E5%92%8C%E5%9F%BA%E7%A1%80%E5%8F%82%E6%95%B0%E8%AE%BE%E5%AE%9A.md)
  - [切换网络](https://github.com/gjmzj/kubeasz/blob/master/docs/op/clean_k8s_network.md)
- 其他：
  - Bug fix: 清理集群时可能出现`Device or resource busy: '/var/run/docker/netns/xxxxxxx'`的错误，可手动umount后重新清理集群
  - Bug fix: #239 harbor调整安装解压工具, 适配多系统 (#240)
   
