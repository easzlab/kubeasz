#!/bin/bash
#------------------------------------------------------------
# - Save docker images to disk or load images from disk
#
# @author:  gjmzj
# @usage:   ./imgutil.sh <image_repository>
# @ref:     https://github.com/easzlab/kubeasz/tree/master/tools/imgutils

set -o nounset
set -o errexit
#set -o xtrace

function logger() {
  TIMESTAMP=$(date +'%Y-%m-%d %H:%M:%S')
  local FNAME=$(basename "${BASH_SOURCE[1]}")
  local SOURCE="\033[36m[$FNAME:${BASH_LINENO[0]}]\033[0m"
  case "$1" in
    debug)
      echo -e "\033[36m$TIMESTAMP\033[0m $SOURCE \033[36mDEBUG $2\033[0m"
      ;;
    info)
      echo -e "\033[36m$TIMESTAMP\033[0m $SOURCE \033[32mINFO $2\033[0m"
      ;;
    warn)
      echo -e "\033[36m$TIMESTAMP\033[0m $SOURCE \033[33mWARN $2\033[0m"
      ;;
    error)
      echo -e "\033[36m$TIMESTAMP\033[0m $SOURCE \033[31mERROR $2\033[0m"
      ;;
    *) ;;
  esac
}

function pull_and_push_images(){
  NS="easzlab"
  [ "$#" -eq 1 ] && NS="$1"
  for item in "${IMAGES[@]}"; do
    image_part="${item##*/}"
    image_name="${image_part%:*}"
    image_tag="${image_part##*:}"
    image_file="$imageDir/${image_name}_${image_tag}.tar"
    if [[ ! -f "$image_file" ]];then
      docker pull "$item" && \
      docker save -o "$image_file" "$item" || \
      { logger error "download $item failed!"; return 1; }
    else
      docker load -i "$image_file"
    fi
    docker tag "$item" "easzlab.io.local:5000/${NS}/${image_part}"
    docker push "easzlab.io.local:5000/${NS}/${image_part}" || \
    { logger error "push easzlab.io.local:5000/${NS}/${image_part} failed!"; return 1; }
  done
}

function main() {
  # 检查是否传入参数
  if [ $# -ne 1 ]; then
    echo "Usage: ./imgutil.sh <image_repository>"
    echo "Example: ./imgutil.sh docker.io/library/nginx:alpine"
    exit 1
  fi

  # 可以设置 IMAGE_DIR 环境变量
  imageDir=${IMAGE_DIR:=/etc/kubeasz/down}

  IMAGES=(\
          "$1" \
        )

  pull_and_push_images
}

main "$@"
