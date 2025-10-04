package authinfra

import (
	"github.com/gabrielmrtt/taski/config"
	authtoken "github.com/gabrielmrtt/taski/internal/auth/infra/token"
	authservice "github.com/gabrielmrtt/taski/internal/auth/service"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
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
	organizationUserRepository := organizationdatabase.NewOrganizationUserBunRepository(options.DbConnection)

	tokenService := authtoken.NewJwtTokenService(authtoken.JwtTokenServiceOptions{
		Secret:            config.GetInstance().JwtSecret,
		ExpirationMinutes: config.GetInstance().JwtExpirationMinutes,
	})

	authService := authservice.NewUserLoginService(userRepository, organizationUserRepository, tokenService)

	// Return the service instead of configuring routes directly
	_ = authService
}
