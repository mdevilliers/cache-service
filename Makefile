CMD      := cache-service
PKG      := github.com/mdevilliers/cache-service
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

EXE_NAME := cache-service

# Versioning
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_SHA    ?= $(shell git rev-parse --short HEAD)
GIT_TAG    ?= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  ?= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

# Binary Name
BIN_OUTDIR ?= ./build/bin
BIN_NAME   ?= cache-service-$(shell go env GOOS)-$(shell go env GOARCH)
ifeq ("$(GIT_TAG)","")
	BIN_VERSION = $(GIT_SHA)
endif
BIN_VERSION ?= ${GIT_TAG}


DOCKER_REGISTRY := mdevilliers

# Docker Tag from Git
DOCKER_IMAGE_TAG  ?= ${GIT_TAG}_$(GIT_SHA)_$(GIT_DIRTY)
DOCKER_BUILD_CMD:=
DOCKER_BUILD_CMD+= $(GO_BUILD_VARS) $(GO_BUILD) $(GO_BUILD_FLAGS)
DOCKER_BUILD_CMD+= -o docker/$(EXE_NAME)
DOCKER_BUILD_CMD+= github.com/mdevilliers/cache-service/cmd/cache-service

# LDFlags
LDFLAGS += -X $(PKG)/internal/version.Timestamp=$(shell date +%s)
LDFLAGS += -X $(PKG)/internal/version.GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(PKG)/internal/version.GitTreeState=${GIT_DIRTY}
LDFLAGS += -X $(PKG)/internal/version.Version=${BIN_VERSION}

# CGO
CGO ?= 1

# Go Build Flags
GOBUILDFLAGS :=
GOBUILDFLAGS += -o $(BIN_OUTDIR)/$(BIN_NAME)

# Linting
OS := $(shell uname)
GOLANGCI_LINT_VERSION=1.37.1
ifeq ($(OS),Darwin)
	GOLANGCI_LINT_ARCHIVE=golangci-lint-$(GOLANGCI_LINT_VERSION)-darwin-amd64.tar.gz
else
	GOLANGCI_LINT_ARCHIVE=golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz
endif

.PHONY: info
info:
	@echo "Version:        ${BIN_VERSION}"
	@echo "Binary Name:    ${BIN_NAME}"
	@echo "Git Tag:        ${GIT_TAG}"
	@echo "Git Commit:     ${GIT_COMMIT}"
	@echo "Git Tree State: ${GIT_DIRTY}"

# Build a statically linked binary
.PHONY: static
static: CGO = 0
static: build

# Build a binary
.PHONY: build
build: CMD = ./cmd/cache-service
build: GOBUILDFLAGS += -ldflags '$(LDFLAGS)'
build:
	@CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) $(CMD)

DOCKER_ARGS:=
DOCKER_ARGS+= --force-rm
DOCKER_ARGS+= --build-arg BIN_VERSION=$(BIN_VERSION)
DOCKER_ARGS+= --build-arg GIT_COMMIT=$(GIT_COMMIT)
DOCKER_ARGS+= --build-arg GIT_SHA=$(GIT_SHA)
DOCKER_ARGS+= --build-arg GIT_TAG=$(GIT_TAG)
DOCKER_ARGS+= --build-arg GIT_DIRTY=$(GIT_DIRTY)
DOCKER_ARGS+= -f ./build/package/Dockerfile
DOCKER_ARGS+= -t $(DOCKER_REGISTRY)/$(CMD):$(DOCKER_IMAGE_TAG)

# Build docker image
.PHONY: image
image:
	docker build $(DOCKER_ARGS) .

.PHONY: test
# Run test suite
test:
ifeq ("$(wildcard $(shell which gocov))","")
	go get github.com/axw/gocov/gocov
endif
	gocov test ${PKG_LIST} | gocov report

# deploys to configured kubernetes instance
.PHONY: deploy
deploy:
	kubectl delete -f k8s/ 2>/dev/null; true
	kubectl create -f k8s/

# the linting gods must be obeyed
.PHONY: lint
lint: $(BIN_OUTDIR)/golangci-lint/golangci-lint
	$(BIN_OUTDIR)/golangci-lint/golangci-lint run

$(BIN_OUTDIR)/golangci-lint/golangci-lint:
	curl -OL https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/$(GOLANGCI_LINT_ARCHIVE)
	mkdir -p $(BIN_OUTDIR)/golangci-lint/
	tar -xf $(GOLANGCI_LINT_ARCHIVE) --strip-components=1 -C $(BIN_OUTDIR)/golangci-lint/
	chmod +x $(BIN_OUTDIR)/golangci-lint
	rm -f $(GOLANGCI_LINT_ARCHIVE)

.PHONY: install_proto_tools
# install a known version of the protoc compiler
# TODO : deploy to a Docker registry and each developer can download the same copy.
install_proto_tools:
	docker build -t $(DOCKER_REGISTRY)/proto_tools -f ./build/package/Dockerfile.proto .

.PHONY: proto
# regenerate protobuf files
proto:
	docker run -v $(PWD)/proto:/go/proto $(DOCKER_REGISTRY)/proto_tools

.PHONY: proto_verify
# verify proto binding has been generated
# The CI will check that no un-generated files have been checked in
# TODO : add stages in CI to check
proto_verify: proto
	git diff --exit-code

.PHONY: mocks
# generate mocks
mocks:
ifeq ("$(wildcard $(shell which counterfeiter))","")
	go get github.com/maxbrunsfeld/counterfeiter/v6
endif
	counterfeiter -o=./proto/v1/mocks/service.go ./proto/v1/service.pb.go CacheClient
	counterfeiter -o=./internal/service/mocks/store.go ./internal/service/service.go store

.PHONY: hack_image_deploy_local
# task to deploy and build a local image using a `kind` environment
# see ./hack/kind/ for details.
# This task has the following steps :
# - build the application and docker image locally
# - deploy the image to the 'kind' cluster
# - set the image in the deployment pod to the latest value
# - force the deployment to redeploy by changing some metadata
hack_image_deploy_local: image deploy
	kind load docker-image $(DOCKER_REGISTRY)/$(CMD):$(DOCKER_IMAGE_TAG)
	kubectl set image deployment/cache-service-deployment cache-service=$(DOCKER_REGISTRY)/cache-service:$(DOCKER_IMAGE_TAG)
	kubectl patch -f ./k8s/cache-service_deployment.yaml -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"

.PHONY: hack_local_redis
# task to run a local redis instace for development
hack_local_redis:
	docker run --rm -p 6379:6379 -d redis:3.2.9
