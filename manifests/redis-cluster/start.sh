#!/bin/sh
set -x

ROOT=$(cd `dirname $0`; pwd)
cd $ROOT

helm install redis \
	--create-namespace \
	--namespace dependency \
	-f ./values.yaml \
	./redis-ha
