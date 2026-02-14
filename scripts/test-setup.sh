#!/bin/bash
set -ex

# K3S setup
if [ -f /etc/redhat-release ]; then
    sudo dnf install -y kernel-modules-extra
    sudo firewall-cmd --permanent --add-port=6443/tcp #apiserver
    sudo firewall-cmd --permanent --zone=trusted --add-source=10.42.0.0/16 #pods
    sudo firewall-cmd --permanent --zone=trusted --add-source=10.43.0.0/16 #services
    sudo firewall-cmd --reload
elif [ -f /etc/debian_version ]; then
    # the following changes are relevant for free tier GA runners
    sudo apt-get update
    sudo apt-get install -y nfs-common
    sudo ufw allow 6443/tcp #apiserver
    sudo ufw allow from 10.42.0.0/16 to any #pods
    sudo ufw allow from 10.43.0.0/16 to any #services
else
    echo "Host OS not supported"
    exit 1
fi

sudo curl -sfL https://get.k3s.io | sh -s - server \
    --disable-cloud-controller \
    --disable=servicelb \
    --etcd-disable-snapshots
sudo chmod 644 /etc/rancher/k3s/k3s.yaml