package userhttpmiddlewares

import (
	"strings"
	"time"

	"github.com/gabrielmrtt/taski/config"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	AuthenticatedUserId string `json:"sub"`
	jwt.RegisteredClaims
}

func GetAuthenticatedUserIdentity(ctx *gin.Context) core.Identity {
	authenticatedUserId := ctx.GetString("authenticated_user_id")

	if authenticatedUserId == "" {
		return core.Identity{}
	}

	return core.NewIdentityFromPublic(authenticatedUserId)
}

func GenerateJwtToken(userIdentity core.Identity) (string, error) {
	expirationMinutes := config.Instance.JwtExpirationMinutes

	now := time.Now()
	claims := &JwtClaims{
		AuthenticatedUserId: userIdentity.Public,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userIdentity.Public,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirationMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Instance.JwtSecret))
}

func extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", core.NewUnauthenticatedError("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", core.NewUnauthenticatedError("invalid authorization header")
	}

	return parts[1], nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtSecret := config.Instance.JwtSecret

		tokenStr, err := extractBearerToken(ctx)
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		var claims JwtClaims

		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, core.NewUnauthenticatedError("invalid token")
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set("authenticated_user_id", claims.AuthenticatedUserId)
		ctx.Next()
	}
}
