package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func newServer(ctx *AppContext) *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := newRequestContext(c, ctx)
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
	ctx := NewAppContext()
	err := ctx.Init()
	if err != nil {
		log.Fatal(err)
	}
	e := newServer(ctx)

	// Routes
	e.GET("/", indexHandler)
	e.GET(EndpointWellKnown, authorizationServerWellKnownHandler)
	e.GET(EndpointJWK, JWKHandler)
	e.POST(EndpointToken, NewTokenHandler)
	return e
}

func Run(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":8000"))
}

func newRequestContext(c echo.Context, ctx *AppContext) RequestContext {
	return RequestContext{c, ctx.JWKManager, ctx.ClientRepo}
}
