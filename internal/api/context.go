package api

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

type RequestContext struct {
	echo.Context
}

func (c *RequestContext) getIssuerUrl() *url.URL {
	host := c.Request().Host
	xForwardedHosts := c.Request().Header["X-Forwarded-Host"]
	if len(xForwardedHosts) > 0 {
		host = xForwardedHosts[0]
	}
	return &url.URL{
		Scheme: c.Scheme(),
		Host:   host,
		Path:   "",
	}
}
