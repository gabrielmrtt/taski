package main

import (
	"fmt"

	"github.com/gabrielmrtt/taski/config"
	"github.com/gabrielmrtt/taski/docs"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_http "github.com/gabrielmrtt/taski/internal/organization/http"
	project_http "github.com/gabrielmrtt/taski/internal/project/http"
	role_http "github.com/gabrielmrtt/taski/internal/role/http"
	storage_http "github.com/gabrielmrtt/taski/internal/storage/http"
	user_http "github.com/gabrielmrtt/taski/internal/user/http"
	workspace_http "github.com/gabrielmrtt/taski/internal/workspace/http"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func shutdownApplication() {
	fmt.Println("Shutting down application...")
	core_database_postgres.DB.Close()
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
		user_http.BootstrapControllers(g)
		storage_http.BootstrapControllers(g)
		organization_http.BootstrapControllers(g)
		role_http.BootstrapControllers(g)
		workspace_http.BootstrapControllers(g)
		project_http.BootstrapControllers(g)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Run(fmt.Sprintf(":%s", appPort))
}

func main() {
	bootstrapApplication()
}
