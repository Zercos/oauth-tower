package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func indexHandler(c echo.Context) error {
	cc := c.(RequestContext)
	response := map[string]string{
		"message": "Welcome to the OAuth-Tower - OAuth 2.0 Authorization Server",
		"status":  "running",
		"config":  cc.getIssuerUrl().String() + EndpointWellKnown,
	}
	return c.JSON(http.StatusOK, response)
}

func authorizationServerWellKnownHandler(c echo.Context) error {
	cc := c.(RequestContext)
	issuer := cc.getIssuerUrl().String()
	data := WellKnownConfiguration{
		Issuer:                issuer,
		AuthorizationEndpoint: issuer + EndpointAuthorization,
		TokenEndpoint:         issuer + EndpointToken,
		IntrospectionEndpoint: issuer + EndpointIntrospection,
		RevocationEndpoint:    issuer + EndpointRevocation,
		JWKsUri:               issuer + EndpointJWK,
		SubjectTypesSupported: []string{SubjectTypePublic},
		ResponseTypesSupported: []string{
			ResponseTypeAuthorizationCodeFlow,
			ResponseTypeImplicitFlowBoth,
			ResponseTypeImplicitFlowIDToken,
			ResponseTypeImplicitFlowToken,
		},
		ClaimsSupported: []string{
			"iss", "sub", "aud", "exp", "iat", "nbf",
		},
	}
	return c.JSON(http.StatusOK, data)
}

func JWKHandler(c echo.Context) error {
	ctx := c.(RequestContext)
	return c.JSON(http.StatusOK, ctx.JWKManager.GetSet())
}

func NewTokenHandler(c echo.Context) error {
	ctx := c.(RequestContext)
	var tokenData RequestDataNewToken
	if err := c.Bind(&tokenData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	authorizer := authorizerByGrantType(tokenData.GrantType, ctx)
	if authorizer == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	newToken, err := authorizer.GenerateJWT(tokenData)
	if err != nil {
		ctx.Logger().Info(err)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, newToken)
}
