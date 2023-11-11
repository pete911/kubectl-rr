package internal

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-rr/internal/k8s"
	"k8s.io/client-go/rest"
	"time"
)

type Pod struct {
	Name       string
	Namespace  string
	Containers []Container
}

type Container struct {
	Name   string
	Image  string
	CPU    Resource
	Memory Resource
}

type Resource struct {
	Metric
	Request string
	Limit   string
}

type Metric struct {
	Current string
	Min     string
	Max     string
}

func GetPods(restConfig *rest.Config, namespace, labelSelector, fieldSelector string) ([]Pod, error) {
	client, err := k8s.NewClient(restConfig)
	if err != nil {
		return nil, err
	}

	prom, err := k8s.NewPrometheus(client)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	k8sPods, err := client.GetPods(ctx, namespace, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}
	return toPods(k8sPods, prom)
}

func toPods(k8sPods []k8s.Pod, prom k8s.Prometheus) ([]Pod, error) {
	var out []Pod
	for _, k8sPod := range k8sPods {
		pod := Pod{
			Name:      k8sPod.Name,
			Namespace: k8sPod.Namespace,
		}
		for _, k8sContainer := range k8sPod.Containers {
			metric, err := toCPUMetric(prom, pod.Namespace, pod.Name, k8sContainer.Name)
			if err != nil {
				return nil, fmt.Errorf("get cpu metric %s/%s container %s: %w", pod.Namespace, pod.Name, k8sContainer.Name, err)
			}
			pod.Containers = append(pod.Containers, Container{
				Name:   k8sContainer.Name,
				Image:  "", // TODO
				CPU:    Resource{Metric: metric, Request: k8sContainer.Requests.Cpu.String(), Limit: k8sContainer.Limits.Cpu.String()},
				Memory: Resource{}, // TODO
			})
		}
		out = append(out, pod)
	}
	return out, nil
}

func toCPUMetric(prom k8s.Prometheus, namespace, pod, container string) (Metric, error) {
	cpu, err := prom.CPU(namespace, pod, container)
	if err != nil {
		return Metric{}, err
	}
	minCpu, err := prom.MinCPU(namespace, pod, container)
	if err != nil {
		return Metric{}, err
	}
	maxCpu, err := prom.MaxCPU(namespace, pod, container)
	if err != nil {
		return Metric{}, err
	}
	return Metric{
		Current: cpu,
		Min:     minCpu,
		Max:     maxCpu,
	}, nil
}
