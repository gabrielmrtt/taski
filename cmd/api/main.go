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

func shutdownApplication() {
	fmt.Println("Shutting down application...")
	coredatabase.DB.Close()
}

func bootstrapApplication() {
	defer shutdownApplication()

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
			DbConnection: coredatabase.DB,
		})
		organizationinfra.BootstrapInfra(organizationinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.DB,
		})
		workspaceinfra.BootstrapInfra(workspaceinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.DB,
		})
		projectinfra.BootstrapInfra(projectinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.DB,
		})
		roleinfra.BootstrapInfra(roleinfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.DB,
		})
		teaminfra.BootstrapInfra(teaminfra.BootstrapInfraOptions{
			RouterGroup:  g,
			DbConnection: coredatabase.DB,
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
