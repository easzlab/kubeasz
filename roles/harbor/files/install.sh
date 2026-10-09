#!/bin/bash
set -o nounset
set -o errexit
#set -o xtrace

HARBOR_HOSTNAME=harbor.easzlab.lo
HARBOR_TLS_PORT=7443
HARBOR_PATH=/var/data
HARBOR_VER=v2.12.4

HARBOR_ADMIN_PASSWORD=$(< /dev/urandom tr -dc 0-9A-Za-z-_ | head -c 16)
HARBOR_DB_PASSWORD=$(< /dev/urandom tr -dc 0-9A-Za-z-_ | head -c 16)

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

BASE=$(cd "$(dirname "$0")"; pwd)
cd "$BASE"
mkdir -p "$BASE/down" "$BASE/certs" "$HARBOR_PATH"

[[ -d "$HARBOR_PATH/harbor" ]] && { logger info "harbor already installed, exit"; exit 0; }

logger debug "install docker-compose"
[[ -f /usr/local/bin/docker-compose ]] || \
cp -f "$BASE/down/docker-compose" /usr/local/bin

logger debug "copy harbor-offline-installer..."
mkdir -p "$HARBOR_PATH"
cp -f "$BASE/down/harbor-offline-installer-$HARBOR_VER.tgz" "$HARBOR_PATH/"

logger debug "extract harbor files..."
[[ -d "$HARBOR_PATH/harbor" ]] || \
{ cd "$HARBOR_PATH"; tar zxf "harbor-offline-installer-$HARBOR_VER.tgz"; }

logger debug "set harbor config: harbor.yml"
cp -f "$BASE/harbor.yml.j2" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s/HARBOR_HOSTNAME/$HARBOR_HOSTNAME/g" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s/HARBOR_TLS_PORT/$HARBOR_TLS_PORT/g" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s%HARBOR_CERTS%$BASE/certs%g" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s/HARBOR_ADMIN_PASSWORD/$HARBOR_ADMIN_PASSWORD/g" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s/HARBOR_DB_PASSWORD/$HARBOR_DB_PASSWORD/g" "$HARBOR_PATH/harbor/harbor.yml"
sed -i "s%HARBOR_PATH%$HARBOR_PATH%g" "$HARBOR_PATH/harbor/harbor.yml"

logger debug "install harbor..."
cd "$HARBOR_PATH/harbor" && ./install.sh
