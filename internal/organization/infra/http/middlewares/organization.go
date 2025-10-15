package organizationhttpmiddlewares

import (
	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/role"
	"github.com/gin-gonic/gin"
)

// UserMustHavePermission is a middleware that checks if the user is part of the organization and has the permission to execute an action
func UserMustHavePermission(permissionSlug string, options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permission role.PermissionSlugs = role.PermissionSlugs(permissionSlug)
		if permission == "" {
			corehttp.NewHttpErrorResponse(ctx, core.NewInternalError("invalid permission slug"))
			ctx.Abort()
			return
		}

		var organizationIdentity *core.Identity = nil
		var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		pathOrganizationId := ctx.Param("organizationId")
		if pathOrganizationId != "" {
			identity := core.NewIdentityFromPublic(pathOrganizationId)
			organizationIdentity = &identity
		} else {
			organizationIdentity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
		}

		if organizationIdentity == nil {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you need to access an organization to execute this action"))
			ctx.Abort()
			return
		}

		repo := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)

		orgUser, err := repo.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: *organizationIdentity,
			UserIdentity:         *authenticatedUserIdentity,
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

// UserMustBeSame is a middleware that checks if the authenticated user is the same as the organization user from the path parameter
func UserMustBeSame(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var organizationIdentity *core.Identity = nil
		var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		pathOrganizationId := ctx.Param("organizationId")
		if pathOrganizationId != "" {
			identity := core.NewIdentityFromPublic(pathOrganizationId)
			organizationIdentity = &identity
		} else {
			organizationIdentity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
		}

		userId := ctx.Param("userId")

		repo := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)

		orgUser, err := repo.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: *organizationIdentity,
			UserIdentity:         *authenticatedUserIdentity,
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

		if !userIdentity.Equals(*authenticatedUserIdentity) {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you can't execute this action"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
