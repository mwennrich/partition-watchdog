# partition-watchdog

***ALPHA version***

## Usage

```text
Usage:
  partition-watchdog [command]

Available Commands:
  check       check connectivity to a partition and stop kube-controller-manager if unavailable
  help        Help about any command

Flags:
      --checkinterval duration   time between nodeReady checks (default 10s)
      --deployment string        name of deployment to scale (default "kube-controller-manager")
  -h, --help                     help for partition-watchdog
      --target string            target to check (e.g. 212.1.5.7:80
      --timeout duration         connection timeout for checks (default 2s)
      --tries int                number of checks partition is considered (un)available (default 5)
```

## Example

```bash
kubectl apply -f deploy/partition-watchdog.yaml
```
