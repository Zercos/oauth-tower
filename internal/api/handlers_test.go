package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../.env")
	db := initDB()
	defer db.Close()
	db.ClearWholeDB()
}

func TestIndexHandler(t *testing.T) {
	// given
	appCtx := NewAppContext()
	e := newServer(appCtx)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := newRequestContext(e.NewContext(req, rec), appCtx)
	expectedConfigUrl := c.getIssuerUrl().String() + EndpointWellKnown

	// when
	err := indexHandler(c)

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
	appCtx := NewAppContext()
	e := newServer(appCtx)
	req := httptest.NewRequest(http.MethodGet, EndpointWellKnown, nil)
	rec := httptest.NewRecorder()
	c := newRequestContext(e.NewContext(req, rec), appCtx)

	// when
	err := authorizationServerWellKnownHandler(c)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res WellKnownConfiguration
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, c.getIssuerUrl().String(), res.Issuer)
	assert.Equal(t, c.getIssuerUrl().String()+EndpointAuthorization, res.AuthorizationEndpoint)
}

func TestJWKHandler(t *testing.T) {
	// given
	appCtx := NewAppContext()
	e := newServer(appCtx)
	req := httptest.NewRequest(http.MethodGet, EndpointJWK, nil)
	rec := httptest.NewRecorder()
	c := newRequestContext(e.NewContext(req, rec), appCtx)

	// when
	err := JWKHandler(c)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res map[string][]map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
}

func TestClientCredentialsNewTokenHandler(t *testing.T) {
	// given
	appCtx := NewAppContext()
	err := appCtx.Init()
	assert.NoError(t, err)
	e := newServer(appCtx)
	newTokenReqJson := `{"client_id":"client1","client_secret":"secret","grant_type":"client_credentials"}`
	req := httptest.NewRequest(http.MethodPost, EndpointToken, strings.NewReader(newTokenReqJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := newRequestContext(e.NewContext(req, rec), appCtx)
	err = c.ClientRepo.AddClient(OAuthClient{"client1", "secret", "http://example.com/callback"})
	assert.NoError(t, err)

	// when
	err = NewTokenHandler(c)

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
	assert.Equal(t, c.getIssuerUrl().String(), claims["iss"])
	assert.Equal(t, "client1", claims["sub"])
	assert.Contains(t, claims, "exp")
	assert.Contains(t, claims, "iat")
}
