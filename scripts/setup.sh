#!/bin/bash
set -ex

# kubectl installation
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
mv -f ./kubectl /usr/local/bin/kubectl

# K3S setup
curl -sfL https://get.k3s.io | sh -s - --token 12345
chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
kubectl create ns longhorn-system

# Longhorn installation
rm -rf /usr/local/bin/longhornctl
curl -L https://github.com/longhorn/cli/releases/download/v1.10.1/longhornctl-linux-amd64 -o longhornctl
chmod +x longhornctl
mv ./longhornctl /usr/local/bin/longhornctl
longhornctl install preflight
longhornctl check preflight
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.10.1/deploy/longhorn.yaml

sleep 300
kubectl get po -A

# install mage
go install github.com/magefile/mage@latest

# K3S uninstallation
# /usr/local/bin/k3s-uninstall.sh