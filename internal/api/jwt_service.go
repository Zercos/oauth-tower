package api

import (
	"fmt"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(issuer string, jwk *jose.JSONWebKey) (ResponseNewToken, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss": issuer,
		},
	)
	if jwk == nil {
		return ResponseNewToken{}, fmt.Errorf("Missing signing key")
	}
	signedToken, err := token.SignedString(jwk.Key)
	if err != nil {
		return ResponseNewToken{}, err
	}
	newToken := ResponseNewToken{
		AccessToken: signedToken,
		TokenType:   "bearer",
		ExpiresIn:   120,
	}
	return newToken, nil
}
