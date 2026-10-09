#!/bin/bash

HARBOR_ENDPOINT="harbor.xgeo-pre.cn-hangzhou.lo:7445"
HARBOR_IP="10.10.1.1"
HARBOR_DOMAIN=$(echo "$HARBOR_ENDPOINT"|cut -d':' -f1)

echo "增加/etc/hosts域名解析: $HARBOR_DOMAIN -> $HARBOR_IP"
sed -i "/$HARBOR_DOMAIN/d" /etc/hosts
echo "$HARBOR_IP $HARBOR_DOMAIN" >> /etc/hosts

S3="insecure_skip_verify = true # by_zsdingtalk"
S2="[plugins.\"io.containerd.grpc.v1.cri\".registry.configs.\"$HARBOR_ENDPOINT\".tls] # by_zsdingtalk"
S1="[plugins.\"io.containerd.grpc.v1.cri\".registry.configs] # by_zsdingtalk"

echo "修改containerd配置文件"
sed -i '/by_zsdingtalk/d' /etc/containerd/config.toml
sed -i "/registry\]/a \ \ \ \ \ \ \ \ \ \ $S3" /etc/containerd/config.toml
sed -i "/registry\]/a \ \ \ \ \ \ \ \ $S2" /etc/containerd/config.toml
sed -i "/registry\]/a \ \ \ \ \ \ $S1" /etc/containerd/config.toml

echo "重启containerd服务"
systemctl restart containerd.service
