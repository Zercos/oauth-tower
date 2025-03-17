package api

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

type RequestContext struct {
	echo.Context
	JWKManager *JWKManager
	ClientRepo IClientRepo
}

func (c *RequestContext) getIssuerUrl() *url.URL {
	host := config.getIssuerHost()
	if host == "" {
		host = c.Request().Host
		xForwardedHosts := c.Request().Header["X-Forwarded-Host"]
		if len(xForwardedHosts) > 0 {
			host = xForwardedHosts[0]
		}
	}
	return &url.URL{
		Scheme: c.Scheme(),
		Host:   host,
		Path:   "",
	}
}

type AppContext struct {
	initiated  bool
	JWKManager *JWKManager
	ClientRepo IClientRepo
}

func (ctx *AppContext) Init() error {
	ctx.initiated = true
	return ctx.JWKManager.LoadKeys()
}

func NewAppContext() *AppContext {
	db := initDB()
	return &AppContext{JWKManager: NewJWKManager(), ClientRepo: NewClientRepository(db)}
}
