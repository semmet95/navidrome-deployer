//go:build mage
// +build mage

package main

import (
	"context"
	"navidrome-deployer/test/util"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/magefile/mage/mg"
)

type Test mg.Namespace

const (
	// helm config
	longhorChartRepoName = "longhorn"
	longhorChartRepoUrl  = "https://charts.longhorn.io"
	releaseName          = "navidrome-deployer"
	releaseNamespace     = "default"

	// app config
	appName                = "navidrome"
	defaultLonghornVersion = "v1.10.1"
)

var (
	kubeConfigPath string
)

func (Test) DeployApp() {
	chartPath, err := filepath.Abs("charts/navidrome-deployer")
	if err != nil {
		panic(err)
	}

	err = util.AddHelmRepo(context.TODO(), &testing.T{}, longhorChartRepoName, longhorChartRepoUrl)
	if err != nil {
		panic(err)
	}

	localChart := util.Chart{
		ReleaseName: releaseName,
		LocalPath:   chartPath,
		Namespace:   releaseNamespace,
	}
	err = util.InstallHelmChartLocal(context.TODO(), &testing.T{}, localChart)
	if err != nil {
		panic(err)
	}

	err = k8s.WaitUntilDeploymentAvailableE(
		&testing.T{},
		&k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  "default",
		},
		appName,
		8,
		15*time.Second,
	)
	if err != nil {
		panic(err)
	}
}
