# Navidrome Deployer

[![E2E Tests](https://github.com/semmet95/navidrome-deployer/actions/workflows/e2e-tests.yml/badge.svg)](https://github.com/semmet95/navidrome-deployer/actions/workflows/e2e-tests.yml)

## Prerequisites

Before deploying `navidrome-deployer`, ensure the following tools are installed:

- [kubectl](https://kubernetes.io/docs/tasks/tools/) - Kubernetes command-line tool
- [go](https://golang.org/doc/install) - Go programming language
- [mage](https://magefile.org/) - Go-based task runner
- [helm](https://helm.sh/docs/intro/install/) - Kubernetes package manager
- [helmfile](https://github.com/roboll/helmfile#installation) - Helm values file manager

## Installation on a cluster
To install the latest release
```bash
helmfile apply -f https://github.com/semmet95/navidrome-deployer/releases/latest/download/helmfile.yaml
```

To install a specific version
```bash
helmfile apply -f https://github.com/semmet95/navidrome-deployer/releases/download/<version>/helmfile.yaml
```

## Local Setup and Deployment

To deploy navidrome-deployer locally, execute the following command:

```bash
./scripts/local-deployment.sh
```
