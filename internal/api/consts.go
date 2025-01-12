package api

// Response Types
const (
	ResponseTypeAuthorizationCodeFlow = "code"
	ResponseTypeImplicitFlowIDToken   = "id_token"
	ResponseTypeImplicitFlowToken     = "token"
	ResponseTypeImplicitFlowBoth      = "id_token token"
)

// Subject Types
const (
	SubjectTypePublic = "public"
)

// Endpoints
const (
	EndpointWellKnown     = "/.well-known/oauth-authorization-server"
	EndpointAuthorization = "/oauth/authorization"
	EndpointToken         = "/oauth/token"
	EndpointIntrospection = "/oauth/introspection"
	EndpointRevocation    = "/oauth/revocation"
	EndpointJWK           = "/oauth/jwks.json"
)
