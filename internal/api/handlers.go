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
	}
	return c.JSON(http.StatusOK, data)
}

func JWKHandler(c echo.Context) error {
	ctx := c.(RequestContext)
	return c.JSON(http.StatusOK, ctx.JWKManager.GetSet())
}

func NewTokenHandler(c echo.Context) error {
	type NewTokenData struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}
	ctx := c.(RequestContext)
	var tokenData NewTokenData
	if err := c.Bind(&tokenData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if tokenData.GrantType == GrantTypeClientCredentials {
		if err := ctx.ClientRepo.AuthenticateClient(tokenData.ClientId, tokenData.ClientSecret); err != nil {
			ctx.Logger().Info(err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return c.JSON(http.StatusOK, struct{}{})
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}
