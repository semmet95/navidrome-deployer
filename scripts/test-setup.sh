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
sudo curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="
  server \
  --kubelet-arg=v=1 \
  --kube-controller-manager-arg=v=1 \
  --kube-apiserver-arg=v=1 \
  --write-kubeconfig-mode=644 \
  --disable traefik \
  --disable servicelb \
  --disable metrics-server \
  --disable-cloud-controller \
  --disable coredns \
  --node-name ci-k3s \
" sh -s -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml