package teaminfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	teamdatabase "github.com/gabrielmrtt/taski/internal/team/infra/database"
	teamhttp "github.com/gabrielmrtt/taski/internal/team/infra/http"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	teamRepository := teamdatabase.NewTeamBunRepository(options.DbConnection)
	organizationUserRepository := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)

	listTeamsService := teamservice.NewListTeamsService(teamRepository)
	getTeamService := teamservice.NewGetTeamService(teamRepository)
	createTeamService := teamservice.NewCreateTeamService(teamRepository, organizationUserRepository, transactionRepository)
	updateTeamService := teamservice.NewUpdateTeamService(teamRepository, organizationUserRepository, transactionRepository)
	deleteTeamService := teamservice.NewDeleteTeamService(teamRepository, transactionRepository)

	TeamHandler := teamhttp.NewTeamHandler(listTeamsService, getTeamService, createTeamService, updateTeamService, deleteTeamService)

	TeamHandler.ConfigureRoutes(options.RouterGroup)
}
