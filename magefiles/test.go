//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

const (
	agentCount     = "1"
	appName        = "navidrome"
	defaultCluster = "test-env"
	serverCount    = "1"
)

func (Test) Setup() {
	createCluster(context.TODO(), defaultCluster)
}

func (Test) DeployApp() {
	mg.Deps(Test.Setup)
	applyManifest(context.TODO(), "app.yml")
	verifyDeploymentHealth(context.TODO(), appName, "default")
}

func (Test) Cleanup() {
	deleteCluster(context.TODO(), defaultCluster)
}

func createCluster(ctx context.Context, name string) {
	output, err := sh.Output(
		"k3d",
		"cluster",
		"create",
		name,
		"--agents",
		agentCount,
		"--servers",
		serverCount,
	)

	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}

func deleteCluster(ctx context.Context, name string) {
	output, err := sh.Output(
		"k3d",
		"cluster",
		"delete",
		name,
	)

	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}

func applyManifest(ctx context.Context, filePath string) {
	output, err := sh.Output(
		"kubectl",
		"apply",
		"-f",
		filePath,
	)

	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}

func verifyDeploymentHealth(ctx context.Context, deployName, namespace string) {
	output, err := sh.Output(
		"kubectl",
		"wait",
		"deployment",
		deployName,
		"-n",
		namespace,
		"--for=condition=Available",
		"--timeout=120s",
	)

	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}
