#!/bin/bash

kubectl get po -A -o wide
echo '-----------------------------------------------------------------------------------------------'
kubectl describe deploy navidrome -n navidrome-system
echo '-----------------------------------------------------------------------------------------------'
kubectl describe deploy filebrowser -n navidrome-system
echo '-----------------------------------------------------------------------------------------------'
kubectl describe job filebrowser-reconfig -n navidrome-system
echo '-----------------------------------------------------------------------------------------------'
kubectl get ClusterIssuer -A
echo '-----------------------------------------------------------------------------------------------'
kubectl get Certificate -A
echo '-----------------------------------------------------------------------------------------------'
kubectl get IngressRoute -A -o wide
echo '-----------------------------------------------------------------------------------------------'
kubectl logs -l app=filebrowser --all-containers=true --prefix=true -n navidrome-system
echo '-----------------------------------------------------------------------------------------------'
kubectl logs -l app=navidrome --all-containers=true --prefix=true -n navidrome-system