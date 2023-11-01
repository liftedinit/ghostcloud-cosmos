#### LINT ####

golangci_version=v1.55.1

lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@echo "--> Installing golangci-lint $(golangci_version) complete"

lint:
	@echo "--> Running linter"
	$(MAKE) lint-install
	@golangci-lint run ./x/...

lint-fix:
	@echo "--> Running linter"
	$(MAKE) lint-install
	@golangci-lint run ./x/... --fix

.PHONY: lint lint-fix

#### FORMAT ####

goimports_version=latest

format-install:
	@echo "--> Installing goimports $(goimports_version)"
	@go install golang.org/x/tools/cmd/goimports@$(goimports_version)
	@echo "--> Installing goimports $(goimports_version) complete"

format:
	@echo "--> Running goimports"
	$(MAKE) format-install
	@find . -name '*.go' -exec goimports -w -local github.com/cosmos/cosmos-sdk,cosmossdk.io,github.com/cometbft,github.com/cosmos.ibc-go,ghostcloud  {} \;

#### COVERAGE ####

coverage:
	@echo "--> Running coverage"
	@go test -race -cpu=$$(nproc) -covermode=atomic -coverprofile=coverage.out $$(go list ./x/...) > /dev/null 2>&1
	@echo "--> Running coverage filter"
	@./scripts/filter-coverage.sh
	@echo "--> Running coverage report"
	@go tool cover -func=coverage-filtered.out
	@echo "--> Running coverage html"
	@go tool cover -html=coverage-filtered.out -o coverage.html
	@echo "--> Coverage report available at coverage.html"
	@echo "--> Cleaning up coverage files"
	@rm coverage.out
	@rm coverage-filtered.out
	@echo "--> Running coverage complete"

.PHONY: coverage