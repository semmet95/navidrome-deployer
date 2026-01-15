//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
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
	longhornNamespace    = "longhorn-system"
	releaseName          = "navidrome-deployer"
	releaseNamespace     = "default"

	// app config
	appName                = "navidrome"
	defaultLonghornVersion = "v1.10.1"
)

var (
	kubeConfigPath string

	allDeployments = [3]util.K8SResource{
		{
			Name:      appName,
			Namespace: "default",
		},
		{
			Name:      "longhorn-driver-deployer",
			Namespace: longhornNamespace,
		},
		{
			Name:      "longhorn-ui",
			Namespace: longhornNamespace,
		},
	}
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

	for _, deploy := range allDeployments {
		opts := &k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  deploy.Namespace,
		}
		err = k8s.WaitUntilDeploymentAvailableE(
			&testing.T{},
			opts,
			deploy.Name,
			8,
			15*time.Second,
		)
		if err != nil {
			logs, logErr := util.GetAppLogs(
				context.TODO(),
				&testing.T{},
				opts,
				deploy.Name,
			)
			if logErr != nil {
				fmt.Println("error while retriving %s app logs: %v", deploy.Name, logErr)
			} else {
				fmt.Println(logs)
			}
			panic(err)
		}
	}
}
