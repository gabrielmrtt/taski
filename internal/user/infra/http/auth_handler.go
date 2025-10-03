package userhttp

import (
	"net/http"

	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	user "github.com/gabrielmrtt/taski/internal/user"
	userhttprequests "github.com/gabrielmrtt/taski/internal/user/infra/http/requests"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserLoginService *userservice.UserLoginService
}

func NewAuthHandler(
	userLoginService *userservice.UserLoginService,
) *AuthHandler {
	return &AuthHandler{
		UserLoginService: userLoginService,
	}
}

type LoginResponse = corehttp.HttpSuccessResponseWithData[user.UserLoginDto]

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
	var request userhttprequests.UserLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.UserLoginService.Execute(request.ToInput())
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
