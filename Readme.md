# cache-service

[![CircleCI](https://circleci.com/gh/mdevilliers/cache-service.svg?style=svg)](https://circleci.com/gh/mdevilliers/cache-service)
[![ReportCard](https://goreportcard.com/badge/github.com/mdevilliers/cache-service)](https://goreportcard.com/report/github.com/mdevilliers/cache-service)

POC for a simple caching service

### Local developement

To update the generated protobuf files :

```
make install_proto_tools
make proto
```

To generate all of the mocks

```
make mocks
```

### Local K8s deployment

Install Kind using the instructions [here](/hack/kind)

```
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
make hack_image_deploy_local
```


