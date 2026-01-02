package util

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
)

type Chart struct {
	ReleaseName   string
	LocalPath     string
	Namespace     string
	ReleaseValues map[string]string
	RepoUrl       string
}

type K3DCluster struct {
	AgentCount  string
	Name        string
	ServerCount string
}

func CreateCluster(ctx context.Context, t *testing.T, cluster K3DCluster) (string, error) {
	cmd := shell.Command{
		Command: "k3d",
		Args: []string{
			"cluster",
			"create",
			cluster.Name,
			"--agents",
			cluster.AgentCount,
			"--servers",
			cluster.ServerCount,
		},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return output, err
	}

	cmd = shell.Command{
		Command: "k3d",
		Args: []string{
			"kubeconfig",
			"write",
			cluster.Name,
		},
	}
	output, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return output, err
	}

	return output, nil
}

func DeleteCluster(ctx context.Context, t *testing.T, name string) (string, error) {
	cmd := shell.Command{
		Command: "k3d",
		Args: []string{
			"cluster",
			"delete",
			name,
		},
	}
	return shell.RunCommandAndGetOutputE(t, cmd)
}

func InstallHelmChartLocal(ctx context.Context, t *testing.T, chart Chart) error {
	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", chart.Namespace),
		SetValues:      chart.ReleaseValues,
	}
	return helm.InstallE(t, options, chart.LocalPath, chart.ReleaseName)
}
