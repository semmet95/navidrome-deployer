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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Test mg.Namespace

const (
	// helm config
	longhornNamespace = "longhorn-system"
	releaseName       = "navidrome-deployer"
	releaseNamespace  = "default"

	// app config
	appName                = "navidrome"
	defaultLonghornVersion = "v1.10.1"
)

var (
	kubeConfigPath string
)

func init() {
	configPath, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		panic("KUBECONFIG environment variable is not set")
	}
	kubeConfigPath = configPath
}

func (Test) CheckDependencies() error {
	longhornDeployments, err := k8s.ListDeploymentsE(
		&testing.T{},
		&k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  longhornNamespace,
		},
		v1.ListOptions{},
	)
	if err != nil {
		return err
	}

	for _, deploy := range longhornDeployments {
		opts := &k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  deploy.Namespace,
		}
		err = k8s.WaitUntilDeploymentAvailableE(
			&testing.T{},
			opts,
			deploy.Name,
			8,
			30*time.Second,
		)
		if err != nil {
			logs, logErr := util.GetDeploymentLogs(
				context.TODO(),
				&testing.T{},
				opts,
				&deploy,
			)
			if logErr != nil {
				fmt.Println("error while retriving %s app logs: %v", deploy.Name, logErr)
			} else {
				fmt.Println(logs)
			}
			panic(err)
		}
	}

	return nil
}

func (Test) DeployApp() {
	mg.Deps(Test.CheckDependencies)
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

	opts := &k8s.KubectlOptions{
		ConfigPath: kubeConfigPath,
		Namespace:  releaseNamespace,
	}
	err = k8s.WaitUntilDeploymentAvailableE(
		&testing.T{},
		opts,
		appName,
		8,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}
}
