package main

import (
	"fmt"

	"github.com/gabrielmrtt/taski/config"
	"github.com/gabrielmrtt/taski/docs"
	authinfra "github.com/gabrielmrtt/taski/internal/auth/infra"
	organizationinfra "github.com/gabrielmrtt/taski/internal/organization/infra"
	projectinfra "github.com/gabrielmrtt/taski/internal/project/infra"
	roleinfra "github.com/gabrielmrtt/taski/internal/role/infra"
	sharedpostgres "github.com/gabrielmrtt/taski/internal/shared/postgres"
	taskinfra "github.com/gabrielmrtt/taski/internal/task/infra"
	teaminfra "github.com/gabrielmrtt/taski/internal/team/infra"
	userinfra "github.com/gabrielmrtt/taski/internal/user/infra"
	workspaceinfra "github.com/gabrielmrtt/taski/internal/workspace/infra"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func bootstrapApplication() {
	engine := gin.New()

	apiVersion := config.GetInstance().ApiVersion
	appPort := config.GetInstance().AppPort

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", apiVersion)

	dbConnection := sharedpostgres.GetPostgresConnection()

	g := engine.Group(fmt.Sprintf("/api/%s", apiVersion))
	{
		authinfra.BootstrapInfra(authinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		userinfra.BootstrapInfra(userinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		organizationinfra.BootstrapInfra(organizationinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		workspaceinfra.BootstrapInfra(workspaceinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		projectinfra.BootstrapInfra(projectinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		roleinfra.BootstrapInfra(roleinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		teaminfra.BootstrapInfra(teaminfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
		taskinfra.BootstrapInfra(taskinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: dbConnection,
		})
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.StaticFile("/swagger.json", "./docs/swagger.json")
	engine.StaticFile("/docs", "./docs/redoc.html")
	engine.Run(fmt.Sprintf(":%s", appPort))
}

func main() {
	bootstrapApplication()
}
