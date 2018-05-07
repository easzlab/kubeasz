#!/bin/bash

# 不同k8s版本使用的'api-versions'版本不同，此脚本用于切换yaml文件使用的'api-versions'

#set -x

show_usage()
{
   echo -e "\nUsage: $0 <-v K8S_VER> <-f YAML_FILE>"
   echo -e "\nK8S_VER: support 1.8/1.9/1.10"
}

#check_arg -------------------------------------------------
K8S_VER=""
YML_FILE=""

while getopts "v:f:" arg
do
   case $arg in
      v)
         K8S_VER=$OPTARG
         ;;
      f)
         if [ -w "$OPTARG" ];then
            YML_FILE=$OPTARG
         else
            echo File:"$OPTARG not found or not writeable."
            exit 1
         fi
         ;;
      ?)
         echo -e "unkown argument"
         show_usage
         exit 1
         ;;
   esac
done

if [ "$K8S_VER" = "" ] || [ "$YML_FILE" = "" ];then
   echo "error argument"
   show_usage
   exit 1
fi

main()
{
   case "$K8S_VER" in
      1.8)
         sed -i 's/apps\/v1/extensions\/v1beta1/g' $YML_FILE
         exit 0
         ;;
      1.9)
         echo "K8s_VER is $K8S_VER"
         exit 0
         ;;
      1.10)
         sed -i 's/extensions\/v1beta1/apps\/v1/g' $YML_FILE
         exit 0
         ;;
      ?)
         ;;
   esac
   echo "not supported K8s_VER:$K8S_VER"
   exit 1
}

main
