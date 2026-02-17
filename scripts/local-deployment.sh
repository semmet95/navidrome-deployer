#!/bin/bash
set -ex

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

source "$SCRIPT_DIR/test-setup.sh"
source "$SCRIPT_DIR/longhorn-preflight.sh"

# build the filebrowser-reconfig image
docker build -t filebrowser-reconfig:snapshot . -f Dockerfile.filebrowser
docker save filebrowser-reconfig:snapshot | sudo k3s ctr images import -

# helmfile setup
if ! helm plugin list | grep -q "^diff[[:space:]]"; then
    helm plugin install "https://github.com/databus23/helm-diff" --verify=false
fi
helmfile apply -f test/helmfile.yaml