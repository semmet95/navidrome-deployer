package util

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Chart struct {
	ReleaseName   string
	LocalPath     string
	Namespace     string
	ReleaseValues map[string]string
	RepoUrl       string
}

func GetDeploymentLogs(ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, deployment *v1.Deployment) (string, error) {
	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	pods, err := k8s.ListPodsE(t, opts, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return "", fmt.Errorf("error while getting all the pods for deployment %s : %v", deployment.Name, err)
	}

	var allLogs []string
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			logs, err := k8s.GetPodLogsE(t, opts, &pod, container.Name)
			if err != nil {
				return "", fmt.Errorf("failed to get logs from pod %s container %s: %v", pod.Name, container.Name, err)
			}
			if logs != "" {
				allLogs = append(allLogs, fmt.Sprintf("=== Pod: %s | Container: %s ===\n%s", pod.Name, container.Name, logs))
			}
		}
	}

	return strings.Join(allLogs, "\n"), nil
}

func InstallHelmChartLocal(ctx context.Context, t *testing.T, chart Chart) error {
	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    k8s.NewKubectlOptions("", "", chart.Namespace),
		SetValues:         chart.ReleaseValues,
	}
	return helm.InstallE(t, options, chart.LocalPath, chart.ReleaseName)
}

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
