# cache-service

[![CircleCI](https://circleci.com/gh/mdevilliers/cache-service.svg?style=svg)](https://circleci.com/gh/mdevilliers/cache-service)

POC of a caching service with a GRPC API.

### Local developement

To update the generated protobuf files :

```
make install_proto_tools
make proto
```

To generate the mocks

```
make mocks
```

### Local K8s deployment

Install Kind using the instructions [here](/hack/kind)

```
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
make hack_image_deploy_local
```
