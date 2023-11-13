# kubectl-rr
Kubernetes' resource recommendation. Command line tool to list pods, output containers requests and limits and also
usage.

Prerequisites, prometheus has to be installed in the cluster.

```
kubectl-rr pod -n kube-system

NAMESPACE    POD                             CONTAINER                CPU/R  CPU/L  CPU    CPU/Min  CPU/Max
kube-system  aws-node-tbfqw                  aws-vpc-cni-init [init]  25m    0      -      -        -
kube-system  aws-node-tbfqw                  aws-node                 25m    0      2.44m  2.70m    2.86m
kube-system  aws-node-tbfqw                  aws-eks-nodeagent        25m    0      0.36m  0.35m    0.36m
kube-system  aws-node-zn9sk                  aws-vpc-cni-init [init]  25m    0      -      -        -
kube-system  aws-node-zn9sk                  aws-node                 25m    0      2.61m  2.22m    2.72m
kube-system  aws-node-zn9sk                  aws-eks-nodeagent        25m    0      0.32m  0.33m    0.34m
kube-system  coredns-d6f79c549-2zlbw         coredns                  100m   0      1.00m  1.09m    1.15m
kube-system  coredns-d6f79c549-btjkj         coredns                  100m   0      1.09m  0.92m    1.11m
kube-system  kube-proxy-pr8nz                kube-proxy               100m   0      0.21m  0.19m    0.23m
kube-system  kube-proxy-szxrk                kube-proxy               100m   0      0.20m  0.21m    0.24m
kube-system  metrics-server-fbb469ccc-jzqn7  metrics-server           100m   0      2.13m  2.08m    2.21m
```
