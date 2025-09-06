package main

import (
	"fmt"

	"github.com/gabrielmrtt/taski/config"
	storage_http "github.com/gabrielmrtt/taski/internal/storage/http"
	user_http "github.com/gabrielmrtt/taski/internal/user/http"
	"github.com/gin-gonic/gin"
)

func bootstrapApplication() {
	engine := gin.Default()

	apiVersion := config.Instance.ApiVersion
	appPort := config.Instance.AppPort

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	g := engine.Group(fmt.Sprintf("/api/%s", apiVersion))
	{
		user_http.BootstrapControllers(g)
		storage_http.BootstrapControllers(g)
	}

	engine.Run(fmt.Sprintf(":%s", appPort))
}

func main() {
	bootstrapApplication()
}
