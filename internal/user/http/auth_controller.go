package user_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_core "github.com/gabrielmrtt/taski/internal/user"
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

type LoginResponse = core_http.HttpSuccessResponseWithData[user_core.UserLoginDto]

// Login godoc
// @Summary Login
// @Description Authenticates an user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body user_http_requests.UserLoginRequest true "Request body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /auth/login [post]
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

func (c *AuthController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/auth")
	{
		g.POST("/login", c.Login)
	}

	return g
}
