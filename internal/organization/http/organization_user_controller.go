package organization_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	organization_http_requests "github.com/gabrielmrtt/taski/internal/organization/http/requests"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

type OrganizationUserController struct {
	ListOrganizationUsersService      *organization_services.ListOrganizationUsersService
	InviteUserToOrganizationService   *organization_services.InviteUserToOrganizationService
	RemoveUserFromOrganizationService *organization_services.RemoveUserFromOrganizationService
	GetOrganizationUserService        *organization_services.GetOrganizationUserService
	UpdateOrganizationUserService     *organization_services.UpdateOrganizationUserService
}

func NewOrganizationUserController(
	listOrganizationUsersService *organization_services.ListOrganizationUsersService,
	inviteUserToOrganizationService *organization_services.InviteUserToOrganizationService,
	removeUserFromOrganizationService *organization_services.RemoveUserFromOrganizationService,
	getOrganizationUserService *organization_services.GetOrganizationUserService,
	updateOrganizationUserService *organization_services.UpdateOrganizationUserService,
) *OrganizationUserController {
	return &OrganizationUserController{
		ListOrganizationUsersService:      listOrganizationUsersService,
		InviteUserToOrganizationService:   inviteUserToOrganizationService,
		RemoveUserFromOrganizationService: removeUserFromOrganizationService,
		GetOrganizationUserService:        getOrganizationUserService,
		UpdateOrganizationUserService:     updateOrganizationUserService,
	}
}

type ListOrganizationUsersResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationUserDto]

// ListOrganizationUsers godoc
// @Summary List organization users
// @Description Lists all users in an organization.
// @Tags Organization User
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param request query organization_http_requests.ListOrganizationUsersRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListOrganizationUsersResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id/user [get]
func (c *OrganizationUserController) ListOrganizationUsers(ctx *gin.Context) {
	var request organization_http_requests.ListOrganizationUsersRequest
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity

	response, err := c.ListOrganizationUsersService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

type InviteUserToOrganizationResponse = core_http.EmptyHttpSuccessResponse

// InviteUserToOrganization godoc
// @Summary Invite user to organization
// @Description Invites a user to an organization.
// @Tags Organization User
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param request body organization_http_requests.InviteUserToOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} InviteUserToOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id/user [post]
func (c *OrganizationUserController) InviteUserToOrganization(ctx *gin.Context) {
	var request organization_http_requests.InviteUserToOrganizationRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	err := c.InviteUserToOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

type RemoveUserFromOrganizationResponse = core_http.EmptyHttpSuccessResponse

// RemoveUserFromOrganization godoc
// @Summary Remove user from organization
// @Description Removes a user from an organization.
// @Tags Organization User
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} RemoveUserFromOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id/user/:user_id [delete]
func (c *OrganizationUserController) RemoveUserFromOrganization(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	userIdentity := core.NewIdentityFromPublic(ctx.Param("user_id"))

	input := organization_services.RemoveUserFromOrganizationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.RemoveUserFromOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

type GetOrganizationUserResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationUserDto]

// GetOrganizationUser godoc
// @Summary Get organization user
// @Description Returns an organization user.
// @Tags Organization User
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} GetOrganizationUserResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id/user/:user_id [get]
func (c *OrganizationUserController) GetOrganizationUser(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	userIdentity := core.NewIdentityFromPublic(ctx.Param("user_id"))

	input := organization_services.GetOrganizationUserInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	response, err := c.GetOrganizationUserService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

type UpdateOrganizationUserResponse = core_http.EmptyHttpSuccessResponse

// UpdateOrganizationUser godoc
// @Summary Update organization user
// @Description Updates an organization user.
// @Tags Organization User
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} UpdateOrganizationUserResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id/user/:user_id [put]
func (c *OrganizationUserController) UpdateOrganizationUser(ctx *gin.Context) {
	var request organization_http_requests.UpdateOrganizationUserRequest

	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	userIdentity := core.NewIdentityFromPublic(ctx.Param("user_id"))

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserIdentity = userIdentity

	err := c.UpdateOrganizationUserService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

func (c *OrganizationUserController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organization_id/user")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", organization_http_middlewares.UserMustHavePermission("organizations:users:view"), c.ListOrganizationUsers)
		g.GET("/:user_id", organization_http_middlewares.UserMustHavePermission("organizations:users:view"), c.GetOrganizationUser)
		g.POST("", organization_http_middlewares.UserMustHavePermission("organizations:users:create"), c.InviteUserToOrganization)
		g.PUT("/:user_id", organization_http_middlewares.UserMustHavePermission("organizations:users:update"), c.UpdateOrganizationUser)
		g.DELETE("/:user_id", organization_http_middlewares.UserMustHavePermission("organizations:users:delete"), c.RemoveUserFromOrganization)
	}

	return g
}
