package projectinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	projecthttp "github.com/gabrielmrtt/taski/internal/project/infra/http"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	projectRepository := projectdatabase.NewProjectBunRepository(options.DbConnection)
	projectUserRepository := projectdatabase.NewProjectUserBunRepository(options.DbConnection)
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	workspaceRepository := workspacedatabase.NewWorkspaceBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)

	listProjectsService := projectservice.NewListProjectsService(projectRepository, workspaceRepository)
	getProjectService := projectservice.NewGetProjectService(projectRepository)
	createProjectService := projectservice.NewCreateProjectService(projectRepository, projectUserRepository, userRepository, workspaceRepository, transactionRepository)
	updateProjectService := projectservice.NewUpdateProjectService(projectRepository, transactionRepository)
	deleteProjectService := projectservice.NewDeleteProjectService(projectRepository, transactionRepository)

	configureRoutesOptions := corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	}

	projectController := projecthttp.NewProjectHandler(listProjectsService, getProjectService, createProjectService, updateProjectService, deleteProjectService)
	projectController.ConfigureRoutes(configureRoutesOptions)
}
