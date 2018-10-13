# OpenAPI - Assert

[![Build Status](https://img.shields.io/travis/fabiorphp/openapi-assert/master.svg?style=flat-square)](https://travis-ci.org/fabiorphp/openapi-assert)
[![Codecov branch](https://img.shields.io/codecov/c/github/fabiorphp/openapi-assert/master.svg?style=flat-square)](https://codecov.io/gh/fabiorphp/openapi-assert)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/fabiorphp/openapi-assert)
[![Go Report Card](https://goreportcard.com/badge/github.com/fabiorphp/openapi-assert?style=flat-square)](https://goreportcard.com/report/github.com/fabiorphp/openapi-assert)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/fabiorphp/openapi-assert/blob/master/LICENSE)

Asserting data against OpenAPI docs. This project is inspired by [PHP Swagger Assertions](https://github.com/Maks3w/SwaggerAssertions).

## Instalation

OpenAPI Assert requires Go 1.11 or later.

```sh
$ go get github.com/fabiorphp/openapi-assert
```

If you want to get an specific version, please use the example below:

```
go get gopkg.in/fabiorphp/openapi-assert.v0
```

## Development

### Requirements

- Install [Go](https://golang.org)

### Makefile
```sh
# Clean up
$ make clean

# Download project dependencies
$ make configure

# Run tests and generates html coverage file
$ make cover

# Format all go files
$ make fmt

# Run tests
$make test
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/fabiorphp/openapi-assert/blob/master/LICENSE) for more details.
