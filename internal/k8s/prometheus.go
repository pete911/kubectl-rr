package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"
)

var (
	promPort = "9090"
)

type PrometheusConfig struct {
	Namespace string
	Labels    string
}

type Prometheus struct {
	forwarder Forwarder
}

func NewPrometheus(client Client, config PrometheusConfig) (Prometheus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.GetPods(ctx, config.Namespace, config.Labels, "")
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

func (p Prometheus) CPU(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])`, namespace, pod, container))
}

// TODO - only call if the container has cpu limit set
func (p Prometheus) ThrottledCPU(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`rate(container_cpu_cfs_throttled_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])`, namespace, pod, container))
}

func (p Prometheus) MinCPU(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`min_over_time(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])[30m:1m])`, namespace, pod, container))
}

func (p Prometheus) MaxCPU(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`max_over_time(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s", container="%s"}[5m])[30m:1m])`, namespace, pod, container))
}

func (p Prometheus) Memory(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`rate(container_memory_working_set_bytes{namespace="%s", pod="%s", container="%s"}[5m])`, namespace, pod, container))
}

func (p Prometheus) MinMemory(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`min_over_time(rate(container_memory_working_set_bytes{namespace="%s", pod="%s", container="%s"}[5m])[30m:1m])`, namespace, pod, container))
}

func (p Prometheus) MaxMemory(namespace, pod, container string) (float64, error) {
	return p.queryOneVector(fmt.Sprintf(`max_over_time(rate(container_memory_working_set_bytes{namespace="%s", pod="%s", container="%s"}[5m])[30m:1m])`, namespace, pod, container))
}

func (p Prometheus) Stop() {
	p.forwarder.Stop()
}

func (p Prometheus) queryOneVector(query string) (float64, error) {
	params := url.Values{"query": []string{query}}
	b, err := p.forwarder.Get("/api/v1/query", params)
	if err != nil {
		return 0, fmt.Errorf("query prometheus vector: %w", err)
	}

	vectorResponse, err := toVectorResponse(b)
	if err != nil {
		return 0, fmt.Errorf("query prometheus vector: %w", err)
	}
	if len(vectorResponse) == 0 {
		// this is ok, for example init containers will most likely have no metrics for past hour for long-running pods
		return 0, nil
	}
	return vectorResponse[0].Value.Value, nil
}

func toVectorResponse(b []byte) ([]VectorResult, error) {
	var response QueryResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}
	if response.Data.ResultType != "vector" {
		return nil, fmt.Errorf("unexpected %s result type, expected vector", response.Data.ResultType)
	}

	var vectorResult []VectorResult
	if err := json.Unmarshal(response.Data.Result, &vectorResult); err != nil {
		return nil, fmt.Errorf("vecotr result: %w", err)
	}
	return vectorResult, nil
}

type QueryResponse struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

type QueryData struct {
	ResultType string          `json:"resultType"` // "matrix" | "vector" | "scalar" | "string"
	Result     json.RawMessage `json:"result"`
}

type VectorResult struct {
	Metric map[string]string `json:"metric"`
	Value  ResultValue       `json:"value"`
}

type ResultValue struct {
	Time  time.Time
	Value float64
}

func (r *ResultValue) UnmarshalJSON(b []byte) error {
	var v []interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf("unmarshal prometheus result value: expected slice with 2 elements, got %d", len(v))
	}
	timestamp, ok := v[0].(float64)
	if !ok {
		return fmt.Errorf("unmarshal prometheus result value: cannot convert timestap %v", v[0])
	}
	value, err := strconv.ParseFloat(v[1].(string), 64)
	if err != nil {
		return fmt.Errorf("unmarshal prometheus result value: convert value %w", err)
	}

	sec, dec := math.Modf(timestamp)
	r.Time = time.Unix(int64(sec), int64(dec*(1e9)))
	r.Value = value
	return nil
}
