package workspace_http_middlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	workspace_database_postgres "github.com/gabrielmrtt/taski/internal/workspace/database/postgres"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
	"github.com/gin-gonic/gin"
)

func UserMustBeInWorkspace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
		workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspace_id"))

		repo := workspace_database_postgres.NewWorkspaceUserPostgresRepository()

		workspaceUser, err := repo.GetWorkspaceUserByIdentity(workspace_repositories.GetWorkspaceUserByIdentityParams{
			WorkspaceIdentity: workspaceIdentity,
			UserIdentity:      userIdentity,
		})
		if err != nil {
			core_http.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if workspaceUser == nil {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this workspace"))
			ctx.Abort()
			return
		}

		ctx.Next()
		return
	}
}
