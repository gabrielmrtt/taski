package core_http

import "github.com/gin-gonic/gin"

type Controller interface {
	ConfigureRoutes(e *gin.Engine)
}
