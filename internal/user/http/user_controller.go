package user_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	user_http_requests "github.com/gabrielmrtt/taski/internal/user/http/requests"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	GetMeService                 *user_services.GetMeService
	ChangeUserPasswordService    *user_services.ChangeUserPasswordService
	UpdateUserCredentialsService *user_services.UpdateUserCredentialsService
	UpdateUserDataService        *user_services.UpdateUserDataService
	DeleteUserService            *user_services.DeleteUserService
}

func NewUserController(
	getMeService *user_services.GetMeService,
	changeUserPasswordService *user_services.ChangeUserPasswordService,
	updateUserCredentialsService *user_services.UpdateUserCredentialsService,
	updateUserDataService *user_services.UpdateUserDataService,
	deleteUserService *user_services.DeleteUserService,
) *UserController {
	return &UserController{
		GetMeService:                 getMeService,
		ChangeUserPasswordService:    changeUserPasswordService,
		UpdateUserCredentialsService: updateUserCredentialsService,
		UpdateUserDataService:        updateUserDataService,
		DeleteUserService:            deleteUserService,
	}
}

type GetMeResponse = core_http.HttpSuccessResponseWithData[user_core.UserDto]

// GetMe godoc
// @Summary Get me
// @Description Returns the authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} GetMeResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /me [get]
func (c *UserController) GetMe(ctx *gin.Context) {
	var request user_http_requests.GetMeRequest

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	input := request.ToInput()
	input.LoggedUserIdentity = authenticatedUserIdentity

	response, err := c.GetMeService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type ChangeUserPasswordResponse = core_http.EmptyHttpSuccessResponse

// ChangeUserPassword godoc
// @Summary Change user password
// @Description Change the authenticated user password.
// @Tags User
// @Accept json
// @Produce json
// @Param request body user_http_requests.ChangeUserPasswordRequest true "Request body"
// @Success 200 {object} ChangeUserPasswordResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /me/password [patch]
func (c *UserController) ChangeUserPassword(ctx *gin.Context) {
	var request user_http_requests.ChangeUserPasswordRequest

	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.UserIdentity = authenticatedUserIdentity
	err := c.ChangeUserPasswordService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type UpdateUserCredentialsResponse = core_http.EmptyHttpSuccessResponse

// UpdateUserCredentials godoc
// @Summary Update user credentials
// @Description Update the authenticated user credentials.
// @Tags User
// @Accept json
// @Produce json
// @Param request body user_http_requests.UpdateUserCredentialsRequest true "Request body"
// @Success 200 {object} UpdateUserCredentialsResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /me/credentials [put]
func (c *UserController) UpdateUserCredentials(ctx *gin.Context) {
	var request user_http_requests.UpdateUserCredentialsRequest

	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.UserIdentity = authenticatedUserIdentity
	err := c.UpdateUserCredentialsService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type UpdateUserDataResponse = core_http.EmptyHttpSuccessResponse

// UpdateUserData godoc
// @Summary Update user data
// @Description Update the authenticated user data.
// @Tags User
// @Accept mpfd
// @Produce json
// @Param request body user_http_requests.UpdateUserDataRequest true "Request body"
// @Success 200 {object} UpdateUserDataResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /me/data [put]
func (c *UserController) UpdateUserData(ctx *gin.Context) {
	var request user_http_requests.UpdateUserDataRequest

	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBind(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.UserIdentity = authenticatedUserIdentity
	err := c.UpdateUserDataService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteUserResponse = core_http.EmptyHttpSuccessResponse

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes the authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} DeleteUserResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /me [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	input := user_services.DeleteUserInput{
		UserIdentity: authenticatedUserIdentity,
	}
	err := c.DeleteUserService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *UserController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/me")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", c.GetMe)
		g.PATCH("/password", c.ChangeUserPassword)
		g.PUT("/credentials", c.UpdateUserCredentials)
		g.PUT("/data", c.UpdateUserData)
		g.DELETE("", c.DeleteUser)
	}

	return g
}
