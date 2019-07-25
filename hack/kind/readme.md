# Kind

Documentation for installation here - https://kind.sigs.k8s.io/

To create a local k8s cluster using the configuration file - 
```
kind create cluster --config ./kind.yaml
```
To configure your kubectl to point at the local installation
```
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")" 
```
