package server

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type WellKnownConfiguration struct {
	Issuer                                     string   `json:"issuer"`
	JWKSURI                                    string   `json:"jwks_uri,omitempty"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint,omitempty"`
	SubjectTypesSupported                      []string `json:"subject_types_supported"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
	ScopesSupported                            []string `json:"scopes_supported,omitempty"`
	ClaimsSupported                            []string `json:"claims_supported,omitempty"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	IntrospectionEndpoint                      string   `json:"introspection_endpoint,omitempty"`
	RevocationEndpoint                         string   `json:"revocation_endpoint,omitempty"`
	RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
}

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

func CreateServer() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &RequestContext{c}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	// Routes
	e.GET("/", index)
	e.GET("/.well-known/oauth-authorization-server", authorizationServerWellKnown)
	return e
}

func Run(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":8000"))
}

func index(c echo.Context) error {
	cc := c.(*RequestContext)
	response := map[string]string{
		"message": "Welcome to the OAuth-Tower - OAuth 2.0 Authorization Server",
		"status":  "running",
		"config":  cc.getIssuerUrl().String() + "/.well-known/oauth-authorization-server",
	}
	return c.JSON(http.StatusOK, response)
}

func authorizationServerWellKnown(c echo.Context) error {
	cc := c.(*RequestContext)
	issuer := cc.getIssuerUrl().String()
	data := WellKnownConfiguration{
		Issuer:                issuer,
		AuthorizationEndpoint: issuer + "/authorization",
		SubjectTypesSupported: []string{SubjectTypePublic},
		ResponseTypesSupported: []string{
			ResponseTypeAuthorizationCodeFlow,
			ResponseTypeImplicitFlowBoth,
			ResponseTypeImplicitFlowIDToken,
			ResponseTypeImplicitFlowToken,
		},
	}
	return c.JSON(http.StatusOK, data)
}
