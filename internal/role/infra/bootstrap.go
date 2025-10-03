package roleinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
	rolehttp "github.com/gabrielmrtt/taski/internal/role/infra/http"
	roleservice "github.com/gabrielmrtt/taski/internal/role/service"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	roleRepository := roledatabase.NewRoleBunRepository(options.DbConnection)
	permissionRepository := roledatabase.NewPermissionBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)

	createRoleService := roleservice.NewCreateRoleService(roleRepository, permissionRepository, transactionRepository)
	updateRoleService := roleservice.NewUpdateRoleService(roleRepository, permissionRepository, transactionRepository)
	deleteRoleService := roleservice.NewDeleteRoleService(roleRepository, transactionRepository)
	listRolesService := roleservice.NewListRolesService(roleRepository)

	RoleHandler := rolehttp.NewRoleHandler(createRoleService, updateRoleService, deleteRoleService, listRolesService)

	RoleHandler.ConfigureRoutes(options.RouterGroup)
}
