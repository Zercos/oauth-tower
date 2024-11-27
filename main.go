package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	// Routes
	e.GET("/", index)

	e.Logger.Fatal(e.Start(":8000"))
}

func index(c echo.Context) error {
	response := map[string]string{
		"message": "Welcome to the OAuth-Tower - OAuth 2.0 Authorization Server",
		"status":  "running",
		"config":  "https://example.com/.well-known/oauth-authorization-server",
	}
	return c.JSON(http.StatusOK, response)
}
