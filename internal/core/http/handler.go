package corehttp

import "github.com/gin-gonic/gin"

type Handler interface {
	ConfigureRoutes(group *gin.RouterGroup)
}
