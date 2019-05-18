package main

import (
	assert "github.com/faabiosr/openapi-assert"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type (
	// ResponseMsg it's a response message payload.
	ResponseMsg struct {
		Msg string `json:"message"`
	}
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// OpenAPI - Assert
	doc, err := assert.LoadFromURI("swagger.yaml")
	assert := assert.New(doc)

	if err != nil {
		e.Logger.Fatal(err)
	}

	// Routes
	e.POST("/user", func(c echo.Context) error {
		err := assert.Request(c.Request())

		if err != nil {
			return c.JSON(http.StatusBadRequest, &ResponseMsg{err.Error()})
		}

		return c.JSON(http.StatusCreated, &ResponseMsg{"user created"})
	})

	// Start server
	e.Logger.Fatal(e.Start("127.0.0.1:9000"))
}
