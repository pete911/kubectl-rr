package k8s

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type Client struct {
	RestConfig *rest.Config
	coreV1     corev1.CoreV1Interface
}

func NewClient(restConfig *rest.Config) (Client, error) {
	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return Client{}, err
	}
	return Client{
		RestConfig: restConfig,
		coreV1:     cs.CoreV1(),
	}, nil
}

func (c Client) GetPods(ctx context.Context, namespace, labelSelector, fieldSelector string) ([]Pod, error) {
	if namespace == "" {
		return c.getAllPods(ctx, labelSelector, fieldSelector)
	}

	podList, err := c.coreV1.Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector, FieldSelector: fieldSelector})
	if err != nil {
		return nil, err
	}
	return toPods(podList.Items), nil
}

func (c Client) getAllPods(ctx context.Context, labelSelector, fieldSelector string) ([]Pod, error) {
	namespaces, err := c.getNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}

	var pods []Pod
	for _, namespace := range namespaces {
		p, err := c.GetPods(ctx, namespace.Name, labelSelector, fieldSelector)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p...)
	}
	return pods, nil
}

func (c Client) getNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	namespaceList, err := c.coreV1.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}
	return namespaceList.Items, nil
}
