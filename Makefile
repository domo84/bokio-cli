BINARY      := bokio
SPEC_DIR    := specs
SPEC_REPO   := https://raw.githubusercontent.com/bokio/bokio-api
SPEC_REF    ?= v1
SPEC_PATH   := api-specification
SPECS       := company-api.yaml general-api.yaml
SPEC_FILES  := $(addprefix $(SPEC_DIR)/,$(SPECS))

.DEFAULT_GOAL := help

## build: compile the CLI binary
.PHONY: build
build:
	go build -o $(BINARY) ./cmd/bokio/

## vet: run static analysis
.PHONY: vet
vet:
	go vet ./...

## test: run the test suite
.PHONY: test
test:
	go test ./...

## specs: download the Bokio OpenAPI specs into specs/ (not committed)
.PHONY: specs
specs:
	@mkdir -p $(SPEC_DIR)
	@for spec in $(SPECS); do \
		echo "Downloading $$spec ($(SPEC_REF))"; \
		curl -fsSL -o $(SPEC_DIR)/$$spec $(SPEC_REPO)/$(SPEC_REF)/$(SPEC_PATH)/$$spec || exit 1; \
	done
	@echo "Specs written to $(SPEC_DIR)/"

## specs-list: list the operations defined by the downloaded specs
.PHONY: specs-list
specs-list: $(SPEC_FILES)
	@for spec in $(SPEC_FILES); do \
		echo "== $$spec"; \
		awk '/^  \//{path=substr($$1,1,length($$1)-1)} /^    (get|post|put|patch|delete):/{sub(":","",$$1); printf "%s %s\n", toupper($$1), path}' $$spec; \
	done

$(SPEC_FILES):
	@$(MAKE) --no-print-directory specs

## specs-clean: remove the downloaded specs
.PHONY: specs-clean
specs-clean:
	rm -rf $(SPEC_DIR)

## clean: remove build artifacts and downloaded specs
.PHONY: clean
clean: specs-clean
	rm -f $(BINARY)

## help: list available targets
.PHONY: help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
