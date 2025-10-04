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

func UserMustHavePermission(permissionSlug string, options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permission role.PermissionSlugs = role.PermissionSlugs(permissionSlug)
		if permission == "" {
			corehttp.NewHttpErrorResponse(ctx, core.NewInternalError("invalid permission slug"))
			ctx.Abort()
			return
		}

		pathOrgnizationId := ctx.Param("organizationId")
		organizationIdentity := authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
		authenticatedUserIdentity := authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		if pathOrgnizationId != "" {
			pathOrganizationIdentity := core.NewIdentityFromPublic(pathOrgnizationId)
			organizationIdentity = &pathOrganizationIdentity
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

func UserMustBeSame(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))
		userId := ctx.Param("userId")
		authenticatedUserIdentity := authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		if organizationIdentity.IsEmpty() || userId == "" {
			ctx.Next()
			return
		}

		repo := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)

		orgUser, err := repo.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			OrganizationIdentity: organizationIdentity,
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
