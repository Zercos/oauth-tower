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
	EndpointWellKnown     = ".well-known/oauth-authorization-server"
	EndpointAuthorization = "authorization"
	EndpointToken         = "token"
	EndpointIntrospection = "introspection"
	EndpointRevocation    = "revocation"
	EndpointJWK           = "jwks.json"
)
