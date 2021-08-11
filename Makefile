.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

# Utilities.
V := 0
Q := $(if $(filter 1,$(V)),,@)
T := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set).
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Image spec to use all building/pushing image targets.
IMAGE_REPOSITORY ?= spotinst
IMAGE_NAME ?= ocean-operator
IMAGE_TAG ?= latest
IMAGE_SPEC := $(IMAGE_REPOSITORY)/$(IMAGE_NAME):$(IMAGE_SPEC)

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

##@ General

.PHONY: help
help: ## Display this help
	$(Q) awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(Q) $(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(Q) $(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(Q) go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(Q) go vet ./...

.PHONY: test
test: manifests generate fmt vet setup-envtest ## Run tests.
	$(Q) go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	$(Q) go build -trimpath -o bin/ocean-operator cmd/ocean-operator/main.go

.PHONY: build-tide
build-tide: generate fmt vet ## Build tide binary.
	$(Q) go build -trimpath -o bin/ocean-tide cmd/ocean-tide/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	$(Q) go run cmd/ocean-operator/main.go manager $(ARGS)

BUILD_IMAGE_TARGETS := build-image-amd build-image-arm
.PHONY: $(BUILD_IMAGE_TARGETS) build-image
$(BUILD_IMAGE_TARGETS): build-image
build-image: ## Build and push Docker image for all platforms
	$(Q) hack/scripts/build.sh
build-image-amd: TARGET_PLATFORM="linux/amd64" ## Build docker image for linux/amd64
build-image-arm: TARGET_PLATFORM="linux/arm64" ## Build docker image for linux/arm64

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install CRDs into the cluster specified in ~/.kube/config.
	$(Q) $(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the cluster specified in ~/.kube/config.
	$(Q) $(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the cluster specified in ~/.kube/config.
	$(Q) cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMAGE_SPEC)
	$(Q) $(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the cluster specified in ~/.kube/config.
	$(Q) $(KUSTOMIZE) build config/default | kubectl delete -f -

.PHONY: controller-gen
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(Q) $(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

.PHONY: kustomize
KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(Q) $(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

.PHONY: setup-envtest
SETUP_ENVTEST = $(shell pwd)/bin/setup-envtest
setup-envtest: ## Download setup-envtest locally if necessary.
	$(Q) $(call go-get-tool,$(SETUP_ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)
	$(Q) $(SETUP_ENVTEST) use -i -p env 1.21.x!
	$(Q) source <($(SETUP_ENVTEST) -i -p env 1.21.x!)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
