package organizationhttp

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	organizationhttprequests "github.com/gabrielmrtt/taski/internal/organization/infra/http/requests"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	"github.com/gin-gonic/gin"
)

type OrganizationUserHandler struct {
	ListOrganizationUsersService      *organizationservice.ListOrganizationUsersService
	InviteUserToOrganizationService   *organizationservice.InviteUserToOrganizationService
	RemoveUserFromOrganizationService *organizationservice.RemoveUserFromOrganizationService
	GetOrganizationUserService        *organizationservice.GetOrganizationUserService
	UpdateOrganizationUserService     *organizationservice.UpdateOrganizationUserService
}

func NewOrganizationUserHandler(
	listOrganizationUsersService *organizationservice.ListOrganizationUsersService,
	inviteUserToOrganizationService *organizationservice.InviteUserToOrganizationService,
	removeUserFromOrganizationService *organizationservice.RemoveUserFromOrganizationService,
	getOrganizationUserService *organizationservice.GetOrganizationUserService,
	updateOrganizationUserService *organizationservice.UpdateOrganizationUserService,
) *OrganizationUserHandler {
	return &OrganizationUserHandler{
		ListOrganizationUsersService:      listOrganizationUsersService,
		InviteUserToOrganizationService:   inviteUserToOrganizationService,
		RemoveUserFromOrganizationService: removeUserFromOrganizationService,
		GetOrganizationUserService:        getOrganizationUserService,
		UpdateOrganizationUserService:     updateOrganizationUserService,
	}
}

type ListOrganizationUsersResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationUserDto]

// ListOrganizationUsers godoc
// @Summary List organization users
// @Description Lists all users in an organization.
// @Tags Organization User
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request query organizationhttprequests.ListOrganizationUsersRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListOrganizationUsersResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/user [get]
func (c *OrganizationUserHandler) ListOrganizationUsers(ctx *gin.Context) {
	var request organizationhttprequests.ListOrganizationUsersRequest
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity

	response, err := c.ListOrganizationUsersService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type InviteUserToOrganizationResponse = corehttp.EmptyHttpSuccessResponse

// InviteUserToOrganization godoc
// @Summary Invite user to organization
// @Description Invites a user to an organization.
// @Tags Organization User
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request body organizationhttprequests.InviteUserToOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} InviteUserToOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/user [post]
func (c *OrganizationUserHandler) InviteUserToOrganization(ctx *gin.Context) {
	var request organizationhttprequests.InviteUserToOrganizationRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	err := c.InviteUserToOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type RemoveUserFromOrganizationResponse = corehttp.EmptyHttpSuccessResponse

// RemoveUserFromOrganization godoc
// @Summary Remove user from organization
// @Description Removes a user from an organization.
// @Tags Organization User
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} RemoveUserFromOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/user/:userId [delete]
func (c *OrganizationUserHandler) RemoveUserFromOrganization(ctx *gin.Context) {
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	userIdentity := core.NewIdentityFromPublic(ctx.Param("userId"))

	input := organizationservice.RemoveUserFromOrganizationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.RemoveUserFromOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type GetOrganizationUserResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationUserDto]

// GetOrganizationUser godoc
// @Summary Get organization user
// @Description Returns an organization user.
// @Tags Organization User
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} GetOrganizationUserResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/user/:userId [get]
func (c *OrganizationUserHandler) GetOrganizationUser(ctx *gin.Context) {
	var request organizationhttprequests.GetOrganizationUserRequest

	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	userIdentity := core.NewIdentityFromPublic(ctx.Param("userId"))

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserIdentity = userIdentity

	response, err := c.GetOrganizationUserService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateOrganizationUserResponse = corehttp.EmptyHttpSuccessResponse

// UpdateOrganizationUser godoc
// @Summary Update organization user
// @Description Updates an organization user.
// @Tags Organization User
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} UpdateOrganizationUserResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/user/:userId [put]
func (c *OrganizationUserHandler) UpdateOrganizationUser(ctx *gin.Context) {
	var request organizationhttprequests.UpdateOrganizationUserRequest

	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	userIdentity := core.NewIdentityFromPublic(ctx.Param("userId"))

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserIdentity = userIdentity

	err := c.UpdateOrganizationUserService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *OrganizationUserHandler) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organizationId/user")
	{
		g.Use(userhttpmiddlewares.AuthMiddleware())

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:view"), c.ListOrganizationUsers)
		g.GET("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:view"), c.GetOrganizationUser)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:create"), c.InviteUserToOrganization)
		g.PUT("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:update"), c.UpdateOrganizationUser)
		g.DELETE("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:delete"), c.RemoveUserFromOrganization)
	}

	return g
}
