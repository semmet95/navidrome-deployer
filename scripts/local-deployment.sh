#!/bin/bash
set -ex

source test-setup.sh

# helmfile setup
if ! helm plugin list | grep -q "^diff[[:space:]]"; then
    helm plugin install "https://github.com/databus23/helm-diff" --verify=false
fi
helmfile apply -f helmfile.yaml