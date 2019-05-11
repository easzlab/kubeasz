#!/bin/bash
#--------------------------------------------------
# Rotate the indices in elastic of the EFK deployment
#
# @author:  gjmzj
# @usage:   ./rotate.sh [num_of_days] (1<num_of_days<999)
# @repo:    https://github.com/kubeasz/mirrorepo/es-index-rotator
# @ref:     https://github.com/easzlab/kubeasz/tree/master/dockerfiles/es-index-rotator

set -o nounset
set -o errexit
set -o xtrace

# max days of logs to keep, default=7
max_days_of_log=7

if [[ "$#" -gt 0 && $1 =~ ^[1-9][0-9]{0,2}$ ]];then
    max_days_of_log=$1
fi

echo -e "\n[INFO] rotate job starts, try to keep $max_days_of_log days of logs."

curl elasticsearch-logging:9200/_cat/indices? > /tmp/indices || \
{ echo "[ERROR] Can not connect to elastic!"; exit 1; }

curr_days_of_log=$(cat /tmp/indices|grep logstash|wc -l)

curr_days_of_log=$((${curr_days_of_log}-2))

if [[ "$max_days_of_log" -gt "$curr_days_of_log" ]];then
    echo "[WARN] No need to rotate the ES indices!"
    exit 0
fi  

first_day=$(date -d "$max_days_of_log days ago" +'%Y.%m.%d')

rotate=$(cat /tmp/indices|grep logstash|cut -d' ' -f3|cut -d'-' -f2|sort|sed -n "1,/$first_day/"p)

for day in $rotate;do
    curl -X DELETE elasticsearch-logging:9200/logstash-$day
done

echo -e "\n[INFO] Success to rotate the ES indices!"
exit 0
