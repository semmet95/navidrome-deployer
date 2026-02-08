#!/bin/bash
set -ex

kubectl create ns longhorn-system
sudo rm -rf /usr/local/bin/longhornctl
curl -L https://github.com/longhorn/cli/releases/download/v1.11.0/longhornctl-linux-amd64 -o longhornctl
chmod +x longhornctl
sudo mv ./longhornctl /usr/local/bin/longhornctl
longhornctl install preflight
longhornctl check preflight
kubectl delete ns longhorn-system