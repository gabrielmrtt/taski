package userhttp

import (
	"net/http"

	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	user "github.com/gabrielmrtt/taski/internal/user"
	userhttprequests "github.com/gabrielmrtt/taski/internal/user/infra/http/requests"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gin-gonic/gin"
)

type UserRegistrationHandler struct {
	RegisterUserService           *userservice.RegisterUserService
	VerifyUserRegistrationService *userservice.VerifyUserRegistrationService
	ForgotUserPasswordService     *userservice.ForgotUserPasswordService
	RecoverUserPasswordService    *userservice.RecoverUserPasswordService
}

func NewUserRegistrationHandler(
	registerUserService *userservice.RegisterUserService,
	verifyUserRegistrationService *userservice.VerifyUserRegistrationService,
	forgotUserPasswordService *userservice.ForgotUserPasswordService,
	recoverUserPasswordService *userservice.RecoverUserPasswordService,
) *UserRegistrationHandler {
	return &UserRegistrationHandler{
		RegisterUserService:           registerUserService,
		VerifyUserRegistrationService: verifyUserRegistrationService,
		ForgotUserPasswordService:     forgotUserPasswordService,
		RecoverUserPasswordService:    recoverUserPasswordService,
	}
}

type RegisterUserResponse = corehttp.HttpSuccessResponseWithData[user.UserDto]

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user and creates a new user registration token. To activate the user, you need to verify the user registration using their user registration token.
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body userhttprequests.RegisterUserRequest true "Request body"
// @Success 200 {object} RegisterUserResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Failure 409 {object} corehttp.HttpErrorResponse
// @Router /user-registration [post]
func (c *UserRegistrationHandler) RegisterUser(ctx *gin.Context) {
	var request userhttprequests.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.RegisterUserService.Execute(request.ToInput())
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type VerifyUserRegistrationResponse = corehttp.EmptyHttpSuccessResponse

// VerifyUserRegistration godoc
// @Summary Verify user registration
// @Description Verifies an user registration using an user registration token. After verifying, the user will be activated and ready to use other services.
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body userhttprequests.VerifyUserRegistrationRequest true "Request body"
// @Success 200 {object} VerifyUserRegistrationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 409 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /user-registration/verify [post]
func (c *UserRegistrationHandler) VerifyUserRegistration(ctx *gin.Context) {
	var request userhttprequests.VerifyUserRegistrationRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.VerifyUserRegistrationService.Execute(request.ToInput())
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type ForgotUserPasswordResponse = corehttp.EmptyHttpSuccessResponse

// ForgotUserPassword godoc
// @Summary Forgot user password
// @Description Creates a new password recovery token.
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body userhttprequests.ForgotUserPasswordRequest true "Request body"
// @Success 200 {object} ForgotUserPasswordResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 409 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /user-registration/forgot-password [post]
func (c *UserRegistrationHandler) ForgotUserPassword(ctx *gin.Context) {
	var request userhttprequests.ForgotUserPasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.ForgotUserPasswordService.Execute(request.ToInput())
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type RecoverUserPasswordResponse = corehttp.EmptyHttpSuccessResponse

// RecoverUserPassword godoc
// @Summary Recover user password
// @Description Recovers an user password using a password recovery token.
// @Tags User Registration
// @Accept json
// @Produce json
// @Param request body userhttprequests.RecoverUserPasswordRequest true "Request body"
// @Success 200 {object} RecoverUserPasswordResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 409 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /user-registration/recover-password [post]
func (c *UserRegistrationHandler) RecoverUserPassword(ctx *gin.Context) {
	var request userhttprequests.RecoverUserPasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	err := c.RecoverUserPasswordService.Execute(request.ToInput())
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserRegistrationHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	g := options.RouterGroup.Group("/user-registration")
	{
		g.POST("/", c.RegisterUser)
		g.POST("/verify", c.VerifyUserRegistration)
		g.POST("/forgot-password", c.ForgotUserPassword)
		g.POST("/recover-password", c.RecoverUserPassword)
	}

	return g
}
