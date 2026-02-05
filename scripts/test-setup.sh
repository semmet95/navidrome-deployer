#!/bin/bash
set -ex

# K3S setup
if [ -f /etc/redhat-release ]; then
    sudo dnf install -y kernel-modules-extra
    systemctl disable firewalld --now
elif [ -f /etc/debian_version ]; then
    sudo ufw disable
else
    echo "Host OS not supported"
    exit 1
fi
sudo curl -sfL https://get.k3s.io | sh -s -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml