package internal

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-rr/internal/k8s"
	"k8s.io/client-go/rest"
	"sort"
	"time"
)

type Pod struct {
	Name           string
	Namespace      string
	Containers     []Container
	InitContainers []Container
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
	Current float64
	Min     float64
	Max     float64
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
	defer prom.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	k8sPods, err := client.GetPods(ctx, namespace, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}
	pods, err := toPods(k8sPods, prom)
	if err != nil {
		return nil, err
	}

	sort.Slice(pods, func(i, j int) bool {
		if pods[i].Namespace != pods[j].Namespace {
			return pods[i].Namespace < pods[j].Namespace
		}
		return pods[i].Name < pods[j].Name
	})
	return pods, nil
}

func toPods(k8sPods []k8s.Pod, prom k8s.Prometheus) ([]Pod, error) {
	var out []Pod
	for _, k8sPod := range k8sPods {
		pod := Pod{
			Name:      k8sPod.Name,
			Namespace: k8sPod.Namespace,
		}
		for _, k8sContainer := range k8sPod.Containers {
			container, err := toContainer(k8sPod, k8sContainer, prom)
			if err != nil {
				return nil, err
			}
			pod.Containers = append(pod.Containers, container)
		}
		for _, k8sContainer := range k8sPod.InitContainers {
			container, err := toContainer(k8sPod, k8sContainer, prom)
			if err != nil {
				return nil, err
			}
			pod.InitContainers = append(pod.InitContainers, container)
		}
		out = append(out, pod)
	}
	return out, nil
}

func toContainer(k8sPod k8s.Pod, k8sContainer k8s.Container, prom k8s.Prometheus) (Container, error) {
	cpuMetric, err := toCPUMetric(prom, k8sPod.Namespace, k8sPod.Name, k8sContainer)
	if err != nil {
		return Container{}, fmt.Errorf("get cpu metric %s/%s container %s: %w", k8sPod.Namespace, k8sPod.Name, k8sContainer.Name, err)
	}

	memoryMetric, err := toMemoryMetric(prom, k8sPod.Namespace, k8sPod.Name, k8sContainer)
	if err != nil {
		return Container{}, fmt.Errorf("get memory metric %s/%s container %s: %w", k8sPod.Namespace, k8sPod.Name, k8sContainer.Name, err)
	}

	// 1.0 - 1 CPU
	// 0.1 - 100m (millicpu)
	// min for requests/limits is 1m - 0.001
	return Container{
		Name:   k8sContainer.Name,
		Image:  k8sContainer.Image,
		CPU:    Resource{Metric: cpuMetric, Request: k8sContainer.Requests.Cpu.String(), Limit: k8sContainer.Limits.Cpu.String()},
		Memory: Resource{Metric: memoryMetric, Request: k8sContainer.Requests.Memory.String(), Limit: k8sContainer.Limits.Memory.String()},
	}, nil
}

func toCPUMetric(prom k8s.Prometheus, namespace, pod string, container k8s.Container) (Metric, error) {
	cpu, err := prom.CPU(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	minCpu, err := prom.MinCPU(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	maxCpu, err := prom.MaxCPU(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	return Metric{
		Current: cpu,
		Min:     minCpu,
		Max:     maxCpu,
	}, nil
}

func toMemoryMetric(prom k8s.Prometheus, namespace, pod string, container k8s.Container) (Metric, error) {
	memory, err := prom.Memory(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	minMemory, err := prom.MinMemory(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	maxMemory, err := prom.MaxMemory(namespace, pod, container.Name)
	if err != nil {
		return Metric{}, err
	}
	return Metric{
		Current: memory,
		Min:     minMemory,
		Max:     maxMemory,
	}, nil
}
