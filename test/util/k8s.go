package util

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	v1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Chart struct {
	ReleaseName   string
	LocalPath     string
	Namespace     string
	ReleaseValues map[string]string
	RepoUrl       string
}

func GetDeploymentPods(ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, deployment *v1.Deployment) ([]corev1.Pod, error) {
	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	pods, err := k8s.ListPodsE(t, opts, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("error while getting all the pods for deployment %s : %v", deployment.Name, err)
	}
	return pods, nil
}

func GetDeploymentLogs(ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, deployment *v1.Deployment) (string, error) {
	pods, err := GetDeploymentPods(ctx, t, opts, deployment)
	if err != nil {
		return "", err
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

func DescribePod(ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, podName string) (string, error) {
	output, err := k8s.RunKubectlAndGetOutputE(t, opts, "describe", "pod", podName)
	if err != nil {
		return "", fmt.Errorf("failed to describe pod %s: %v", podName, err)
	}
	return output, nil
}

func WaitUntilJobCompletes(ctx context.Context, t *testing.T, opts *k8s.KubectlOptions, jobName string) error {
	return retry.New(
		retry.Attempts(10),
		retry.Delay(30*time.Second),
	).Do(
		func() error {
			var errMsg string
			job, err := k8s.GetJobE(t, opts, jobName)
			if err != nil {
				return err
			}

			for _, jobCondition := range job.Status.Conditions {
				fmt.Printf("job condition type is %s and status is %s\n", jobCondition.Type, jobCondition.Status)
				if jobCondition.Type == batchV1.JobComplete && jobCondition.Status == corev1.ConditionTrue {
					return nil
				} else {
					errMsg = fmt.Sprintf(
						"job completion status is %v with message: %s",
						jobCondition.Status,
						jobCondition.Message,
					)
					fmt.Println(errMsg)
					return errors.New(errMsg)
				}
			}
			errMsg = fmt.Sprintf("failed to iterate over job conditions: %v", job.Status.Conditions)
			fmt.Println(errMsg)
			return errors.New(errMsg)
		},
	)
}
