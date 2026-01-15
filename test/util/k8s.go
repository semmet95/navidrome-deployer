package util

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

type Chart struct {
	ReleaseName   string
	LocalPath     string
	Namespace     string
	ReleaseValues map[string]string
	RepoUrl       string
}

func InstallHelmChartLocal(ctx context.Context, t *testing.T, chart Chart) error {
	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    k8s.NewKubectlOptions("", "", chart.Namespace),
		SetValues:         chart.ReleaseValues,
	}
	return helm.InstallE(t, options, chart.LocalPath, chart.ReleaseName)
}

// func VerifyDaemonSet(
// 	ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, name string, retryCount int, waitDuration time.Duration) error {
// 	ds := k8s.GetDaemonSet(t, opts, "test")
// 	ds.Status.DesiredNumberScheduled
// }

func VerifyDeployment(
	ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, name string, retryCount int, waitDuration time.Duration) error {
	return k8s.WaitUntilDeploymentAvailableE(
		t,
		opts,
		name,
		retryCount,
		waitDuration,
	)
}
