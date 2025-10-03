package organizationhttpmiddlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/role"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	"github.com/gin-gonic/gin"
)

func UserMustHavePermission(permissionSlug string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permission role.PermissionSlugs = role.PermissionSlugs(permissionSlug)
		if permission == "" {
			corehttp.NewHttpErrorResponse(ctx, core.NewInternalError("invalid permission slug"))
			ctx.Abort()
			return
		}

		organizationIdentity := GetOrganizationIdentityFromPath(ctx)
		if organizationIdentity.IsEmpty() {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("organizationId path parameter is required"))
			ctx.Abort()
			return
		}

		authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		repo := organizationdatabase.NewOrganizationUserBunRepository(coredatabase.GetPostgresConnection())

		orgUser, err := repo.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: organizationIdentity,
			UserIdentity:         authenticatedUserIdentity,
		})
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if orgUser == nil || !orgUser.IsActive() {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this organization"))
			ctx.Abort()
			return
		}

		if !orgUser.CanExecuteAction(role.PermissionSlugs(permissionSlug)) {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you can't execute this action"))
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
		authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		if organizationIdentity.IsEmpty() || userId == "" {
			ctx.Next()
			return
		}

		repo := organizationdatabase.NewOrganizationUserBunRepository(coredatabase.GetPostgresConnection())

		orgUser, err := repo.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: organizationIdentity,
			UserIdentity:         authenticatedUserIdentity,
		})
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if orgUser == nil {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this organization"))
			ctx.Abort()
			return
		}

		userIdentity := core.NewIdentityFromPublic(userId)

		if !userIdentity.Equals(authenticatedUserIdentity) {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you can't execute this action"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
