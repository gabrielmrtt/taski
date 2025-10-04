package authinfra

import (
	"github.com/gabrielmrtt/taski/config"
	authhttp "github.com/gabrielmrtt/taski/internal/auth/infra/http"
	authtoken "github.com/gabrielmrtt/taski/internal/auth/infra/token"
	authservice "github.com/gabrielmrtt/taski/internal/auth/service"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
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

	userLoginService := authservice.NewUserLoginService(userRepository, organizationUserRepository, tokenService)
	accessOrganizationService := authservice.NewAccessOrganizationService(organizationUserRepository, tokenService)

	handler := authhttp.NewAuthHandler(userLoginService, accessOrganizationService, tokenService)

	configureRoutesOptions := corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	}

	handler.ConfigureRoutes(configureRoutesOptions)
}
