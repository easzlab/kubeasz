#!/bin/bash

# 把如下脚本放入：/etc/cron.d/ 目录，每分钟自动同步
# echo '* * * * * root /bin/bash "/root/repository/git-sync.sh" -u user -p eZm******AQ > "/root/repository/sync.log" 2>&1'

set -o nounset
set -o errexit

SOURCE_GIT="github.com/mypro1/"
TARGET_GIT="git@192.168.0.2:"

function usage() {
  echo -e "\033[33mUsage:\033[0m git_sync [options] [args]"
  cat <<EOF
  option:
    -u <user>       set user
    -p <token>      set private token
EOF
}

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

# 同步特定分支
function sync_branch(){
    mkdir -p "$BASE_DIR/$1"
    for ((i=0;i<${#PROJECTS[@]};i++))
    do
        PROJECT_NAME=${PROJECTS[i]}
        PRO_GIT_NAME=$PROJECT_NAME.git

        pull_and_push $1 $PROJECT_NAME || pull_and_push_force $1 $PROJECT_NAME
    done
}

# 参数 $1 代表分支名
# 参数 $2 代表项目名
function pull_and_push(){
    [[ -d "$BASE_DIR/$1/$2" ]] || { logger warn "project not existed"; return 1; }
    logger info "normal pull_and_push $2 $1"
    cd "$BASE_DIR/$1/$2"
    git pull origin $1 || { logger error "git pull $2 $1"; return 1; }
    git push secondary $1 || { logger error "git push $2 $1"; return 1; }
}


function pull_and_push_force(){
    logger warn "force pull_and_push $2 $1"
    cd "$BASE_DIR/$1"
    rm -rf "$2"
    git clone -b $1 "$FROM_GIT$PRO_GIT_NAME" || { logger error "git pull $2 $1"; exit 1; }
    cd "$2"
    git remote add secondary "$TARGET_GIT$PRO_GIT_NAME"
    git push secondary $1 --force || { logger error "git push $2 $1"; exit 1; }
}


function main() {
    BASE_DIR=/tmp/mygit_sync

    [[ "$#" -eq 0 ]] && { usage >&2; exit 1; }

    USER=""
    TOKEN=""

    while getopts "u:p:f" OPTION; do
      case "$OPTION" in
        u)
          USER="$OPTARG"
          ;;
        p)
          TOKEN="$OPTARG"
          ;;
        ?)
          usage
          exit 1
          ;;
      esac
    done

    FROM_GIT="http://$USER:$TOKEN@$SOURCE_GIT"

    # 1. 同步项目master、release分支
    PROJECTS=("zscluster" "setup")
    sync_branch master
    sync_branch release

    # 4. done
    logger debug "sync finished"
}

main "$@"
