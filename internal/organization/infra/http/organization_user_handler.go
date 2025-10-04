package organizationhttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	organizationhttprequests "github.com/gabrielmrtt/taski/internal/organization/infra/http/requests"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
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
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))
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

	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))

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
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))
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

	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))
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

	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organizationId"))
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

func (c *OrganizationUserHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/organization/:organizationId/user")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:view", middlewareOptions), c.ListOrganizationUsers)
		g.GET("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:view", middlewareOptions), c.GetOrganizationUser)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:create", middlewareOptions), c.InviteUserToOrganization)
		g.PUT("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:update", middlewareOptions), c.UpdateOrganizationUser)
		g.DELETE("/:userId", organizationhttpmiddlewares.UserMustHavePermission("organizations:users:delete", middlewareOptions), c.RemoveUserFromOrganization)
	}

	return g
}
