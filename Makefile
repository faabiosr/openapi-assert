.DEFAULT_GOAL := test

# Clean up
clean:
	@rm -fR ./vendor/ ./cover.*
.PHONY: clean

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./cover.out -o ./cover.html
	@test -f ./cover.out && rm ./cover.out;
.PHONY: cover

# Format all go files
fmt:
	@gofmt -s -w -l $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)
.PHONY: fmt

# GolangCI Linter
lint:
	@golangci-lint run -v ./...
.PHONY: lint

# Run tests
test:
	@go test -v -race -coverprofile=./cover.out -covermode=atomic ./...
.PHONY: test
