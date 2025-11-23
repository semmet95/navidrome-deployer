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
	defaultCluster = "test-env"
	serverCount    = "1"
)

func (Test) Setup() {
	createCluster(context.TODO(), defaultCluster)
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
