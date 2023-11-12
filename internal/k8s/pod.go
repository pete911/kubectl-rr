package k8s

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Namespaces(pods []Pod) []string {
	namespaces := make(map[string]struct{})
	for _, v := range pods {
		namespaces[v.Namespace] = struct{}{}
	}

	var out []string
	for k := range namespaces {
		out = append(out, k)
	}
	return out
}

func toPods(pods []v1.Pod) []Pod {
	var out []Pod
	for _, pod := range pods {
		out = append(out, toPod(pod))
	}
	return out
}

type Pod struct {
	Name           string
	Namespace      string
	Containers     []Container
	InitContainers []Container
}

func toPod(pod v1.Pod) Pod {
	return Pod{
		Name:           pod.Name,
		Namespace:      pod.Namespace,
		Containers:     toContainers(pod.Spec.Containers),
		InitContainers: toContainers(pod.Spec.InitContainers),
	}
}

func toContainers(containers []v1.Container) []Container {
	var out []Container
	for _, c := range containers {
		out = append(out, toContainer(c))
	}
	return out
}

type Container struct {
	Name     string
	Image    string
	Requests Resource
	Limits   Resource
}

type Resource struct {
	Cpu    *resource.Quantity
	Memory *resource.Quantity
}

func toContainer(container v1.Container) Container {
	return Container{
		Name:  container.Name,
		Image: container.Image,
		Requests: Resource{
			Cpu:    container.Resources.Requests.Cpu(),
			Memory: container.Resources.Requests.Memory(),
		},
		Limits: Resource{
			Cpu:    container.Resources.Limits.Cpu(),
			Memory: container.Resources.Limits.Memory(),
		},
	}
}
