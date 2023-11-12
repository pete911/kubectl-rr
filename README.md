# kubectl-rr
Kubernetes' resource recommendation. Command line tool to list pods, output containers requests and limits and also
usage.

Prerequisites, prometheus has to be installed in the cluster.

```
namespace prometheus pod prometheus-kube-state-metrics-56f5765bcf-qjx66
containers
  kube-state-metrics
    cpu requests 0 limits 0 current: 0.00146 min: 0.00118 max: 0.00183
namespace prometheus pod prometheus-prometheus-node-exporter-6bgz4
containers
  node-exporter
    cpu requests 0 limits 0 current: 0.00108 min: 0.00072 max: 0.00129
namespace prometheus pod prometheus-prometheus-pushgateway-5b7b9f67bb-f59t2
containers
  pushgateway
    cpu requests 0 limits 0 current: 0.00064 min: 0.00040 max: 0.00076
namespace prometheus pod prometheus-server-7bbd49dd-4dtnc
containers
  prometheus-server-configmap-reload
    cpu requests 0 limits 0 current: 0.00015 min: 0.00006 max: 0.00014
  prometheus-server
    cpu requests 0 limits 0 current: 0.00619 min: 0.00541 max: 0.01269
```
