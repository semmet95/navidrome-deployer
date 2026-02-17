#!/bin/bash
set -ex

kubectl get po -A -o wide
kubectl describe deploy longhorn -n navidrome-system
kubectl describe deploy navidrome -n navidrome-system
kubectl describe deploy filebrowser -n navidrome-system
kubectl describe job filebrowser-reconfig -n navidrome-system