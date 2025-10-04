package userhttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	user "github.com/gabrielmrtt/taski/internal/user"
	userhttprequests "github.com/gabrielmrtt/taski/internal/user/infra/http/requests"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	GetMeService                 *userservice.GetMeService
	ChangeUserPasswordService    *userservice.ChangeUserPasswordService
	UpdateUserCredentialsService *userservice.UpdateUserCredentialsService
	UpdateUserDataService        *userservice.UpdateUserDataService
	DeleteUserService            *userservice.DeleteUserService
}

func NewUserHandler(
	getMeService *userservice.GetMeService,
	changeUserPasswordService *userservice.ChangeUserPasswordService,
	updateUserCredentialsService *userservice.UpdateUserCredentialsService,
	updateUserDataService *userservice.UpdateUserDataService,
	deleteUserService *userservice.DeleteUserService,
) *UserHandler {
	return &UserHandler{
		GetMeService:                 getMeService,
		ChangeUserPasswordService:    changeUserPasswordService,
		UpdateUserCredentialsService: updateUserCredentialsService,
		UpdateUserDataService:        updateUserDataService,
		DeleteUserService:            deleteUserService,
	}
}

type GetMeResponse = corehttp.HttpSuccessResponseWithData[user.UserDto]

// GetMe godoc
// @Summary Get me
// @Description Returns the authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} GetMeResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /me [get]
func (c *UserHandler) GetMe(ctx *gin.Context) {
	var request userhttprequests.GetMeRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input userservice.GetMeInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.AuthenticatedUserIdentity = *authenticatedUserIdentity

	response, err := c.GetMeService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type ChangeUserPasswordResponse = corehttp.EmptyHttpSuccessResponse

// ChangeUserPassword godoc
// @Summary Change user password
// @Description Change the authenticated user password.
// @Tags User
// @Accept json
// @Produce json
// @Param request body userhttprequests.ChangeUserPasswordRequest true "Request body"
// @Success 200 {object} ChangeUserPasswordResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /me/password [patch]
func (c *UserHandler) ChangeUserPassword(ctx *gin.Context) {
	var request userhttprequests.ChangeUserPasswordRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input userservice.ChangeUserPasswordInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.UserIdentity = *authenticatedUserIdentity
	err := c.ChangeUserPasswordService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type UpdateUserCredentialsResponse = corehttp.EmptyHttpSuccessResponse

// UpdateUserCredentials godoc
// @Summary Update user credentials
// @Description Update the authenticated user credentials.
// @Tags User
// @Accept json
// @Produce json
// @Param request body userhttprequests.UpdateUserCredentialsRequest true "Request body"
// @Success 200 {object} UpdateUserCredentialsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /me/credentials [put]
func (c *UserHandler) UpdateUserCredentials(ctx *gin.Context) {
	var request userhttprequests.UpdateUserCredentialsRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input userservice.UpdateUserCredentialsInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.UserIdentity = *authenticatedUserIdentity

	err := c.UpdateUserCredentialsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type UpdateUserDataResponse = corehttp.EmptyHttpSuccessResponse

// UpdateUserData godoc
// @Summary Update user data
// @Description Update the authenticated user data.
// @Tags User
// @Accept mpfd
// @Produce json
// @Param request body userhttprequests.UpdateUserDataRequest true "Request body"
// @Success 200 {object} UpdateUserDataResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /me/data [put]
func (c *UserHandler) UpdateUserData(ctx *gin.Context) {
	var request userhttprequests.UpdateUserDataRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input userservice.UpdateUserDataInput

	if err := ctx.ShouldBind(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.UserIdentity = *authenticatedUserIdentity

	err := c.UpdateUserDataService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteUserResponse = corehttp.EmptyHttpSuccessResponse

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes the authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} DeleteUserResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /me [delete]
func (c *UserHandler) DeleteUser(ctx *gin.Context) {
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input userservice.DeleteUserInput = userservice.DeleteUserInput{
		UserIdentity: *authenticatedUserIdentity,
	}

	err := c.DeleteUserService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/me")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", c.GetMe)
		g.PATCH("/password", c.ChangeUserPassword)
		g.PUT("/credentials", c.UpdateUserCredentials)
		g.PUT("/data", c.UpdateUserData)
		g.DELETE("", c.DeleteUser)
	}

	return g
}
