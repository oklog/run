.PHONY: deps
deps:
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/google/yamlfmt/cmd/yamlfmt@v0.10.0
	@go mod tidy

.PHONY: check
check:
	@staticcheck -checks all

.PHONY: test
test:
	@bash -c "set -e; set -o pipefail; go test -v -cover -covermode atomic -race ./..."

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(shell find . -name "*.go")
	@yamlfmt .
