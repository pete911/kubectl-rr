# kubectl-rr
Kubernetes' resource recommendation. Command line tool to list pods, output containers requests and limits and also
usage.

Prerequisites, prometheus has to be installed in the cluster.

```
namespace kube-system pod coredns-5d78c9869d-drmv4
  container coredns
    cpu requests 100m limits 0 current: 0.0024475560790870273 min: 0.0009372885548091062 max: 0.002529520302870373
namespace kube-system pod etcd-rr-test-control-plane
  container etcd
    cpu requests 100m limits 0 current: 0.026811003644753445 min: 0.01103471804517143 max: 0.02824301838711446
namespace kube-system pod kindnet-8tz7r
  container kindnet-cni
    cpu requests 100m limits 100m current: 0.0008620287157887875 min: 0.0003609862491599064 max: 0.0008620287157887875
namespace kube-system pod kube-apiserver-rr-test-control-plane
  container kube-apiserver
    cpu requests 250m limits 0 current: 0.05591350447152222 min: 0.022103391698499973 max: 0.05722482040961385
namespace kube-system pod kube-controller-manager-rr-test-control-plane
  container kube-controller-manager
    cpu requests 200m limits 0 current: 0.018073922814116995 min: 0.00727293128565153 max: 0.019075219536946856
```
