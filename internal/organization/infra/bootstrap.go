package organizationinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	organizationhttp "github.com/gabrielmrtt/taski/internal/organization/infra/http"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
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
	organizationRepository := organizationdatabase.NewOrganizationBunRepository(options.DbConnection)
	organizationUserRepository := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)
	roleRepository := roledatabase.NewRoleBunRepository(options.DbConnection)
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	workspaceRepository := workspacedatabase.NewWorkspaceBunRepository(options.DbConnection)
	workspaceUserRepository := workspacedatabase.NewWorkspaceUserBunRepository(options.DbConnection)
	projectRepository := projectdatabase.NewProjectBunRepository(options.DbConnection)
	projectUserRepository := projectdatabase.NewProjectUserBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)

	createOrganizationService := organizationservice.NewCreateOrganizationService(organizationRepository, organizationUserRepository, roleRepository, userRepository, transactionRepository)
	getOrganizationService := organizationservice.NewGetOrganizationService(organizationRepository)
	listOrganizationsService := organizationservice.NewListOrganizationsService(organizationRepository)
	updateOrganizationService := organizationservice.NewUpdateOrganizationService(organizationRepository, transactionRepository)
	deleteOrganizationService := organizationservice.NewDeleteOrganizationService(organizationRepository, transactionRepository)
	inviteUserToOrganizationService := organizationservice.NewInviteUserToOrganizationService(organizationRepository, organizationUserRepository, userRepository, roleRepository, workspaceRepository, workspaceUserRepository, projectRepository, projectUserRepository, transactionRepository)
	removeUserFromOrganizationService := organizationservice.NewRemoveUserFromOrganizationService(organizationRepository, organizationUserRepository, userRepository, transactionRepository)
	acceptOrganizationUserInvitationService := organizationservice.NewAcceptOrganizationUserInvitationService(organizationUserRepository, workspaceUserRepository, projectUserRepository, transactionRepository)
	refuseOrganizationUserInvitationService := organizationservice.NewRefuseOrganizationUserInvitationService(organizationUserRepository, workspaceUserRepository, projectUserRepository, transactionRepository)
	listOrganizationUsersService := organizationservice.NewListOrganizationUsersService(organizationUserRepository)
	getOrganizationUserService := organizationservice.NewGetOrganizationUserService(organizationUserRepository)
	updateOrganizationUserService := organizationservice.NewUpdateOrganizationUserService(organizationRepository, organizationUserRepository, roleRepository, workspaceRepository, projectRepository, workspaceUserRepository, projectUserRepository, transactionRepository)
	listMyOrganizationInvitesService := organizationservice.NewListMyOrganizationInvitesService(organizationRepository)

	OrganizationHandler := organizationhttp.NewOrganizationHandler(listOrganizationsService, getOrganizationService, createOrganizationService, updateOrganizationService, deleteOrganizationService)
	OrganizationUserHandler := organizationhttp.NewOrganizationUserHandler(listOrganizationUsersService, inviteUserToOrganizationService, removeUserFromOrganizationService, getOrganizationUserService, updateOrganizationUserService)
	OrganizationInvitesHandler := organizationhttp.NewOrganizationInvitesHandler(listMyOrganizationInvitesService, acceptOrganizationUserInvitationService, refuseOrganizationUserInvitationService)

	configureRoutesOptions := corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	}

	OrganizationHandler.ConfigureRoutes(configureRoutesOptions)
	OrganizationUserHandler.ConfigureRoutes(configureRoutesOptions)
	OrganizationInvitesHandler.ConfigureRoutes(configureRoutesOptions)
}
