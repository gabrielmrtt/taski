package main

import (
	"fmt"

	"github.com/gabrielmrtt/taski/config"
	user_http "github.com/gabrielmrtt/taski/internal/user/http"
	"github.com/gin-gonic/gin"
)

func bootstrapApplication() {
	engine := gin.Default()

	apiVersion := config.Instance.ApiVersion
	appPort := config.Instance.AppPort

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	engine.Group(fmt.Sprintf("/api/%s", apiVersion))
	{
		user_http.BootstrapControllers(engine)
	}

	engine.Run(fmt.Sprintf(":%s", appPort))
}

func main() {
	bootstrapApplication()
}
