package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {
	// given
	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := RequestContext{e.NewContext(req, rec)}
	expectedConfigUrl := c.getIssuerUrl().String() + "/" + EndpointWellKnown

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
	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/"+EndpointWellKnown, nil)
	rec := httptest.NewRecorder()
	c := RequestContext{e.NewContext(req, rec)}

	// when
	err := authorizationServerWellKnownHandler(c)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res WellKnownConfiguration
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, c.getIssuerUrl().String(), res.Issuer)
	assert.Equal(t, c.getIssuerUrl().String()+"/"+EndpointAuthorization, res.AuthorizationEndpoint)
}
