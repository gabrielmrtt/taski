package workspace_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	workspace_database_postgres "github.com/gabrielmrtt/taski/internal/workspace/database/postgres"
	workspace_services "github.com/gabrielmrtt/taski/internal/workspace/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	workspaceRepository := workspace_database_postgres.NewWorkspacePostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	listWorkspacesService := workspace_services.NewListWorkspacesService(workspaceRepository)
	getWorkspaceService := workspace_services.NewGetWorkspaceService(workspaceRepository)
	createWorkspaceService := workspace_services.NewCreateWorkspaceService(workspaceRepository, transactionRepository)
	updateWorkspaceService := workspace_services.NewUpdateWorkspaceService(workspaceRepository, transactionRepository)
	deleteWorkspaceService := workspace_services.NewDeleteWorkspaceService(workspaceRepository, transactionRepository)

	workspaceController := NewWorkspaceController(listWorkspacesService, getWorkspaceService, createWorkspaceService, updateWorkspaceService, deleteWorkspaceService)
	workspaceController.ConfigureRoutes(g)
}
