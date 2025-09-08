package role_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	role_database_postgres "github.com/gabrielmrtt/taski/internal/role/database/postgres"
	role_services "github.com/gabrielmrtt/taski/internal/role/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	roleRepository := role_database_postgres.NewRolePostgresRepository()
	permissionRepository := role_database_postgres.NewPermissionPostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	createRoleService := role_services.NewCreateRoleService(roleRepository, permissionRepository, transactionRepository)
	updateRoleService := role_services.NewUpdateRoleService(roleRepository, permissionRepository, transactionRepository)
	deleteRoleService := role_services.NewDeleteRoleService(roleRepository, transactionRepository)
	listRolesService := role_services.NewListRolesService(roleRepository)

	roleController := NewRoleController(createRoleService, updateRoleService, deleteRoleService, listRolesService)

	roleController.ConfigureRoutes(g)
}
