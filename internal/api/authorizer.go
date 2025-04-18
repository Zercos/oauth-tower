package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func authorizerByGrantType(grantType string, ctx RequestContext) GrantTypeAuthorizer {
	switch grantType {
	case GrantTypeAuthorizationCode:
		return &AuthorizationCodeAuthorizer{ctx}
	case GrantTypeClientCredentials:
		return &ClientCredentialsAuthorizer{ctx}
	default:
		return nil
	}
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

func (a *AuthorizationCodeAuthorizer) GenerateJWT(tokenData RequestDataNewToken) (ResponseNewToken, error) {
	var newToken ResponseNewToken
	if err := a.ctx.ClientRepo.AuthenticateClient(tokenData.ClientId, tokenData.ClientSecret); err != nil {
		return newToken, err
	}

	userId, err := a.ctx.TokenRepo.GetUserIdForToken(tokenData.Code)
	if err != nil || userId == "" {
		return newToken, errors.New("invalid_grant")
	}

	token := a.makeJWT(tokenData, userId)
	signedToken, err := a.ctx.JWKManager.SignToken(token)
	if err != nil {
		return newToken, err
	}
	err = a.ctx.TokenRepo.RemoveToken(tokenData.Code)
	if err != nil {
		return newToken, err
	}
	newToken = ResponseNewToken{
		AccessToken: signedToken,
		TokenType:   "bearer",
		ExpiresIn:   config.getExpireTokenSec(),
	}
	return newToken, nil
}

func (a *AuthorizationCodeAuthorizer) makeJWT(tokenData RequestDataNewToken, userId string) *jwt.Token {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss": a.ctx.getIssuerUrl().String(),
			"sub": userId,
			"exp": time.Now().Add(time.Second * time.Duration(config.getExpireTokenSec())).Unix(),
			"iat": time.Now().Unix(),
			"aud": tokenData.ClientId,
		},
	)
	return token
}
