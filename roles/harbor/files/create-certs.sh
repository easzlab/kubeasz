#!/bin/bash
set -o nounset
set -o errexit
#set -o xtrace
CFSSL_VER=v1.6.3

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

function prepare(){
    cat > $BASE/certs/ca-config.json << EOF
{
  "signing": {
    "default": {
      "expiry": "438000h"
    },
    "profiles": {
      "mtls": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "438000h"
      }
    }
  }
}
EOF

    cat > $BASE/certs/ca-csr.json << EOF
{
    "CN": "mtls-ca",
    "key": {
      "algo": "rsa",
      "size": 2048
    },
    "names": [
      {
        "C": "CN",
        "ST": "ZJ",
        "L": "HZ",
        "O": "DT",
        "OU": "EASZLAB"
      }
    ],
    "ca": {
      "expiry": "876000h"
    }
}
EOF

    cat > $BASE/certs/server-csr.json << EOF
{
  "CN": "reg-server",
  "hosts": [
    "registry.easzlab.lo",
    "harbor.easzlab.lo",
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "ZJ",
      "L": "HZ",
      "O": "DT",
      "OU": "EASZLAB"
    }
  ]
}
EOF
}

function create_certs() {
  # download
  logger info "download image: cfssl-utils"
  if [[ ! -f "$BASE/down/cfssl-utils-$CFSSL_VER.tar" ]];then
    docker pull "easzlab/cfssl-utils:$CFSSL_VER" && \
    docker save -o "$BASE/down/cfssl-utils-$CFSSL_VER.tar" "easzlab/cfssl-utils:$CFSSL_VER"
  fi
  logger debug "load image..."
  docker load -i "$BASE/down/cfssl-utils-$CFSSL_VER.tar" > /dev/null

  # clean
  docker ps -a --format="{{ .Names }}"|grep cfssl-utils > /dev/null && \
  logger info "save current certs in backup" && \
  cp -r "$BASE/certs" "$BASE/backup/certs.$(date +'%Y%m%d%H%M%S')" && \
  logger info "stop&remove container: cfssl-utils" && \
  docker rm -f cfssl-utils > /dev/null

  # create
  logger info "start container: cfssl-utils"
  docker run -d \
           --name cfssl-utils \
           --volume "$BASE/certs":/certs/ \
       "easzlab/cfssl-utils:$CFSSL_VER" \
           sleep 3600

  logger info "create ca.pem/ca-key.pem..."
  docker exec -it cfssl-utils sh -c 'cd /certs && cfssl gencert -initca ca-csr.json | cfssljson -bare ca'

  logger info "create server.pem/server-key.pem..."
  docker exec -it cfssl-utils sh -c 'cd /certs && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=mtls server-csr.json | cfssljson -bare server'
}


BASE=$(cd "$(dirname "$0")"; pwd)
cd "$BASE"
mkdir -p "$BASE/down" "$BASE/certs" "$BASE/backup"

prepare
create_certs
