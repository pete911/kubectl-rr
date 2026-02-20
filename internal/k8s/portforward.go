package k8s

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type Forwarder struct {
	portForwarder *portforward.PortForwarder
	stopChan      chan struct{}
	host          string
	httpClient    *http.Client
}

func StartForwarder(restConfig *rest.Config, namespace, podName, podPort string) (Forwarder, error) {
	serverURL, err := url.Parse(restConfig.Host)
	if err != nil {
		return Forwarder{}, fmt.Errorf("parse rest config host: %w", err)
	}
	serverURL.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)

	transport, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return Forwarder{}, err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf(":%s", podPort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return Forwarder{}, err
	}

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			fmt.Printf("forward ports: %v\n", err)
		}
	}()

	<-readyChan
	ports, err := forwarder.GetPorts()
	if err != nil {
		return Forwarder{}, fmt.Errorf("get ports: %w", err)
	}
	if len(ports) != 1 {
		return Forwarder{}, fmt.Errorf("returned %d ports, expected 1", err)
	}

	return Forwarder{
		portForwarder: forwarder,
		stopChan:      stopChan,
		host:          fmt.Sprintf("http://localhost:%d", ports[0].Local),
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (f Forwarder) Get(path string, params url.Values) ([]byte, error) {
	path, err := url.JoinPath(f.host, path)
	if err != nil {
		return nil, fmt.Errorf("join request path: %w", err)
	}
	if len(params) != 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	resp, err := f.httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%s unexpected %s status", resp.Status, path)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return b, nil
}

func (f Forwarder) Stop() {
	f.stopChan <- struct{}{}
}
