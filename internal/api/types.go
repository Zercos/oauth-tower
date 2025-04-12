package api

type WellKnownConfiguration struct {
	Issuer                                     string   `json:"issuer"`
	JWKsUri                                    string   `json:"jwks_uri,omitempty"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint,omitempty"`
	SubjectTypesSupported                      []string `json:"subject_types_supported"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
	ScopesSupported                            []string `json:"scopes_supported,omitempty"`
	ClaimsSupported                            []string `json:"claims_supported,omitempty"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	IntrospectionEndpoint                      string   `json:"introspection_endpoint,omitempty"`
	RevocationEndpoint                         string   `json:"revocation_endpoint,omitempty"`
	RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
}

type RequestDataNewToken struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type RequestQueryParamAuthorize struct {
	ClientId     string `query:"client_id"`
	ResponseType string `query:"response_type"`
	RedirectURI  string `query:"redirect_uri"`
	State        string `query:"state"`
}

type RequestDataNewLogin struct {
	Redirect string `form:"redirect"`
	Username string `form:"username"`
	Password string `form:"password"`
}

type ResponseNewToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint16 `json:"expires_in"`
}

type GrantTypeAuthorizer interface {
	GenerateJWT(tokenData RequestDataNewToken) (ResponseNewToken, error)
}

type IClientRepo interface {
	GetClient(clientId string) (OAuthClient, error)
	AddClient(client OAuthClient, checkExists bool) error
	AuthenticateClient(clientId string, clientSecret string) error
}

type IUserRepo interface {
	GetUser(username string) (UserModel, error)
	AuthenticateUser(username string, password string) error
	AddUser(user NewUser, checkExists bool) error
}
