//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"navidrome-deployer/test/util"
	"os"
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
	navidromeNamespace = "navidrome-system"
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

func (Test) CheckDeployments() error {
	longhornDeployments, err := k8s.ListDeploymentsE(
		&testing.T{},
		&k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  longhornNamespace,
		},
		v1.ListOptions{},
	)
	if err != nil {
		panic(err)
	}

	navidromeDeployments, err := k8s.ListDeploymentsE(
		&testing.T{},
		&k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  navidromeNamespace,
		},
		v1.ListOptions{},
	)
	if err != nil {
		panic(err)
	}

	for _, deploy := range append(longhornDeployments, navidromeDeployments...) {
		opts := &k8s.KubectlOptions{
			ConfigPath: kubeConfigPath,
			Namespace:  deploy.Namespace,
		}
		err = k8s.WaitUntilDeploymentAvailableE(
			&testing.T{},
			opts,
			deploy.Name,
			12,
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
				fmt.Printf("error while retrieving %s app logs: %v\n", deploy.Name, logErr)
			} else {
				fmt.Println(logs)
			}
			panic(err)
		}
	}

	return nil
}

func (Test) DeployApp() {
	mg.Deps(Test.CheckDeployments)
}
