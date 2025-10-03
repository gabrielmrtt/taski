package main

import (
	"fmt"

	"github.com/gabrielmrtt/taski/config"
	"github.com/gabrielmrtt/taski/docs"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	organizationinfra "github.com/gabrielmrtt/taski/internal/organization/infra"
	projectinfra "github.com/gabrielmrtt/taski/internal/project/infra"
	roleinfra "github.com/gabrielmrtt/taski/internal/role/infra"
	teaminfra "github.com/gabrielmrtt/taski/internal/team/infra"
	userinfra "github.com/gabrielmrtt/taski/internal/user/infra"
	workspaceinfra "github.com/gabrielmrtt/taski/internal/workspace/infra"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func bootstrapApplication() {
	engine := gin.New()

	apiVersion := config.Instance.ApiVersion
	appPort := config.Instance.AppPort

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", apiVersion)

	g := engine.Group(fmt.Sprintf("/api/%s", apiVersion))
	{
		userinfra.BootstrapInfra(userinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
		})
		organizationinfra.BootstrapInfra(organizationinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
		})
		workspaceinfra.BootstrapInfra(workspaceinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
		})
		projectinfra.BootstrapInfra(projectinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
		})
		roleinfra.BootstrapInfra(roleinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
		})
		teaminfra.BootstrapInfra(teaminfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.GetPostgresConnection(),
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
