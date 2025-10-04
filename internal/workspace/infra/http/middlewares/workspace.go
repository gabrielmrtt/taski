package workspacehttpmiddlewares

import (
	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
	"github.com/gin-gonic/gin"
)

func UserMustBeInWorkspace(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdentity := authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
		workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))

		repo := workspacedatabase.NewWorkspaceUserBunRepository(options.DbConnection)

		workspaceUser, err := repo.GetWorkspaceUserByIdentity(workspacerepo.GetWorkspaceUserByIdentityParams{
			WorkspaceIdentity: workspaceIdentity,
			UserIdentity:      *userIdentity,
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
