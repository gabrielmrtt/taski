package projecthttpmiddlewares

import (
	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	"github.com/gin-gonic/gin"
)

func UserMustBeInProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))
		authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

		repo := projectdatabase.NewProjectUserBunRepository(coredatabase.GetPostgresConnection())

		projectUser, err := repo.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
			ProjectIdentity: projectIdentity,
			UserIdentity:    authenticatedUserIdentity,
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
