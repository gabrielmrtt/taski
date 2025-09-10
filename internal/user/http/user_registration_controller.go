package user_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_core "github.com/gabrielmrtt/taski/internal/user"
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

type RegisterUserResponse = core_http.HttpSuccessResponseWithData[user_core.UserDto]

// RegisterUser godoc
// @Summary Register a new user
// @Schemes
// @Description Register a new user
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body user_http_requests.RegisterUserRequest true "Register User Request"
// @Success 200 {object} RegisterUserResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Failure 409 {object} core_http.HttpErrorResponse
// @Router /user-registration [post]
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

type VerifyUserRegistrationResponse = core_http.EmptyHttpSuccessResponse

// VerifyUserRegistration godoc
// @Summary Verify user registration
// @Schemes
// @Description Verify user registration
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body user_http_requests.VerifyUserRegistrationRequest true "Verify User Registration Request"
// @Success 200 {object} VerifyUserRegistrationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Failure 409 {object} core_http.HttpErrorResponse
// @Router /user-registration/verify [post]
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

type ForgotUserPasswordResponse = core_http.EmptyHttpSuccessResponse

// ForgotUserPassword godoc
// @Summary Forgot user password
// @Schemes
// @Description Forgot user password
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body user_http_requests.ForgotUserPasswordRequest true "Forgot User Password Request"
// @Success 200 {object} ForgotUserPasswordResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Failure 409 {object} core_http.HttpErrorResponse
// @Router /user-registration/forgot-password [post]
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

type RecoverUserPasswordResponse = core_http.EmptyHttpSuccessResponse

// RecoverUserPassword godoc
// @Summary Recover user password
// @Schemes
// @Description Recover user password
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body user_http_requests.RecoverUserPasswordRequest true "Recover User Password Request"
// @Success 200 {object} RecoverUserPasswordResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Failure 409 {object} core_http.HttpErrorResponse
// @Router /user-registration/recover-password [post]
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

func (c *UserRegistrationController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/user-registration")
	{
		g.POST("/", c.RegisterUser)
		g.POST("/verify", c.VerifyUserRegistration)
		g.POST("/forgot-password", c.ForgotUserPassword)
		g.POST("/recover-password", c.RecoverUserPassword)
	}

	return g
}
