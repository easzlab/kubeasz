#!/bin/bash
#  Create & manage k8s clusters
#  shellcheck disable=SC2155

set -o nounset
set -o errexit
#set -o xtrace

function usage() {
    echo -e "\033[33mUsage:\033[0m ezctl COMMAND [args]"
    cat <<EOF
-------------------------------------------------------------------------------------
Cluster setups:
    list		             to list all of the managed clusters
    checkout    <cluster>            to switch default kubeconfig of the cluster
    new         <cluster>            to start a new k8s deploy with name 'cluster'
    setup       <cluster>  <step>    to setup a cluster, also supporting a step-by-step way
    start       <cluster>            to start all of the k8s services stopped by 'ezctl stop'
    stop        <cluster>            to stop all of the k8s services temporarily
    upgrade     <cluster>            to upgrade the k8s cluster
    destroy     <cluster>            to destroy the k8s cluster
    backup      <cluster>            to backup the cluster state (etcd snapshot)
    restore     <cluster>            to restore the cluster state from backups
    start-aio		             to quickly setup an all-in-one cluster with default settings

Cluster ops:
    add-etcd    <cluster>  <ip>      to add a etcd-node to the etcd cluster
    add-master  <cluster>  <ip>      to add a master node to the k8s cluster
    add-node    <cluster>  <ip>      to add a work node to the k8s cluster
    del-etcd    <cluster>  <ip>      to delete a etcd-node from the etcd cluster
    del-master  <cluster>  <ip>      to delete a master node from the k8s cluster
    del-node    <cluster>  <ip>      to delete a work node from the k8s cluster

Extra operation:
    kca-renew   <cluster>            to force renew CA certs and all the other certs (with caution)
    kcfg-adm    <cluster>  <args>    to manage client kubeconfig of the k8s cluster

Use "ezctl help <command>" for more information about a given command.
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

function help-info() {
    case "$1" in
        (setup)
	    usage-setup
            ;;
        (add-etcd)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-etcd.md'"
            ;;
        (add-master)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-master.md'"
            ;;
        (add-node)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-node.md'"
            ;;
        (del-etcd)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-etcd.md'"
            ;;
        (del-master)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-master.md'"
            ;;
        (del-node)
            echo -e "read more > 'https://github.com/easzlab/kubeasz/blob/master/docs/op/op-node.md'"
            ;;
        (kca-renew)
	    echo -e "WARNNING: this command should be used with caution"
	    echo -e "force to recreate CA certs and all of the others certs used in the cluster"
	    echo -e "it should be used only when the admin.conf leaked"
            ;;
        (kcfg-adm)
	    usage-kcfg-adm
            ;;
        (*)
            echo -e "todo: help info $1"
            ;;
    esac
}

function usage-kcfg-adm(){
  echo -e "\033[33mUsage:\033[0m ezctl kcfg-adm <cluster> <args>"
  cat <<EOF
available <args>:
    -A     to add a client kubeconfig with a newly created user
    -D     to delete a client kubeconfig with the existed user
    -L     to list all of the users
    -e     to set expiry of the user certs in hours (ex. 24h, 8h, 240h)
    -t     to set a user-type (admin or view)
    -u     to set a user-name prefix

examples: ./ezctl kcfg-adm test-k8s -L
          ./ezctl kcfg-adm default -A -e 240h -t admin -u jack
          ./ezctl kcfg-adm default -D -u jim-202101162141
EOF
}

function usage-setup(){
  echo -e "\033[33mUsage:\033[0m ezctl setup <cluster> <step>"
  cat <<EOF
available steps:
    01  prepare            to prepare CA/certs & kubeconfig & other system settings
    02  etcd               to setup the etcd cluster
    03  container-runtime  to setup the container runtime(docker or containerd)
    04  kube-master        to setup the master nodes
    05  kube-node          to setup the worker nodes
    06  network            to setup the network plugin
    07  cluster-addon      to setup other useful plugins
    90  all                to run 01~07 all at once
    10  ex-lb              to install external loadbalance for accessing k8s from outside
    11  harbor             to install a new harbor server or to integrate with an existed one

examples: ./ezctl setup test-k8s 01  (or ./ezctl setup test-k8s prepare)
	  ./ezctl setup test-k8s 02  (or ./ezctl setup test-k8s etcd)
          ./ezctl setup test-k8s all
          ./ezctl setup test-k8s 04 -t restart_master
EOF
}

### Cluster setups functions ##############################

function new() {
    # check if already existed
    [[ -d "clusters/$1" ]] && { logger error "cluster:$1 already existed, if cluster:$1 setup failed, try 'rm -rf $BASE/clusters/$1' first!"; exit 1; }

    logger debug "generate custom cluster files in $BASE/clusters/$1"
    mkdir -p "clusters/$1"
    cp example/hosts.multi-node "clusters/$1/hosts"
    sed -i "s/_cluster_name_/$1/g" "clusters/$1/hosts"
    cp example/config.yml "clusters/$1/config.yml"

    logger debug "set versions"
    # these variables were imported with eval and have the same name in ezdown
    eval $(sed '/V.[rR]=.*\./!d' ezdown)
    k8sVer=$(echo "$K8S_BIN_VER"|sed 's/v//g')
    registryMirror=true

    grep registry-mirrors /etc/docker/daemon.json > /dev/null 2>&1 || { logger debug "disable registry mirrors"; registryMirror=false; }

    sed -i -e "s/__k8s_ver__/$k8sVer/g" \
       -e "s/__flannel__/$flannelVer/g" \
	   -e "s/__calico__/$calicoVer/g" \
	   -e "s/__cilium__/$ciliumVer/g" \
	   -e "s/__kube_ovn__/$kubeOvnVer/g" \
	   -e "s/__kube_router__/$kubeRouterVer/g" \
	   -e "s/__coredns__/$corednsVer/g" \
	   -e "s/__pause__/$pauseVer/g" \
	   -e "s/__dns_node_cache__/$dnsNodeCacheVer/g" \
	   -e "s/__dashboard__/$dashboardVer/g" \
	   -e "s/__local_path_provisioner__/$localpathProvisionerVer/g" \
	   -e "s/__nfs_provisioner__/$nfsProvisionerVer/g" \
	   -e "s/__prom_chart__/$promChartVer/g" \
	   -e "s/__minio_chart__/$minioOperatorVer/g" \
	   -e "s/__kubeapps_chart__/$kubeappsVer/g" \
	   -e "s/__kubeblocks_ver__/$kubeblocksVer/g" \
	   -e "s/__harbor__/$HARBOR_VER/g" \
	   -e "s/^ENABLE_MIRROR_REGISTRY.*$/ENABLE_MIRROR_REGISTRY: $registryMirror/g" \
	   -e "s/__metrics__/$metricsVer/g" "clusters/$1/config.yml"


    logger debug "cluster $1: files successfully created."
    logger info "next steps 1: to config '$BASE/clusters/$1/hosts'"
    logger info "next steps 2: to config '$BASE/clusters/$1/config.yml'"
}

function setup() {
    [[ -d "clusters/$1" ]] || { logger error "invalid config, run 'ezctl new $1' first"; return 1; }
    [[ -f "bin/kube-apiserver" ]] || { logger error "no binaries founded, run 'ezdown -D' fist"; return 1; }

    # for extending usage
    EXTRA_ARGS=$(echo "$*"|sed "s/$1 $2//g"|sed "s/^ *//g")

    PLAY_BOOK="dummy.yml"
    case "$2" in
      (01|prepare)
          PLAY_BOOK="01.prepare.yml"
          ;;
      (02|etcd)
          PLAY_BOOK="02.etcd.yml"
          ;;
      (03|container-runtime)
          PLAY_BOOK="03.runtime.yml"
          ;;
      (04|kube-master)
          PLAY_BOOK="04.kube-master.yml"
          ;;
      (05|kube-node)
          PLAY_BOOK="05.kube-node.yml"
          ;;
      (06|network)
          PLAY_BOOK="06.network.yml"
          ;;
      (07|cluster-addon)
          PLAY_BOOK="07.cluster-addon.yml"
          ;;
      (90|all)
          PLAY_BOOK="90.setup.yml"
          ;;
      (10|ex-lb)
          PLAY_BOOK="10.ex-lb.yml"
          ;;
      (11|harbor)
          PLAY_BOOK="11.harbor.yml"
          ;;
      (*)
          usage-setup
          exit 1
          ;;
    esac

    COMMAND="ansible-playbook -i clusters/$1/hosts -e @clusters/$1/config.yml $EXTRA_ARGS playbooks/$PLAY_BOOK"
    echo "$COMMAND"

    k8s_ver=$(bin/kube-apiserver --version|cut -d' ' -f2)
    etcd_ver=v$(bin/etcd --version|grep 'etcd Version'|cut -d' ' -f3)
    network_cni=$(grep CLUSTER_NETWORK "clusters/$1/hosts"|cut -d'"' -f2|sed 's/-//g')
    network_cni_ver=$(grep -i "${network_cni}Ver" ezdown|cut -d'=' -f2|head -n1)

    cat <<EOF
*** Component Version *********************
*******************************************
*   kubernetes: ${k8s_ver}
*   etcd: ${etcd_ver}
*   ${network_cni}: ${network_cni_ver}
*******************************************
EOF

    logger info "cluster:$1 setup step:$2 begins in 5s, press any key to abort:\n"
    ! (read -r -t5 -n1) || { logger warn "setup abort"; return 1; }

    ${COMMAND} || return 1
}

function cmd() {
    [[ -d "clusters/$1" ]] || { logger error "invalid config, run 'ezctl new $1' first"; return 1; }

    PLAY_BOOK="dummy.yml"
    case "$2" in
      (start)
          PLAY_BOOK="91.start.yml"
          ;;
      (stop)
          PLAY_BOOK="92.stop.yml"
          ;;
      (upgrade)
          PLAY_BOOK="93.upgrade.yml"
          ;;
      (backup)
          PLAY_BOOK="94.backup.yml"
          ;;
      (restore)
          PLAY_BOOK="95.restore.yml"
          ;;
      (destroy)
          PLAY_BOOK="99.clean.yml"
          ;;
      (*)
          usage
          exit 1
          ;;
    esac

    COMMAND="ansible-playbook -i clusters/$1/hosts -e @clusters/$1/config.yml playbooks/$PLAY_BOOK"
    echo "$COMMAND"

    logger info "cluster:$1 $2 begins in 5s, press any key to abort:\n"
    ! (read -r -t5 -n1) || { logger warn "$2 abort"; return 1; }

    ${COMMAND} || return 1
}


function list() {
    [[ -d ./clusters ]] || { logger error "cluster not found, run 'ezctl new' first"; return 1; }
    [[ -f ~/.kube/config ]] || { logger error "kubeconfig not found, run 'ezctl setup' first"; return 1; }
    which md5sum > /dev/null 2>&1 || { logger error "md5sum not found"; return 1; }

    CLUSTERS=$(cd clusters && echo -- *)
    CFG_MD5=$(sed '/server/d' ~/.kube/config|md5sum|cut -d' ' -f1)
    cd "$BASE"

    logger info "list of managed clusters:"
    i=1; for c in $CLUSTERS;
    do
        if [[ -f "clusters/$c/kubectl.kubeconfig" ]];then
            c_md5=$(sed '/server/d' "clusters/$c/kubectl.kubeconfig"|md5sum|cut -d' ' -f1)
            if [[ "$c_md5" = "$CFG_MD5" ]];then
                echo -e "==> cluster $i:\t$c (\033[32mcurrent\033[0m)"
            else
                echo -e "==> cluster $i:\t$c"
            fi
            ((i++))
        fi
    done
}


function checkout() {
    [[ -d "clusters/$1" ]] || { logger error "invalid config, run 'ezctl new $1' first"; return 1; }
    [[ -f "clusters/$1/kubectl.kubeconfig" ]] || { logger error "invalid kubeconfig, run 'ezctl setup $1' first"; return 1; }
    logger info "set default kubeconfig: cluster $1 (\033[32mcurrent\033[0m)"
    /bin/cp -f "clusters/$1/kubectl.kubeconfig" ~/.kube/config
}

### in-cluster operation functions ##############################

function add-node() {
    # check new node's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error "Invalid ip add:$2"; return 1; }

    # check if the new node already exsited
    sed -n '/^\[kube_master/,/^\[harbor/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " && { logger error "node $2 already existed in $BASE/clusters/$1/hosts"; return 2; }

    logger info "add $2 into 'kube_node' group"
    NODE_INFO="${@:2}"
    sed -i "/\[kube_node/a $NODE_INFO" "$BASE/clusters/$1/hosts"

    logger info "start to add a work node:$2 into cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/22.addnode.yml" -e "NODE_TO_ADD=$2" -e "@clusters/$1/config.yml"
}

function add-master() {
    # check new master's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error  "Invalid ip add:$2"; return 1; }

    # check if the new master already exsited
    sed -n '/^\[kube_master/,/^\[kube_node/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " && { logger error "master $2 already existed!"; return 2; }

    logger info "add $2 into 'kube_master' group"
    MASTER_INFO="${@:2}"
    sed -i "/\[kube_master/a $MASTER_INFO" "$BASE/clusters/$1/hosts"

    logger info "start to add a master node:$2 into cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/23.addmaster.yml" -e "NODE_TO_ADD=$2" -e "@clusters/$1/config.yml"

    logger info "re-setting /etc/hosts for all nodes"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/90.setup.yml" -t set_hosts -e "@clusters/$1/config.yml"

    logger info "reconfigure and restart 'kube-lb' service"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/90.setup.yml" -t restart_kube-lb -e "@clusters/$1/config.yml"

    logger info "reconfigure and restart 'ex-lb' service"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/10.ex-lb.yml" -t restart_lb -e "@clusters/$1/config.yml"
}

function add-etcd() {
    # check new node's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error  "Invalid ip add:$2"; return 1; }

    # check if the new node already exsited
    sed -n '/^\[etcd/,/^\[kube_master/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " && { logger error "etcd $2 already existed!"; return 2; }

    logger info "add $2 into 'etcd' group"
    ETCD_INFO="${@:2}"
    sed -i "/\[etcd/a $ETCD_INFO" "$BASE/clusters/$1/hosts"

    logger info "start to add a etcd node:$2 into cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/21.addetcd.yml" -e "NODE_TO_ADD=$2" -e "@clusters/$1/config.yml"

    logger info "reconfig &restart the etcd cluster"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/02.etcd.yml" -t restart_etcd -e "@clusters/$1/config.yml"

    logger info "restart apiservers to use the new etcd cluster"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/04.kube-master.yml" -t restart_master -e "@clusters/$1/config.yml"
}

function del-etcd() {
    # check node's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error  "Invalid ip add:$2"; return 1; }

    # check if the deleting node exsited
    sed -n '/^\[etcd/,/^\[kube_master/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " || { logger error "etcd $2 not existed!"; return 2; }

    logger warn "start to delete the etcd node:$2 from cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/31.deletcd.yml" -e "ETCD_TO_DEL=$2" -e "CLUSTER=$1" -e "@clusters/$1/config.yml"

    logger info "reconfig &restart the etcd cluster"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/02.etcd.yml" -t restart_etcd -e "@clusters/$1/config.yml"

    logger info "restart apiservers to use the new etcd cluster"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/04.kube-master.yml" -t restart_master -e "@clusters/$1/config.yml"
}

function del-node() {
    # check node's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error "Invalid ip add:$2"; return 2; }

    # check if the deleting node exsited
    sed -n '/^\[kube_master/,/^\[harbor/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " || { logger error "node $2 not existed in $BASE/clusters/$1/hosts"; return 2; }

    logger warn "start to delete the node:$2 from cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/32.delnode.yml" -e "NODE_TO_DEL=$2" -e "CLUSTER=$1" -e "@clusters/$1/config.yml"
}

function del-master() {
    # check node's address regexp
    [[ $2 =~ ^(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})(\.(2(5[0-5]{1}|[0-4][0-9]{1})|[0-1]?[0-9]{1,2})){3}$ ]] || { logger error "Invalid ip add:$2"; return 2; }

    # check if the deleting master exsited
    sed -n '/^\[kube_master/,/^\[kube_node/p' "$BASE/clusters/$1/hosts"|grep -E "^$2$|^$2 " || { logger error "master $2 not existed!"; return 2; }

    NODE_NAME=$(bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" get node -owide|grep " $2 "|awk '{print $1}')

    logger warn "start to delete the master:$2 from cluster:$1"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/33.delmaster.yml" -e "NODE_TO_DEL=$2" -e "CLUSTER=$1" -e "@clusters/$1/config.yml"

    logger info "reconfig kubeconfig in ansible manage node"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/roles/deploy/deploy.yml" -t create_kctl_cfg -e "@clusters/$1/config.yml"

    logger info "reconfigure and restart 'kube-lb' service"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/90.setup.yml" -t restart_kube-lb -e "@clusters/$1/config.yml"

    logger info "reconfigure and restart 'ex-lb' service"
    ansible-playbook -i "$BASE/clusters/$1/hosts" "$BASE/playbooks/10.ex-lb.yml" -t restart_lb -e "@clusters/$1/config.yml"

    logger info "delete the master-node: $NODE_NAME"
    bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" delete node "$NODE_NAME"
}


function start-aio(){
    set +u
    # Check ENV 'HOST_IP', exists if the CMD 'ezctl' running in a docker container
    if [[ -z $HOST_IP ]];then
        # ezctl runs in a host machine, get host's ip
        HOST_IF=$(ip route|grep default|head -n1|cut -d' ' -f5)
        HOST_IP=$(ip a|grep "$HOST_IF$"|head -n1|awk '{print $2}'|cut -d'/' -f1)
    fi
    set -u
    logger info "get local host ipadd: $HOST_IP"

    # allow ssh login using key locally
    if [[ ! -e /root/.ssh/id_rsa ]]; then
      logger debug "generate ssh key pair"
      ssh-keygen -t rsa -b 2048 -N '' -f /root/.ssh/id_rsa > /dev/null
      cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
      ssh-keyscan -t ecdsa -H "$HOST_IP" >> /root/.ssh/known_hosts
    fi

    new default
    /bin/cp -f example/hosts.allinone "clusters/default/hosts"
    sed -i "s/_cluster_name_/default/g" "clusters/default/hosts"
    sed -i "s/192.168.1.1/$HOST_IP/g" "clusters/default/hosts"

    setup default all
}

### Extra functions #############################################
function renew-ca() {
    [[ -d "clusters/$1" ]] || { logger error "invalid cluster, run 'ezctl new $1' first"; return 1; }

    logger warn "WARNNING: this script should be used with greate caution"
    logger warn "WARNNING: it will recreate CA certs and all of the others certs used in the cluster"

    COMMAND="ansible-playbook -i clusters/$1/hosts -e @clusters/$1/config.yml -e CHANGE_CA=true playbooks/96.update-certs.yml -t force_change_certs"
    echo "$COMMAND"
    logger info "cluster:$1 process begins in 5s, press any key to abort:\n"
    ! (read -r -t5 -n1) || { logger warn "process abort"; return 1; }

    ${COMMAND} || return 1
}


EXPIRY="4800h"        # default cert will expire in 200 days
USER_TYPE="admin"     # admin/view, admin=clusterrole:cluster-admin view=clusterrole:view
USER_NAME="user"
function kcfg-adm() {
    OPTIND=2
    ACTION=""
    while getopts "ADLe:t:u:" OPTION; do
        case $OPTION in
          A)
            ACTION="add-kcfg $1"
            ;;
          D)
            ACTION="del-kcfg $1"
            ;;
          L)
            ACTION="list-kcfg $1"
            ;;
          e)
            EXPIRY="$OPTARG"
	    [[ $OPTARG =~ ^[1-9][0-9]*h$ ]] || { logger error "'-e' must be set like '2h, 5h, 50000h, ...'"; exit 1; }
            ;;
          t)
            USER_TYPE="$OPTARG"
	    [[ $OPTARG =~ ^(admin|view)$ ]] || { logger error "'-t' can only be set as 'admin' or 'view'"; exit 1; }
            ;;
          u)
            USER_NAME="$OPTARG"
            ;;
          ?)
            help-info kcfg-adm
            return 1
            ;;
        esac
    done

    [[ "$ACTION" == "" ]] && { logger error "illegal option"; help-info kcfg-adm; exit 1; }

    logger info "$ACTION"
    ${ACTION} || { logger error "$ACTION fail"; return 1; }
    logger info "$ACTION success"
}

function add-kcfg(){
    USER_NAME="$USER_NAME"-$(date +'%Y%m%d%H%M')
    logger info "add-kcfg in cluster:$1 with user:$USER_NAME"
    ansible-playbook -i "clusters/$1/hosts" -e "@clusters/$1/config.yml" -e "CUSTOM_EXPIRY=$EXPIRY" \
                     -e "USER_TYPE=$USER_TYPE" -e "USER_NAME=$USER_NAME" -e "ADD_KCFG=true" \
                     -t add-kcfg "roles/deploy/deploy.yml"
}

function del-kcfg(){
    logger info "del-kcfg in cluster:$1 with user:$USER_NAME"
    CRB=$(bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" get clusterrolebindings -ojsonpath="{.items[?(@.subjects[0].name == '$USER_NAME')].metadata.name}") && \
    bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" delete clusterrolebindings "$CRB" && \
    /bin/rm -f "clusters/$1/ssl/users/$USER_NAME"*
}

function list-kcfg(){
    logger info "list-kcfg in cluster:$1"
    ADMINS=$(bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" get clusterrolebindings -ojsonpath='{.items[?(@.roleRef.name == "cluster-admin")].subjects[*].name}')
    VIEWS=$(bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" get clusterrolebindings -ojsonpath='{.items[?(@.roleRef.name == "view")].subjects[*].name}')
    ALL=$(bin/kubectl --kubeconfig="clusters/$1/kubectl.kubeconfig" get clusterrolebindings -ojsonpath='{.items[*].subjects[*].name}')

    printf "\n%-30s %-15s %-20s\n" USER TYPE "EXPIRY(+8h if in Asia/Shanghai)"
    echo "---------------------------------------------------------------------------------"

    for u in $ADMINS; do
       if [[ $u =~ ^.*-[0-9]{12}$ ]];then
          t=$(bin/cfssl-certinfo -cert "clusters/$1/ssl/users/$u.pem"|grep not_after|awk '{print $2}'|sed 's/"//g'|sed 's/,//g')
          printf "%-30s %-15s %-20s\n" "$u" cluster-admin "$t"
       fi
    done;

    for u in $VIEWS; do
       if [[ $u =~ ^.*-[0-9]{12}$ ]];then
          t=$(bin/cfssl-certinfo -cert "clusters/$1/ssl/users/$u.pem"|grep not_after|awk '{print $2}'|sed 's/"//g'|sed 's/,//g')
          printf "%-30s %-15s %-20s\n" "$u" view "$t"
       fi
    done;

    for u in $ALL; do
       if [[ $u =~ ^.*-[0-9]{12}$ ]];then
          [[ $ADMINS == *$u* ]] || [[ $VIEWS == *$u* ]] || {
             t=$(bin/cfssl-certinfo -cert "clusters/$1/ssl/users/$u.pem"|grep not_after|awk '{print $2}'|sed 's/"//g'|sed 's/,//g')
             printf "%-30s %-15s %-20s\n" "$u" unknown "$t"
          }
       fi
    done;
    echo ""
}


### Main Lines ##################################################
function main() {
  BASE="/etc/kubeasz"
  [[ -d "$BASE" ]] || { logger error "invalid dir:$BASE, try: 'ezdown -D'"; exit 1; }
  cd "$BASE"

  # check bash shell
  readlink /proc/$$/exe|grep -q "bash" || { logger error "you should use bash shell only"; exit 1; }

  # check 'ansible' executable
  which ansible > /dev/null 2>&1 || { logger error "need 'ansible', try: 'pip install ansible'"; usage; exit 1; }

  [ "$#" -gt 0 ] || { usage >&2; exit 2; }

  case "$1" in
      ### in-cluster operations #####################
      (add-etcd)
          [ "$#" -gt 2 ] || { usage >&2; exit 2; }
          add-etcd "${@:2}"
          ;;
      (add-master)
          [ "$#" -gt 2 ] || { usage >&2; exit 2; }
          add-master "${@:2}"
          ;;
      (add-node)
          [ "$#" -gt 2 ] || { usage >&2; exit 2; }
          add-node "${@:2}"
          ;;
      (del-etcd)
          [ "$#" -eq 3 ] || { usage >&2; exit 2; }
          del-etcd "$2" "$3"
          ;;
      (del-master)
          [ "$#" -eq 3 ] || { usage >&2; exit 2; }
          del-master "$2" "$3"
          ;;
      (del-node)
          [ "$#" -eq 3 ] || { usage >&2; exit 2; }
          del-node "$2" "$3"
          ;;
      ### cluster-wide operations #######################
      (checkout)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          checkout "$2"
          ;;
      (list)
          [ "$#" -eq 1 ] || { usage >&2; exit 2; }
          list
          ;;
      (new)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          new "$2"
          ;;
      (setup)
          [ "$#" -ge 3 ] || { usage-setup >&2; exit 2; }
          setup "${@:2}"
          ;;
      (start)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" start
          ;;
      (stop)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" stop
          ;;
      (upgrade)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" upgrade
          ;;
      (backup)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" backup
          ;;
      (restore)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" restore
          ;;
      (destroy)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          cmd "$2" destroy
          ;;
      (start-aio)
          [ "$#" -eq 1 ] || { usage >&2; exit 2; }
          start-aio
          ;;
      ### extra operations ##############################
      (kca-renew)
          [ "$#" -eq 2 ] || { usage >&2; exit 2; }
          renew-ca "$2"
          ;;
      (kcfg-adm)
          [ "$#" -gt 2 ] || { usage-kcfg-adm >&2; exit 2; }
          kcfg-adm "${@:2}"
          ;;
      (help)
          [ "$#" -gt 1 ] || { usage >&2; exit 2; }
          help-info "$2"
          exit 0
          ;;
      (*)
          usage
          exit 0
          ;;
  esac
 }

main "$@"
