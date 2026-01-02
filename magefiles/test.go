//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"navidrome-deployer/test/util"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/magefile/mage/mg"
)

type Test mg.Namespace

const (
	// cluster config
	agentCount     = "1"
	defaultCluster = "test-env"
	serverCount    = "1"

	// helm config
	releaseName      = "navidrome-deployer"
	releaseNamespace = "default"

	// app config
	appName = "navidrome"
)

var (
	kubeConfigPath string
)

func (Test) Setup() {
	cluster := util.K3DCluster{
		AgentCount:  agentCount,
		Name:        defaultCluster,
		ServerCount: serverCount,
	}
	output, err := util.CreateCluster(context.TODO(), &testing.T{}, cluster)
	if err != nil {
		fmt.Println(output)
		panic(err)
	}

	kubeConfigPath = output
	err = os.Setenv("KUBECONFIG", kubeConfigPath)
	if err != nil {
		panic(err)
	}
}

func (Test) DeployApp() {
	mg.Deps(Test.Setup)

	chartPath, err := filepath.Abs("charts/navidrome-deployer")
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

func (Test) Cleanup() {
	output, err := util.DeleteCluster(context.TODO(), &testing.T{}, defaultCluster)
	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}
