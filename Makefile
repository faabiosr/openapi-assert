.DEFAULT_GOAL := test

# Clean up
clean:
	@rm -fR ./vendor/ ./coverage.*
.PHONY: clean

# Download project dependencies
configure:
	dep ensure -v
.PHONY: configure

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./coverage.text -o ./coverage.html
.PHONY: cover

# Format all go files
fmt:
	@gofmt -s -w -l $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)
.PHONY: fmt

# Run tests
test:
	@go test -v -race -coverprofile=./coverage.text -covermode=atomic $(shell go list ./...)
.PHONY: test
