package project_http_middlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	project_database_postgres "github.com/gabrielmrtt/taski/internal/project/database/postgres"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

func UserMustBeInProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))
		authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

		repo := project_database_postgres.NewProjectUserPostgresRepository()

		projectUser, err := repo.GetProjectUserByIdentity(project_repositories.GetProjectUserByIdentityParams{
			ProjectIdentity: projectIdentity,
			UserIdentity:    authenticatedUserIdentity,
		})

		if err != nil {
			core_http.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if projectUser == nil {
			core_http.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this project"))
			ctx.Abort()
			return
		}

		ctx.Next()
		return
	}
}
