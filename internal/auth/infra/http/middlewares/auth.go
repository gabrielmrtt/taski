package authhttpmiddlewares

import (
	"strings"

	"github.com/gabrielmrtt/taski/config"
	"github.com/gabrielmrtt/taski/internal/auth"
	authtoken "github.com/gabrielmrtt/taski/internal/auth/infra/token"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gin-gonic/gin"
)

func GetAuthenticatedUserIdentity(ctx *gin.Context) *core.Identity {
	authenticatedUserId := ctx.GetString("authenticated_user_id")

	if authenticatedUserId == "" {
		return nil
	}

	identity := core.NewIdentityFromPublic(authenticatedUserId)
	return &identity
}

func GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx *gin.Context) *core.Identity {
	authenticatedUserOrganizationId := ctx.GetString("authenticated_user_organization_id")

	if authenticatedUserOrganizationId == "" {
		return nil
	}

	identity := core.NewIdentityFromPublic(authenticatedUserOrganizationId)
	return &identity
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

func AuthMiddleware(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtSecret := config.GetInstance().JwtSecret

		tokenService := authtoken.NewJwtTokenService(authtoken.JwtTokenServiceOptions{
			Secret:            jwtSecret,
			ExpirationMinutes: config.GetInstance().JwtExpirationMinutes,
		})

		tokenStr, err := extractBearerToken(ctx)
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		var claims *auth.TokenClaims

		claims, err = tokenService.ValidateToken(tokenStr)
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set("authenticated_user_id", claims.AuthenticatedUserId)
		ctx.Set("authenticated_user_organization_id", *claims.AuthenticatedUserOrganizationId)
		ctx.Next()
	}
}
