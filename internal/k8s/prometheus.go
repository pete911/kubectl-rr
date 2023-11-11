package k8s

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"
	"net/url"
	"strings"
	"time"
)

var (
	// TODO - make configurable, either flags or config file
	promNamespace = ""
	promLabels    = map[string]string{
		"app.kubernetes.io/name":     "prometheus",
		"app.kubernetes.io/instance": "prometheus",
	}
	promPort = "9090"
)

type Prometheus struct {
	forwarder Forwarder
}

func NewPrometheus(client Client) (Prometheus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.GetPods(ctx, promNamespace, toLabels(promLabels), "")
	if err != nil {
		return Prometheus{}, err
	}

	namespaces := Namespaces(pods)
	if len(namespaces) != 1 {
		return Prometheus{}, fmt.Errorf("prometheus expected in one namespace, but found in %v", namespaces)
	}

	forwarder, err := StartForwarder(client.RestConfig, pods[0].Namespace, pods[0].Name, promPort)
	if err != nil {
		return Prometheus{}, fmt.Errorf("start port forward: %w", err)
	}
	return Prometheus{forwarder: forwarder}, nil
}

func (p Prometheus) CPU(namespace, pod, container string) (string, error) {
	return p.query(fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m]))`, namespace, pod, container))
}

func (p Prometheus) MinCPU(namespace, pod, container string) (string, error) {
	return p.query(fmt.Sprintf(`sum(min_over_time(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])[30m:1m]))`, namespace, pod, container))
}

func (p Prometheus) MaxCPU(namespace, pod, container string) (string, error) {
	return p.query(fmt.Sprintf(`sum(max_over_time(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])[30m:1m]))`, namespace, pod, container))
}

func (p Prometheus) Stop() {
	p.forwarder.Stop()
}

func (p Prometheus) query(query string) (string, error) {
	params := url.Values{"query": []string{query}}
	b, err := p.forwarder.Get("/api/v1/query", params)
	if err != nil {
		return "", err
	}

	var response QueryResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return "", err
	}

	if response.Status != "success" {
		return "", fmt.Errorf("unexpected status %s from prometheus", response.Status)
	}
	if response.Data.ResultType != "vector" {
		return "", fmt.Errorf("unexpected result type %s from prometheus", response.Data.ResultType)
	}
	return fmt.Sprintf("%v", response.Data.Results[0].Value[1]), err
}

func toLabels(l map[string]string) string {
	var out []string
	for k, v := range l {
		out = append(out, fmt.Sprintf("%s=%s", k, v)) // TODO - fix this
	}
	return strings.Join(out, ",")
}

type QueryResponse struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

type QueryData struct {
	ResultType string        `json:"resultType"` // "matrix" | "vector" | "scalar" | "string"
	Results    []QueryResult `json:"result"`
}

type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"` // TODO - fix this
}
