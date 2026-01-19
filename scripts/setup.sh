#!/bin/bash
set -ex

# K3S setup
sudo curl -sfL https://get.k3s.io | sh -s - --token 12345
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
kubectl create ns longhorn-system

# Longhorn setup
rm -rf /usr/local/bin/longhornctl
curl -L https://github.com/longhorn/cli/releases/download/v1.10.1/longhornctl-linux-amd64 -o longhornctl
chmod +x longhornctl
mv ./longhornctl /usr/local/bin/longhornctl
longhornctl install preflight
longhornctl check preflight

# install mage
go install github.com/magefile/mage@latest