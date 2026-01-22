#!/bin/bash
set -ex

source local-deployment.sh

# helmfile setup
if ! helm plugin list | grep -q "^diff[[:space:]]"; then
    helm plugin install "https://github.com/databus23/helm-diff" --verify=false || helm plugin install "https://github.com/databus23/helm-diff"
fi
helmfile apply -f helmfile.yaml