#!/bin/bash
set -ex

# K3S setup
if [ -f /etc/redhat-release ]; then
    sudo dnf install -y kernel-modules-extra
    systemctl disable firewalld --now
elif [ -f /etc/debian_version ]; then
    # the following changes are relevant for free tier GA runners
    sudo ufw disable
    sudo apt-get update
    sudo apt-get install -y nfs-common
else
    echo "Host OS not supported"
    exit 1
fi

sudo curl -sfL https://get.k3s.io | sh -s - server --disable traefik --disable servicelb --disable metrics-server --disable-cloud-controller
sudo chmod 644 /etc/rancher/k3s/k3s.yaml