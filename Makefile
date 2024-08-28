# Define packages path
PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

.PHONY: all require tidy fmt goimports vet staticcheck govulncheck test

all: require tidy fmt goimports vet staticcheck govulncheck test

require:
	@type "goimports" > /dev/null 2>&1 || (echo 'goimports not found: to install it, run "go install golang.org/x/tools/cmd/goimports@latest"'; exit 1)
	@type "staticcheck" > /dev/null 2>&1 || (echo 'staticcheck not found: to install it, run "go install honnef.co/go/tools/cmd/staticcheck@latest"'; exit 1)
	@type "govulncheck" > /dev/null 2>&1 || (echo 'govulncheck not found: to install it, run "go install golang.org/x/vuln/cmd/govulncheck@latest"'; exit 1)
	@type "gocyclo" > /dev/null 2>&1 || (echo 'gocyclo not found: to install it, run "go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"'; exit 1)
	@type "golangci-lint" > /dev/null 2>&1 || (echo 'golangci-lint not found: to install it, run "brew install golangci-lint"'; exit 1)

tidy:
	@echo "=> Executing go mod tidy"
	@go mod tidy

fmt:
	@echo "=> Executing go fmt"
	@go fmt ./...

goimports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

vet:
	@echo "=> Executing go vet"
	@go vet ./...

staticcheck:
	@echo "=> Executing staticcheck"
	@staticcheck ./...

govulncheck:
	@echo "=> Executing govulncheck"
	@govulncheck ./... || true

test:
	@go test -v ./... -coverprofile=coverage.out
