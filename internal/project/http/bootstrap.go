package project_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	project_database_postgres "github.com/gabrielmrtt/taski/internal/project/database/postgres"
	project_services "github.com/gabrielmrtt/taski/internal/project/services"
	workspace_database_postgres "github.com/gabrielmrtt/taski/internal/workspace/database/postgres"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	projectRepository := project_database_postgres.NewProjectPostgresRepository()
	workspaceRepository := workspace_database_postgres.NewWorkspacePostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	listProjectsService := project_services.NewListProjectsService(projectRepository, workspaceRepository)
	getProjectService := project_services.NewGetProjectService(projectRepository)
	createProjectService := project_services.NewCreateProjectService(projectRepository, workspaceRepository, transactionRepository)
	updateProjectService := project_services.NewUpdateProjectService(projectRepository, transactionRepository)
	deleteProjectService := project_services.NewDeleteProjectService(projectRepository, transactionRepository)

	projectController := NewProjectController(listProjectsService, getProjectService, createProjectService, updateProjectService, deleteProjectService)
	projectController.ConfigureRoutes(g)
}
