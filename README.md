# OpenAPI - Assert

[![Build Status](https://img.shields.io/travis/faabiosr/openapi-assert/master.svg?style=flat-square)](https://travis-ci.org/faabiosr/openapi-assert)
[![Codecov branch](https://img.shields.io/codecov/c/github/faabiosr/openapi-assert/master.svg?style=flat-square)](https://codecov.io/gh/faabiosr/openapi-assert)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/faabiosr/openapi-assert)
[![Go Report Card](https://goreportcard.com/badge/github.com/faabiosr/openapi-assert?style=flat-square)](https://goreportcard.com/report/github.com/faabiosr/openapi-assert)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/faabiosr/openapi-assert/blob/master/LICENSE)

## Description

openapi-assert is a Go package that provides a affordable way to validate http requests and responses data throught OpenAPI Schema Specification (Swagger) and the project was inspired by [PHP Swagger Assertions](https://github.com/Maks3w/SwaggerAssertions). It has the following features:

* Assert request and response media types
* Assert request and response headers
* Assert request query strings
* Assert request and response body.
* Assert the entire http request and response object.

## Requirements
OpenAPI Assert requires Go 1.11 or later.

## Instalation

Use go get.
```sh
$ go get github.com/faabiosr/openapi-assert
```

Then import the package into your own code:
```
import "github.com/faabiosr/openapi-assert"
```

## Usage
The package provides methods that allow you to assert raw data using swagger files.

See it in action:

```go
package main

import (
    assert "github.com/faabiosr/openapi-assert"
    "log"
    "net/http"
)

func main() {
    doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")

    if err != nil {
        log.Fatal(err)
    }

    log.Println(
        assert.RequestMediaType("text/html", doc, "/pet", http.MethodPost),
    )
}
```

If you want to assert data many times, it is recommended to create instance of assert:

```go
package main

import (
    assert "github.com/faabiosr/openapi-assert"
    "log"
    "net/http"
)

func main() {
    doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")

    if err != nil {
        log.Fatal(err)
    }

    assert := assert.New(doc)

    log.Println(
        assert.RequestMediaType("text/html", "/pet", http.MethodPost),
    )

    log.Println(
        assert.RequestMediaType("image/gif", "/v2/pet", http.MethodPost),
    )
}
```

Asserting http request object using the swagger schema file:

```go
package main

import (
	"fmt"
	assert "github.com/faabiosr/openapi-assert"
	"log"
	"net/http"
)

func main() {
	doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")

	if err != nil {
		log.Fatal(err)
	}

	assert := assert.New(doc)

	http.HandleFunc("/v2/pet", func(w http.ResponseWriter, r *http.Request) {
		err := assert.Request(r)

		fmt.Fprint(w, err)
	})

	log.Fatal(
		http.ListenAndServe("127.0.0.1:9000", nil),
	)
}
```

Asserting http request object using the swagger schema file:

```go
package main

import (
	assert "github.com/faabiosr/openapi-assert"
	"log"
	"net/http"
)

func main() {
	doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")

	if err != nil {
		log.Fatal(err)
	}

	assert := assert.New(doc)

	res, err := http.Get("https://petstore.swagger.io/v2/pet/111111422")

	if err != nil {
		log.Fatal(err)
	}

	log.Println(assert.Response(res))
}
```

## Examples
* Simple example with [Echo Framework](https://github.com/faabiosr/openapi-assert/blob/master/_examples/echo/main.go)


## Development

### Requirements

- Install [Go](https://golang.org)
- Install [GolangCI-Lint](https://github.com/golangci/golangci-lint#install) - Linter

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

# GolangCI-Lint
$ make lint

# Run tests
$make test
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/faabiosr/openapi-assert/blob/master/LICENSE) for more details.
