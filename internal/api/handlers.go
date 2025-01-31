package api

import (
	"net/http"

	"github.com/golang-jwt/jwt"
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
	type NewTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
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
		token := jwt.NewWithClaims(
			jwt.SigningMethodRS256,
			jwt.MapClaims{
				"iss": ctx.getIssuerUrl(),
			},
		)
		jwk := ctx.JWKManager.GetSignKey()
		if jwk == nil {
			ctx.Logger().Info("Missing signing key")
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		signedToken, err := token.SignedString(jwk.Key)
		if err != nil {
			ctx.Logger().Info(err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		resp := NewTokenResponse{
			AccessToken: signedToken,
			TokenType:   "bearer",
			ExpiresIn:   120,
		}
		return c.JSON(http.StatusOK, resp)
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}
