
# UNCOMMENT these when we are running in K8S againf
docker build -t martian . 
docker tag martian 192.168.1.19:30500/martian:latest
docker push 192.168.1.19:30500/martian:latest

MARTIAN_CONTAINER=`kubectl get pods | grep martian | awk '{ print $1 }'`
kubectl delete pod $MARTIAN_CONTAINER