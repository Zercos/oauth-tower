package api

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func authorizerByGrantType(grantType string, ctx RequestContext) GrantTypeAuthorizer {
	if grantType == GrantTypeClientCredentials {
		return &ClientCredentialsAuthorizer{ctx}
	}

	return nil
}

type ClientCredentialsAuthorizer struct {
	ctx RequestContext
}

func (a *ClientCredentialsAuthorizer) GenerateJWT(tokenData RequestDataNewToken) (ResponseNewToken, error) {
	if err := a.ctx.ClientRepo.AuthenticateClient(tokenData.ClientId, tokenData.ClientSecret); err != nil {
		return ResponseNewToken{}, err
	}
	token := a.makeJWT(tokenData.ClientId)
	signedToken, err := a.ctx.JWKManager.SignToken(token)
	if err != nil {
		return ResponseNewToken{}, err
	}
	newToken := ResponseNewToken{
		AccessToken: signedToken,
		TokenType:   "bearer",
		ExpiresIn:   config.getExpireTokenSec(),
	}
	return newToken, nil
}

func (a *ClientCredentialsAuthorizer) makeJWT(clientId string) *jwt.Token {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss": a.ctx.getIssuerUrl().String(),
			"sub": clientId,
			"exp": time.Now().Add(time.Second * time.Duration(config.getExpireTokenSec())).Unix(),
			"iat": time.Now().Unix(),
		},
	)
	return token
}
