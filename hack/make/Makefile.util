##@ Utility

define TODO_HELP_INFO
# Show to-do items per file.
#
# Example:
#   make todo
endef
.PHONY: todo
ifeq ($(HELP),y)
todo:
	$(Q) echo "$$TODO_HELP_INFO"
else
todo: ## Show to-do items per file
	$(Q) grep \
		--exclude-dir=vendor \
		--exclude-dir=.idea \
		--text \
		--color \
		-nRo \
		-E 'TODO.*' \
		.
endif

define VARS_HELP_INFO
# Output all variables.
#
# Example:
#   make vars
#
# See:
#   [1] https://stackoverflow.com/a/59097246
#   [2] https://www.gnu.org/software/make/manual/html_node/Origin-Function.html
endef
.PHONY: vars
ifeq ($(HELP),y)
vars:
	$(Q) echo "$$VARS_HELP_INFO"
else
vars: ## Output all variables
	$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
endif

define HELP_INFO
# Display this help.
#
# Example:
#   make help
endef
.PHONY: help
ifeq ($(HELP),y)
help:
	$(Q) echo "$$HELP_INFO"
else
help: ## Display this help
	$(Q) awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
endif
