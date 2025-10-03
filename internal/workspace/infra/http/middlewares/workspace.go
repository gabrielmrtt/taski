package workspacehttpmiddlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
	"github.com/gin-gonic/gin"
)

func UserMustBeInWorkspace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
		workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))

		repo := workspacedatabase.NewWorkspaceUserBunRepository(coredatabase.GetPostgresConnection())

		workspaceUser, err := repo.GetWorkspaceUserByIdentity(workspacerepo.GetWorkspaceUserByIdentityParams{
			WorkspaceIdentity: workspaceIdentity,
			UserIdentity:      userIdentity,
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
