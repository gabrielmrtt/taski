package user_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_http_requests "github.com/gabrielmrtt/taski/internal/user/http/requests"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
	"github.com/gin-gonic/gin"
)

type UserRegistrationController struct {
	RegisterUserService           *user_services.RegisterUserService
	VerifyUserRegistrationService *user_services.VerifyUserRegistrationService
	ForgotUserPasswordService     *user_services.ForgotUserPasswordService
	RecoverUserPasswordService    *user_services.RecoverUserPasswordService
}

func NewUserRegistrationController(
	registerUserService *user_services.RegisterUserService,
	verifyUserRegistrationService *user_services.VerifyUserRegistrationService,
	forgotUserPasswordService *user_services.ForgotUserPasswordService,
	recoverUserPasswordService *user_services.RecoverUserPasswordService,
) *UserRegistrationController {
	return &UserRegistrationController{
		RegisterUserService:           registerUserService,
		VerifyUserRegistrationService: verifyUserRegistrationService,
		ForgotUserPasswordService:     forgotUserPasswordService,
		RecoverUserPasswordService:    recoverUserPasswordService,
	}
}

func (c *UserRegistrationController) RegisterUser(ctx *gin.Context) {
	var request user_http_requests.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.RegisterUserService.Execute(request.ToInput())
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *UserRegistrationController) VerifyUserRegistration(ctx *gin.Context) {
	var request user_http_requests.VerifyUserRegistrationRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.VerifyUserRegistrationService.Execute(request.ToInput())
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserRegistrationController) ForgotUserPassword(ctx *gin.Context) {
	var request user_http_requests.ForgotUserPasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.ForgotUserPasswordService.Execute(request.ToInput())
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserRegistrationController) RecoverUserPassword(ctx *gin.Context) {
	var request user_http_requests.RecoverUserPasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.RecoverUserPasswordService.Execute(request.ToInput())
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserRegistrationController) ConfigureRoutes(group *gin.RouterGroup) {
	g := group.Group("/user-registration")
	{
		g.POST("/", c.RegisterUser)
		g.POST("/verify", c.VerifyUserRegistration)
		g.POST("/forgot-password", c.ForgotUserPassword)
		g.POST("/recover-password", c.RecoverUserPassword)
	}
}
