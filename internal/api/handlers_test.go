package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../.env")
	db := initDB()
	defer db.Close()
	db.ClearWholeDB()
}

func mkTestRequestCtx(t *testing.T, req *http.Request, rec *httptest.ResponseRecorder) RequestContext {
	appCtx := NewAppContext()
	assert.NoError(t, appCtx.Init())
	e := newServer(appCtx)
	return newRequestContext(e.NewContext(req, rec), appCtx)
}

func TestIndexHandler(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)
	expectedConfigUrl := ctx.getIssuerUrl().String() + EndpointWellKnown

	// when
	err := indexHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "running", res["status"])
	assert.Equal(t, expectedConfigUrl, res["config"])
}

func TestAuthorizationServerWellKnownHandler(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, EndpointWellKnown, nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)

	// when
	err := authorizationServerWellKnownHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res WellKnownConfiguration
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, ctx.getIssuerUrl().String(), res.Issuer)
	assert.Equal(t, ctx.getIssuerUrl().String()+EndpointAuthorization, res.AuthorizationEndpoint)
}

func TestJWKHandler(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, EndpointJWK, nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)

	// when
	err := JWKHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res map[string][]map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
}

func TestClientCredentialsNewTokenHandler(t *testing.T) {
	// given
	newTokenReqJson := `{"client_id":"client1","client_secret":"secret","grant_type":"client_credentials"}`
	req := httptest.NewRequest(http.MethodPost, EndpointToken, strings.NewReader(newTokenReqJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)
	err := ctx.ClientRepo.AddClient(OAuthClient{"client1", "secret", "http://example.com/callback"}, true)
	assert.NoError(t, err)

	// when
	err = NewTokenHandler(ctx)

	// then
	assert.NoError(t, err)

	var res map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	tokenType := res["token_type"]
	assert.Equal(t, "bearer", tokenType)
	tokenString := res["access_token"].(string)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	assert.NoError(t, err)
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, ctx.getIssuerUrl().String(), claims["iss"])
	assert.Equal(t, "client1", claims["sub"])
	assert.Contains(t, claims, "exp")
	assert.Contains(t, claims, "iat")
}

func TestAuthorizeRedirectsToLoginWhenNotAuthenticated(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, EndpointAuthorization+"?response_type=code&client_id=client1&redirect_uri=http://example.com/callback", nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)

	err := ctx.ClientRepo.AddClient(OAuthClient{"client1", "secret", "http://example.com/callback"}, true)
	assert.NoError(t, err)

	// when
	err = AuthorizationHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Location"), "/login?redirect=")
}

func TestAuthorizeFailsWithInvalidClient(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, EndpointAuthorization+"?response_type=code&client_id=wrong&redirect_uri=http://localhost:8081/callback", nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)

	// when
	err := AuthorizationHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid_client")
}

func TestAuthorizeReturnsAuthCodeWhenAuthenticated(t *testing.T) {
	// given

	req := httptest.NewRequest(http.MethodGet, EndpointAuthorization+"?response_type=code&client_id=client1&redirect_uri=http://example.com/callback", nil)
	rec := httptest.NewRecorder()
	ctx := mkTestRequestCtx(t, req, rec)
	ctx.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	err := ctx.ClientRepo.AddClient(OAuthClient{"client1", "secret", "http://example.com/callback"}, true)
	assert.NoError(t, err)

	// Manually create session with user_id
	sess, err := session.Get("session", ctx)
	assert.NoError(t, err)
	sess.Values["user_id"] = "testuser"
	sess.Save(req, rec)

	// when
	err = AuthorizationHandler(ctx)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	location := rec.Header().Get("Location")
	assert.Contains(t, location, "code=")
}
