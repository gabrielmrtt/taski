package workspaceinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	workspacehttp "github.com/gabrielmrtt/taski/internal/workspace/infra/http"
	workspaceservice "github.com/gabrielmrtt/taski/internal/workspace/service"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	workspaceRepository := workspacedatabase.NewWorkspaceBunRepository(options.DbConnection)
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	workspaceUserRepository := workspacedatabase.NewWorkspaceUserBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)

	listWorkspacesService := workspaceservice.NewListWorkspacesService(workspaceRepository)
	getWorkspaceService := workspaceservice.NewGetWorkspaceService(workspaceRepository)
	createWorkspaceService := workspaceservice.NewCreateWorkspaceService(workspaceRepository, userRepository, workspaceUserRepository, transactionRepository)
	updateWorkspaceService := workspaceservice.NewUpdateWorkspaceService(workspaceRepository, transactionRepository)
	deleteWorkspaceService := workspaceservice.NewDeleteWorkspaceService(workspaceRepository, transactionRepository)

	configureRoutesOptions := corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	}

	WorkspaceHandler := workspacehttp.NewWorkspaceHandler(listWorkspacesService, getWorkspaceService, createWorkspaceService, updateWorkspaceService, deleteWorkspaceService)
	WorkspaceHandler.ConfigureRoutes(configureRoutesOptions)
}
