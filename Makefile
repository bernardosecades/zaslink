SHELL := /bin/bash # Use bash syntax

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Swagger UI
run-openapi-ui: ## runs swagger ui: http://localhost:4000/
	docker rm -f sharesecret-swagger 2>/dev/null || true && \
    docker run --name sharesecret-swagger -p 4000:8080 \
    --restart unless-stopped \
    -e SWAGGER_JSON=/docs/openapi/secret.yaml \
    -v $(PWD)/docs/openapi:/docs/openapi \
    swaggerapi/swagger-ui:v5.17.14

## Quality
check-quality: ## runs code quality checks
	make lint
	make fmt
	make vet

# Append || true below if blocking local developement
lint: ## go linting. Update and use specific lint tool and options
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	golangci-lint run -c ./.golangci.yml

lint-fix: ## go linting. Update and use specific lint tool and options
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	golangci-lint run -c ./.golangci.yml --fix

vet: ## go vet
	go vet ./...

fmt: ## runs go formatter
	go fmt ./...

tidy: ## run go mod tidy
	go mod tidy
## Test
test-all: ## runs tests and create generates coverage report
	make tidy
	go test ./... --tags=unit,integration -coverprofile=coverage.out

test-integration:
	make tidy
	go test ./... --tags=integration

test-unit:
	make tidy
	go test ./... --tags=unit

coverage: ## displays test coverage report in html mode
	make test-all
	go tool cover -html=coverage.out

.PHONY: all test-all
## All
all: ## quality checks and tests
	make check-quality
	make test-all

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
