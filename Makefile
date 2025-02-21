# Used for temporary files such as coverage files.
TMP_DIR ?= .tmp

COVERAGE_FILE := $(TMP_DIR)/coverage.txt
COVERAGE_HTML_FILE := $(TMP_DIR)/coverage.html

# golangci-lint version to use on this project.
GOLANGCI_VERSION ?= v1.64.5
GOLANGCI := $(TMP_DIR)/golangci-lint

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
ifeq ($(UNAME_ARCH),x86_64)
ifeq ($(UNAME_OS),Darwin)
OPEN_CMD := open
endif
ifeq ($(UNAME_OS),Linux)
OPEN_CMD := xdg-open
endif
endif

.DEFAULT_GOAL := help
.PHONY: test
test: ## Run all the tests
	@echo "--> Running tests..."
	@mkdir -p $(dir $(COVERAGE_FILE))
	@go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) -race -failfast -timeout=30s ./...

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@echo "--> Creating HTML coverage report at $(COVERAGE_HTML_FILE)..."
	@mkdir -p $(dir $(COVERAGE_FILE)) $(dir $(COVERAGE_HTML_FILE))
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_FILE)
ifndef COVEROPEN
	@echo "--> Open HTML coverage report: $(OPEN_CMD) $(COVERAGE_HTML_FILE)"
else
	$(OPEN_CMD) $(COVERAGE_HTML_FILE)
endif

build: clean ## Build the app
	@echo "--> Building..."
	@go build ./...

.PHONY: clean
clean: ## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -f $(GOLANGCI) $(COVERAGE_FILE) $(COVERAGE_HTML_FILE)
	@go clean
	@go mod tidy -v

$(GOLANGCI):
	@echo "--> Installing golangci $(GOLANGCI_VERSION)..."
	@mkdir -p $(dir $(GOLANGCI))
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $(GOLANGCI)) $(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI) ## Run linter
	@echo "--> Running linter golangci $(GOLANGCI_VERSION)..."
	@$(GOLANGCI) run

.PHONY: ci
ci: lint test cover ## Run all the tests and code checks

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
