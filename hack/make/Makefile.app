##@ Application

define RUN_HELP_INFO
# Run the operator locally.
#
# Example:
#   make run
endef
.PHONY: run
ifeq ($(HELP),y)
run:
	$(Q) echo "$$RUN_HELP_INFO"
else
run: ## Run the operator locally
	$(Q) OPERATOR_NAME=$(OPERATOR_NAME) $(OPERATORSDK) run \
		--go-ldflags "$(GO_LDFLAGS)" \
		--operator-flags "$(OPERATOR_FLAGS)" \
		--local
endif

define BUILD_HELP_INFO
# Compile code and build artifacts.
#
# See:
#   [1] https://github.com/operator-framework/operator-sdk/issues/1867
#
# Example:
#   make build
endef
.PHONY: build
ifeq ($(HELP),y)
build:
	$(Q) echo "$$BUILD_HELP_INFO"
else
build: ## Compile code and build artifacts
	$(Q) $(call goreleaser_release,build,--snapshot)
	$(Q) $(CONTAINER_ENGINE) build -f build/Dockerfile -t $(OPERATOR_IMAGE) .
endif

define GEN_HELP_INFO
# Generate Kubernetes code for CRs and/or CRDs for API's.
#
# Example:
#   make gen-code
#   make gen-crds
endef
GEN_TARGETS := gen-code gen-crds
.PHONY: $(GEN_TARGETS) gen
ifeq ($(HELP),y)
$(GEN_TARGETS) gen:
	$(Q) echo "$$GEN_HELP_INFO"
else
$(GEN_TARGETS): TYPE=$(MAKECMDGOALS:gen-%=%)
$(GEN_TARGETS): gen
gen: tools
	$(Q) if [[ "$(TYPE)" == "code" ]] ; then $(GOGEN) -v $(PKGS); fi
	$(Q) if [[ ! -z "$(strip $(GENERATOR))" ]] ; then $(OPERATORSDK) generate $(strip $(GENERATOR)); fi
gen-code: GENERATOR=k8s  ## Generates Kubernetes code for CRs
gen-crds: GENERATOR=crds ## Generates CRDs for API's
endif

define INSTALL_HELP_INFO
# Install all resources (CRDs, RBAC and Operator).
#
# Example:
#   make install
#   make install INCLUSTER=y
endef
.PHONY: install
ifeq ($(HELP),y)
install:
	$(Q) echo "$$INSTALL_HELP_INFO"
else
install: ## Install all resources (CRDs, RBAC and Operator)
	@echo "===> Creating namespace"
	$(Q) $(KUBECTL) create namespace $(NAMESPACE) || true

	@echo "===> Applying CRDs"
	$(Q) find $(OPERATOR_DEPLOY_DIR)/crds/ -type f -name "*crd.yaml" \
		-exec $(KUBECTL) apply -n $(NAMESPACE) -f {} \;

ifeq ($(INCLUSTER),y)
	@echo "===> Applying RBAC and service account"
	$(Q) $(KUBECTL) apply -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/role.yaml
	$(Q) $(KUBECTL) apply -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/role_binding.yaml
	$(Q) $(KUBECTL) apply -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/service_account.yaml

	@echo "===> Applying operator"
	$(Q) $(KUBECTL) apply -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/operator.yaml
endif
endif

define UNINSTALL_HELP_INFO
# Uninstall all resources (CRDs, RBAC and Operator).
#
# Example:
#   make uninstall
#   make uninstall INCLUSTER=y
endef
.PHONY: uninstall
ifeq ($(HELP),y)
uninstall:
	$(Q) echo "$$UNINSTALL_HELP_INFO"
else
uninstall: ## Uninstall all resources (CRDs, RBAC and Operator)
	@echo "===> Deleting CRDs"
	$(Q) find $(OPERATOR_DEPLOY_DIR)/crds/ -type f -name "*crd.yaml" \
		-exec $(KUBECTL) delete --ignore-not-found -n $(NAMESPACE) -f {} \;

ifeq ($(INCLUSTER),y)
	@echo "===> Deleting RBAC and service account"
	$(Q) $(KUBECTL) delete --ignore-not-found -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/role.yaml
	$(Q) $(KUBECTL) delete --ignore-not-found -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/role_binding.yaml
	$(Q) $(KUBECTL) delete --ignore-not-found -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/service_account.yaml

	@echo "===> Deleting operator"
	$(Q) $(KUBECTL) delete --ignore-not-found -n $(NAMESPACE) -f $(OPERATOR_DEPLOY_DIR)/operator.yaml
endif
	@echo "===> Deleting namespace"
	$(Q) $(KUBECTL) delete --ignore-not-found namespace $(NAMESPACE)
endif

define BUNDLE_HELP_INFO
# Build the Operator metadata bundle.

# See:
#   [1] https://operatorhub.io/bundle
#   [2] https://redhat-connect.gitbook.io/certified-operator-guide
#
# Example:
#   make bundle
endef
.PHONY: bundle
ifeq ($(HELP),y)
bundle:
	$(Q) echo "$$BUNDLE_HELP_INFO"
else
bundle: ## Build the Operator bundle
	@echo "===> Cleaning up the existing bundle"
	$(Q) rm -rf $(OPERATOR_BUNDLE_DIR) && mkdir -p $(OPERATOR_BUNDLE_DIR)

	@echo "===> Generating a CSV file for the Operator"
	$(Q) $(OPERATORSDK) generate csv --operator-name $(OPERATOR_NAME) --csv-version $(VERSION)

	@echo "===> Copying the Operator package and CSV files"
	$(Q) cp $(OPERATOR_DEPLOY_DIR)/olm-catalog/$(OPERATOR_NAME)/$(OPERATOR_NAME).package.yaml $(OPERATOR_BUNDLE_DIR)
	$(Q) cp $(OPERATOR_DEPLOY_DIR)/olm-catalog/$(OPERATOR_NAME)/$(VERSION)/*.yaml $(OPERATOR_BUNDLE_DIR)

	@echo "===> Copying CRDs"
	$(Q) cp $(OPERATOR_DEPLOY_DIR)/crds/ocean.spot.io_clusters_crd.yaml \
		$(OPERATOR_BUNDLE_DIR)/$(OPERATOR_NAME).ocean.spot.io_clusters.crd.yaml
	$(Q) cp $(OPERATOR_DEPLOY_DIR)/crds/ocean.spot.io_launchspecs_crd.yaml \
		$(OPERATOR_BUNDLE_DIR)/$(OPERATOR_NAME).ocean.spot.io_launchspecs.crd.yaml

	@echo "===> Updating current CSV version to $(VERSION)"
	$(Q) sed -i -E "s/currentCSV: .*/currentCSV: $(OPERATOR_NAME).v$(VERSION)/" \
		$(OPERATOR_BUNDLE_DIR)/$(OPERATOR_NAME).package.yaml

	@echo "===> NOTE: Please verify the Operator CSV and build the metadata bundle:"
	@echo "$$ cd bundle/ && zip $(OPERATOR_NAME)-metadata *.yaml && unzip -l $(OPERATOR_NAME)-metadata.zip"
endif
