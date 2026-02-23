#!/bin/bash
set -x

kubectl get po -A -o wide
kubectl describe deploy navidrome -n navidrome-system
kubectl describe deploy filebrowser -n navidrome-system
kubectl describe job filebrowser-reconfig -n navidrome-system
kubectl get ClusterIssuer -A
kubectl get Certificate -A
kubectl get ing -A -o wide