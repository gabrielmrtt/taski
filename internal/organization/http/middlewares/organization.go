package organization_http_middlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_database_postgres "github.com/gabrielmrtt/taski/internal/organization/database/postgres"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organizationId := ctx.Param("organization_id")

		if organizationId == "" {
			ctx.Next()
			return
		}

		organizationIdentity := core.NewIdentityFromPublic(organizationId)
		authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

		repo := organization_database_postgres.NewOrganizationPostgresRepository()

		hasUser, err := repo.CheckIfOrganizationHasUser(organizationIdentity, authenticatedUserIdentity)

		if err != nil {
			core_http.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if !hasUser {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("user is not part of the organization"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
