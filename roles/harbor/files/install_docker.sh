#!/bin/bash
set -o nounset
set -o errexit
#set -o xtrace
DOCKER_VER=28.5.2
DOCKER_URL="https://mirrors.aliyun.com/docker-ce/linux/static/stable/x86_64/docker-${DOCKER_VER}.tgz"
DOCKER_ROOT="/data/lib/docker"

function logger() {
  TIMESTAMP=$(date +'%Y-%m-%d %H:%M:%S')
  case "$1" in
    debug)
      echo -e "$TIMESTAMP \033[36mDEBUG\033[0m $2"
      ;;
    info)
      echo -e "$TIMESTAMP \033[32mINFO\033[0m $2"
      ;;
    warn)
      echo -e "$TIMESTAMP \033[33mWARN\033[0m $2"
      ;;
    error)
      echo -e "$TIMESTAMP \033[31mERROR\033[0m $2"
      ;;
    *)
      ;;
  esac
}

function install_docker() {
  mkdir -p /usr/local/bin

  if [[ -f "$BASE/down/docker-${DOCKER_VER}.tgz" ]];then
    logger info "docker binaries already existed"
  else
    logger info "downloading docker binaries: $DOCKER_VER"
    mkdir -p "$BASE/down"
    if [[ -e /usr/bin/wget ]];then
      wget -c --no-check-certificate "$DOCKER_URL" || { logger error "downloading docker failed"; exit 1; }
    else
      curl -k -C- -O --retry 3 "$DOCKER_URL" || { logger error "downloading docker failed"; exit 1; }
    fi
    /bin/mv docker-${DOCKER_VER}.tgz "$BASE/down"
  fi

  # check if a container runtime is already installed
  systemctl status docker|grep Active|grep -q running && { logger warn "docker is already running."; return 0; }

  tar zxf "$BASE/down/docker-${DOCKER_VER}.tgz" && \
  /bin/mv -f docker/* /usr/local/bin
  rm -rf docker

  logger info "generate docker service file"
  cat > /etc/systemd/system/docker.service << EOF
[Unit]
Description=Docker Application Container Engine
Documentation=http://docs.docker.io
[Service]
Environment="PATH=/usr/local/bin:/bin:/sbin:/usr/bin:/usr/sbin"
ExecStartPre=/bin/sh -c 'groupadd docker > /dev/null 2>&1 || echo ""'
ExecStart=/usr/local/bin/dockerd
ExecStartPost=/sbin/iptables -I FORWARD -s 0.0.0.0/0 -j ACCEPT
ExecStartPost=/bin/echo "1" > /proc/sys/net/ipv4/ip_forward
ExecReload=/bin/kill -s HUP \$MAINPID
Restart=on-failure
RestartSec=5
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
Delegate=yes
KillMode=process
[Install]
WantedBy=multi-user.target
EOF

  # configuration for dockerd
  logger info "generate docker config file"
  mkdir -p /etc/docker
  cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://hub1.nat.tf",
    "https://docker.1panel.live",
    "https://hub.rat.dev",
    "https://docker.amingg.com"
  ],
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "80m",
    "max-file": "1"
    },
  "data-root": "$DOCKER_ROOT"
}
EOF

  if [[ -e /etc/centos-release || -e /etc/redhat-release ]]; then
    logger info "turn off selinux in CentOS/Redhat"
    getenforce|grep Disabled || setenforce 0
    sed -i 's/^SELINUX=.*$/SELINUX=disabled/g' /etc/selinux/config
  fi

  logger info "enable and start docker"
  systemctl enable docker > /dev/null 2>&1 || logger warn "docker.service already enabled"
  systemctl daemon-reload && systemctl restart docker && sleep 5
}


BASE=$(cd "$(dirname "$0")"; pwd)
cd "$BASE"
mkdir -p "$BASE/down"

install_docker
