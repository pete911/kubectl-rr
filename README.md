# kubectl-rr
Kubernetes' resource recommendation. Command line tool to list pods, output containers requests and limits and also
usage.

Prerequisites, prometheus has to be installed in the cluster.

```
kubectl-rr pod -n kube-system

NAMESPACE    POD                             CONTAINER                CPU/R  CPU/L  CPU    CPU/Min  CPU/Max  Memory/R  Memory/L  Memory   Memory/Min  Memory/Max
kube-system  aws-node-tbfqw                  aws-vpc-cni-init [init]  25m    0      -      -        -        0         0         -        -           -
kube-system  aws-node-tbfqw                  aws-node                 25m    0      2.80m  2.73m    2.89m    0         0         227.2Ki  66.7        451.8Ki
kube-system  coredns-d6f79c549-2zlbw         coredns                  100m   0      1.13m  1.08m    1.14m    70Mi      170Mi     4.3Ki    49.7        133.0Ki
kube-system  coredns-d6f79c549-btjkj         coredns                  100m   0      1.12m  0.91m    1.12m    70Mi      170Mi     64.9Ki   -           134.5Ki
kube-system  kube-proxy-pr8nz                kube-proxy               100m   0      0.20m  0.18m    0.24m    0         0         61.0Ki   -           137.1Ki
kube-system  kube-proxy-szxrk                kube-proxy               100m   0      0.23m  0.20m    0.25m    0         0         66.0Ki   -           134.5Ki
kube-system  metrics-server-fbb469ccc-jzqn7  metrics-server           100m   0      2.10m  1.82m    2.19m    200Mi     0         78.1Ki   -           165.0Ki
```
