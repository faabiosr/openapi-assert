package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	assert "github.com/faabiosr/openapi-assert"
	mw "github.com/faabiosr/openapi-assert/middleware/echo"
)

type (
	// ResponseMsg it's a response message payload.
	ResponseMsg struct {
		Msg string `json:"message"`
	}

	// User payload.
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
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

	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(mw.Assert(doc))

	// Routes
	e.POST("/user", func(c echo.Context) error {
		//		err := assert.Request(c.Request())
		//
		//		if err != nil {
		//			return c.JSON(http.StatusBadRequest, &ResponseMsg{err.Error()})
		//		}

		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, &ResponseMsg{err.Error()})
		}

		msg := fmt.Sprintf("user %s with email %s was created", user.Name, user.Email)

		return c.JSON(http.StatusCreated, &ResponseMsg{msg})
	})

	// Start server
	e.Logger.Fatal(e.Start("127.0.0.1:9000"))
}
