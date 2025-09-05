package core_http

import "github.com/gin-gonic/gin"

type Controller interface {
	ConfigureRoutes(group *gin.RouterGroup)
}
