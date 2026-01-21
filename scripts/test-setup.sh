#!/bin/bash
set -ex

# K3S setup

if [ -f /etc/redhat-release ]; then
    sudo dnf install -y kernel-modules-extra
    systemctl disable firewalld --now
elif [ -f /etc/debian_version ]; then
    ufw disable
else
    echo "Host OS not supported"
    exit 1
fi
sudo curl -sfL https://get.k3s.io | sh -s -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Longhorn preflight check and cleanup
kubectl create ns longhorn-system
sudo rm -rf /usr/local/bin/longhornctl
curl -L https://github.com/longhorn/cli/releases/download/v1.10.1/longhornctl-linux-amd64 -o longhornctl
chmod +x longhornctl
sudo mv ./longhornctl /usr/local/bin/longhornctl
longhornctl install preflight
longhornctl check preflight
kubectl delete ns longhorn-system

# install mage
go install github.com/magefile/mage@latest

# helmfile setup
# if ! helm plugin list | grep -q "^diff[[:space:]]"; then
#     helm plugin install "https://github.com/databus23/helm-diff" --verify=false || helm plugin install "https://github.com/databus23/helm-diff"
# fi
# helmfile apply -f helmfile.yaml