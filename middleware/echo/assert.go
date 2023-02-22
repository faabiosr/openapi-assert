// Package echo providers a middleware for echo framework.
package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"

	assert "github.com/faabiosr/openapi-assert"
)

// AssertConfig defines the config for Assert Assert.
type AssertConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper

	// OpenAPI Document
	Document assert.Document
}

// DefaultAssertConfig is the default Assert middleware config.
var DefaultAssertConfig = AssertConfig{
	Skipper: mw.DefaultSkipper,
}

// Assert returns middleware that uses the openapi-assert
// package to assert echo HTTP requests.
func Assert(doc assert.Document) echo.MiddlewareFunc {
	c := DefaultAssertConfig
	c.Document = doc

	return AssertWithConfig(c)
}

// AssertWithConfig returns an Assert middleware with config.
func AssertWithConfig(cfg AssertConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultAssertConfig.Skipper
	}

	if cfg.Document == nil {
		panic("echo: assert middleware requires an openapi-assert document")
	}

	assert := assert.New(cfg.Document)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			if err := assert.Request(ctx.Request()); err != nil {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  err.Error(),
					Internal: err,
				}
			}

			return next(ctx)
		}
	}
}
