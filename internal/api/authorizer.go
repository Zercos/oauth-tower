package api

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
	jwk := a.ctx.JWKManager.GetSignKey()
	issuer := a.ctx.getIssuerUrl().String()
	return generateNewJWT(issuer, jwk)
}
