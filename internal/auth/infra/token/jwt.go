package authtoken

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/auth"
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenServiceOptions struct {
	Secret            string
	ExpirationMinutes int64
}

type JwtClaims struct {
	AuthenticatedUserId             string
	AuthenticatedUserOrganizationId *string
	jwt.RegisteredClaims
}

type JwtTokenService struct {
	Secret            string
	ExpirationMinutes int64
}

func NewJwtTokenService(options JwtTokenServiceOptions) *JwtTokenService {
	return &JwtTokenService{Secret: options.Secret, ExpirationMinutes: options.ExpirationMinutes}
}

func (s *JwtTokenService) GenerateToken(claims auth.TokenClaims) (string, error) {
	expirationMinutes := s.ExpirationMinutes

	now := time.Now()
	jwtClaims := &JwtClaims{
		AuthenticatedUserId:             claims.AuthenticatedUserId,
		AuthenticatedUserOrganizationId: claims.AuthenticatedUserOrganizationId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.AuthenticatedUserId,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirationMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(s.Secret))
}

func (s *JwtTokenService) ValidateToken(token string) (*auth.TokenClaims, error) {
	var tokenClaims *auth.TokenClaims
	var jwtClaims JwtClaims

	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, core.NewUnauthenticatedError("invalid token")
		}

		return []byte(s.Secret), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, core.NewUnauthenticatedError("invalid token")
	}

	tokenClaims = &auth.TokenClaims{
		AuthenticatedUserId:             jwtClaims.AuthenticatedUserId,
		AuthenticatedUserOrganizationId: jwtClaims.AuthenticatedUserOrganizationId,
	}

	return tokenClaims, nil
}
