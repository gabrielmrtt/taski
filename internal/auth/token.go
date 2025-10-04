package auth

type TokenClaims struct {
	AuthenticatedUserId             string
	AuthenticatedUserOrganizationId *string
}

type TokenService interface {
	GenerateToken(claims TokenClaims) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}
