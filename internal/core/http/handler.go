package corehttp

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type Handler interface {
	ConfigureRoutes(options ConfigureRoutesOptions)
}

type ConfigureRoutesOptions struct {
	DbConnection *bun.DB
	RouterGroup  *gin.RouterGroup
}
