package team_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_database_postgres "github.com/gabrielmrtt/taski/internal/organization/database/postgres"
	team_database_postgres "github.com/gabrielmrtt/taski/internal/team/database/postgres"
	team_services "github.com/gabrielmrtt/taski/internal/team/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	teamRepository := team_database_postgres.NewTeamPostgresRepository()
	organizationUserRepository := organization_database_postgres.NewOrganizationUserPostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	listTeamsService := team_services.NewListTeamsService(teamRepository)
	getTeamService := team_services.NewGetTeamService(teamRepository)
	createTeamService := team_services.NewCreateTeamService(teamRepository, organizationUserRepository, transactionRepository)
	updateTeamService := team_services.NewUpdateTeamService(teamRepository, organizationUserRepository, transactionRepository)
	deleteTeamService := team_services.NewDeleteTeamService(teamRepository, transactionRepository)

	teamController := NewTeamController(listTeamsService, getTeamService, createTeamService, updateTeamService, deleteTeamService)

	teamController.ConfigureRoutes(g)
}
