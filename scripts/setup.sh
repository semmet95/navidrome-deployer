#!/bin/bash
set -ex

# install k3d
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# install mage
go install github.com/magefile/mage@latest