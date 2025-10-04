package projecthttpmiddlewares

import (
	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gin-gonic/gin"
)

func UserMustBeInProject(options corehttp.MiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))
		authenticatedUserIdentity := authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		repo := projectdatabase.NewProjectUserBunRepository(options.DbConnection)

		projectUser, err := repo.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
			ProjectIdentity: projectIdentity,
			UserIdentity:    *authenticatedUserIdentity,
		})

		if err != nil {
			corehttp.NewHttpErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if projectUser == nil {
			corehttp.NewHttpErrorResponse(ctx, core.NewUnauthorizedError("you're not part of this project"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
