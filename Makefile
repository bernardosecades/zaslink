APP_NAME=api-share-secret
APP_EXECUTABLE=./out/${APP_NAME}
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
SHELL := /bin/bash # Use bash syntax

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Quality
check-quality: ## runs code quality checks
	make lint
	make fmt
	make vet

# Append || true below if blocking local developement
lint: ## go linting. Update and use specific lint tool and options
	golangci-lint run

lint-fix: ## go linting. Update and use specific lint tool and options
	golangci-lint run --fix

vet: ## go vet
	go vet ./...

fmt: ## runs go formatter
	go fmt ./...

## Test
test: ## runs tests and create generates coverage report
	make tidy
	go test -v -timeout 10m ./... -coverprofile=coverage.out -json > report.json

coverage: ## displays test coverage report in html mode
	make test
	go tool cover -html=coverage.out

.PHONY: all test build vendor
## All
all: ## runs setup, quality checks and builds
	make check-quality
	make test
	make build

.PHONY: help
## Help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
