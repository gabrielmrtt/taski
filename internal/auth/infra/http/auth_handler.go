package authhttp

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/auth"
	authhttprequests "github.com/gabrielmrtt/taski/internal/auth/infra/http/requests"
	authservice "github.com/gabrielmrtt/taski/internal/auth/service"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserLoginService *authservice.UserLoginService
}

func NewAuthHandler(
	userLoginService *authservice.UserLoginService,
) *AuthHandler {
	return &AuthHandler{
		UserLoginService: userLoginService,
	}
}

type LoginResponse = corehttp.HttpSuccessResponseWithData[auth.UserAuthDto]

// Login godoc
// @Summary Login
// @Description Authenticates an user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body userhttprequests.UserLoginRequest true "Request body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /auth/login [post]
func (c *AuthHandler) Login(ctx *gin.Context) {
	var request authhttprequests.UserLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.UserLoginService.Execute(authservice.UserLoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *AuthHandler) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/auth")
	{
		g.POST("/login", c.Login)
	}

	return g
}
