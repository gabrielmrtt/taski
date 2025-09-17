package organization_http_middlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_database_postgres "github.com/gabrielmrtt/taski/internal/organization/database/postgres"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

func UserMustHavePermission(permissionSlug string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permission role_core.PermissionSlugs = role_core.PermissionSlugs(permissionSlug)
		if permission == "" {
			core_http.NewHttpErrorResponse(ctx, core.NewInternalError("invalid permission slug"))
			ctx.Abort()
			return
		}

		organizationIdentity := GetOrganizationIdentityFromPath(ctx)
		if organizationIdentity.IsEmpty() {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("organizationId path parameter is required"))
			ctx.Abort()
			return
		}

		authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

		repo := organization_database_postgres.NewOrganizationUserPostgresRepository()

		orgUser, err := repo.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: organizationIdentity,
			UserIdentity:         authenticatedUserIdentity,
		})
		if err != nil {
			core_http.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if orgUser == nil || !orgUser.IsActive() {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this organization"))
			ctx.Abort()
			return
		}

		if !orgUser.CanExecuteAction(role_core.PermissionSlugs(permissionSlug)) {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you can't execute this action"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func GetOrganizationIdentityFromPath(ctx *gin.Context) core.Identity {
	organizationId := ctx.Param("organizationId")
	if organizationId == "" {
		return core.Identity{}
	}

	return core.NewIdentityFromPublic(organizationId)
}

func UserMustBeSame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organizationIdentity := GetOrganizationIdentityFromPath(ctx)
		userId := ctx.Param("userId")
		authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

		if organizationIdentity.IsEmpty() || userId == "" {
			ctx.Next()
			return
		}

		repo := organization_database_postgres.NewOrganizationUserPostgresRepository()

		orgUser, err := repo.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: organizationIdentity,
			UserIdentity:         authenticatedUserIdentity,
		})
		if err != nil {
			core_http.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if orgUser == nil {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this organization"))
			ctx.Abort()
			return
		}

		userIdentity := core.NewIdentityFromPublic(userId)

		if !userIdentity.Equals(authenticatedUserIdentity) {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you can't execute this action"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
