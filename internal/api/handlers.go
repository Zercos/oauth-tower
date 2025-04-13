package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var authCodes = map[string]string{}

func IndexHandler(c echo.Context) error {
	cc := c.(RequestContext)
	response := map[string]string{
		"message": "Welcome to the OAuth-Tower - OAuth 2.0 Authorization Server",
		"status":  "running",
		"config":  cc.getIssuerUrl().String() + EndpointWellKnown,
	}
	return c.JSON(http.StatusOK, response)
}

func AuthorizationServerWellKnownHandler(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported_grant_type")
	}
	newToken, err := authorizer.GenerateJWT(tokenData)
	if err != nil {
		ctx.Logger().Info(err)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, newToken)
}

func AuthorizationHandler(c echo.Context) error {
	ctx := c.(RequestContext)
	var reqData RequestQueryParamAuthorize
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if reqData.ResponseType != ResponseTypeAuthorizationCodeFlow {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "unsupported_response_type"})
	}
	client, err := ctx.ClientRepo.GetClient(reqData.ClientId)
	if err != nil || client.RedirectURI != reqData.RedirectURI {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid_client"})
	}

	sess, err := session.Get("session", c)
	if err != nil || sess.Values["user_id"] == nil {
		// Not authenticated -> redirect to login with original URL
		originalQuery := c.Request().URL.RawQuery
		loginRedirect := fmt.Sprintf("%s?redirect=/authorize?%s", EndpointAuthorizationLogin, url.QueryEscape(originalQuery))
		return c.Redirect(http.StatusFound, loginRedirect)
	}

	userID := sess.Values["user_id"].(string)
	code := uuid.New().String()
	authCodes[code] = userID

	redirectWithParams := reqData.RedirectURI + "?code=" + code
	if reqData.State != "" {
		redirectWithParams += "&state=" + reqData.State
	}
	return c.Redirect(http.StatusFound, redirectWithParams)
}

func LoginPageHandler(c echo.Context) error {
	redirect := c.QueryParam("redirect")
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Redirect": redirect,
	})
}

func UserLoginHandler(c echo.Context) error {
	ctx := c.(RequestContext)
	var loginData RequestDataNewLogin
	if err := c.Bind(&loginData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := ctx.UserRepo.AuthenticateUser(loginData.Username, loginData.Password); err != nil {
		ctx.Logger().Error(err)
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}
	user, err := ctx.UserRepo.GetUser(loginData.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	sess, _ := session.Get("session", c)
	sess.Values["user_id"] = user.UserId
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, loginData.Redirect)

}
