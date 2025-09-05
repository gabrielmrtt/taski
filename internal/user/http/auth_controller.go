package user_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_http_requests "github.com/gabrielmrtt/taski/internal/user/http/requests"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserLoginService *user_services.UserLoginService
}

func NewAuthController(
	userLoginService *user_services.UserLoginService,
) *AuthController {
	return &AuthController{
		UserLoginService: userLoginService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request user_http_requests.UserLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.UserLoginService.Execute(request.ToInput())
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *AuthController) ConfigureRoutes(engine *gin.Engine) {
	engine.Group("/auth")
	{
		engine.POST("/login", c.Login)
	}
}
