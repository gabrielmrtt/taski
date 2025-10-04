package workspacehttpmiddlewares

import (
	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
	"github.com/gin-gonic/gin"
)

// UserMustBeInWorkspace is a middleware that checks if the authenticated user is part of the workspace
func UserMustBeInWorkspace(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
		var workspaceIdentity *core.Identity = nil

		pathWorkspaceId := ctx.Param("workspaceId")
		if pathWorkspaceId != "" {
			identity := core.NewIdentityFromPublic(pathWorkspaceId)
			workspaceIdentity = &identity
		}

		repo := workspacedatabase.NewWorkspaceUserBunRepository(options.DbConnection)

		workspaceUser, err := repo.GetWorkspaceUserByIdentity(workspacerepo.GetWorkspaceUserByIdentityParams{
			WorkspaceIdentity: *workspaceIdentity,
			UserIdentity:      *authenticatedUserIdentity,
		})
		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if workspaceUser == nil {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this workspace"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
