package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func newServer() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := RequestContext{c}
			return next(cc)
		}
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, path=${path}, method=${method} uri=${uri}, status=${status}, error=${error}\n",
	}))
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:L${line} ${message}")
	return e
}

func CreateServer() *echo.Echo {
	e := newServer()

	// Routes
	e.GET("/", indexHandler)
	e.GET("/"+EndpointWellKnown, authorizationServerWellKnownHandler)
	return e
}

func Run(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":8000"))
}
