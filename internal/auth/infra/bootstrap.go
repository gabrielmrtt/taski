package authinfra

import (
	authhttp "github.com/gabrielmrtt/taski/internal/auth/infra/http"
	authservice "github.com/gabrielmrtt/taski/internal/auth/service"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	authService := authservice.NewUserLoginService(userRepository)
	authController := authhttp.NewAuthHandler(authService)
	authController.ConfigureRoutes(options.RouterGroup)
}
