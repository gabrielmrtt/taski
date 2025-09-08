package organization_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_database_postgres "github.com/gabrielmrtt/taski/internal/organization/database/postgres"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	role_database_postgres "github.com/gabrielmrtt/taski/internal/role/database/postgres"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	organizationRepository := organization_database_postgres.NewOrganizationPostgresRepository()
	roleRepository := role_database_postgres.NewRolePostgresRepository()
	userRepository := user_database_postgres.NewUserPostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	createOrganizationService := organization_services.NewCreateOrganizationService(organizationRepository, roleRepository, userRepository, transactionRepository)
	getOrganizationService := organization_services.NewGetOrganizationService(organizationRepository)
	listOrganizationsService := organization_services.NewListOrganizationsService(organizationRepository)
	updateOrganizationService := organization_services.NewUpdateOrganizationService(organizationRepository, transactionRepository)
	deleteOrganizationService := organization_services.NewDeleteOrganizationService(organizationRepository, transactionRepository)

	organizationController := NewOrganizationController(listOrganizationsService, getOrganizationService, createOrganizationService, updateOrganizationService, deleteOrganizationService)

	organizationController.ConfigureRoutes(g)
}
