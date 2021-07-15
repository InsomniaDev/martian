#!/bin/sh

yarn build
docker build -t portal . 
docker tag portal 192.168.1.19:30500/portal:latest
docker push 192.168.1.19:30500/portal:latest

PORTAL_CONTAINER=`kubectl get pods | grep portal | awk '{ print $1 }'`
kubectl delete pod $PORTAL_CONTAINER