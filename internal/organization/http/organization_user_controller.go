package organization_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	organization_http_requests "github.com/gabrielmrtt/taski/internal/organization/http/requests"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

type OrganizationUserController struct {
	ListOrganizationUsersService            *organization_services.ListOrganizationUsersService
	InviteUserToOrganizationService         *organization_services.InviteUserToOrganizationService
	RemoveUserFromOrganizationService       *organization_services.RemoveUserFromOrganizationService
	AcceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService
	RefuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService
}

func NewOrganizationUserController(
	listOrganizationUsersService *organization_services.ListOrganizationUsersService,
	inviteUserToOrganizationService *organization_services.InviteUserToOrganizationService,
	removeUserFromOrganizationService *organization_services.RemoveUserFromOrganizationService,
	acceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService,
	refuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService,
) *OrganizationUserController {
	return &OrganizationUserController{
		ListOrganizationUsersService:            listOrganizationUsersService,
		InviteUserToOrganizationService:         inviteUserToOrganizationService,
		RemoveUserFromOrganizationService:       removeUserFromOrganizationService,
		AcceptOrganizationUserInvitationService: acceptOrganizationUserInvitationService,
		RefuseOrganizationUserInvitationService: refuseOrganizationUserInvitationService,
	}
}

func (c *OrganizationUserController) ListOrganizationUsers(ctx *gin.Context) {
	var request organization_http_requests.ListOrganizationUsersRequest
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	response, err := c.ListOrganizationUsersService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

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

func (c *OrganizationUserController) AcceptOrganizationUserInvitation(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	userIdentity := core.NewIdentityFromPublic(ctx.Param("user_id"))

	input := organization_services.AcceptOrganizationUserInvitationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.AcceptOrganizationUserInvitationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

func (c *OrganizationUserController) RefuseOrganizationUserInvitation(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	userIdentity := core.NewIdentityFromPublic(ctx.Param("user_id"))

	input := organization_services.RefuseOrganizationUserInvitationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.RefuseOrganizationUserInvitationService.Execute(input)
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

		g.GET("", organization_http_middlewares.BlockIfUserIsNotPartOfOrganization(), c.ListOrganizationUsers)

		g.POST("", organization_http_middlewares.BlockIfUserIsNotPartOfOrganization(), c.InviteUserToOrganization)
		g.DELETE("/:user_id", organization_http_middlewares.BlockIfUserIsNotPartOfOrganization(), c.RemoveUserFromOrganization)

		g.PATCH("/:user_id/accept-invitation", organization_http_middlewares.BlockIfUserIsNotSameOrganizationUser(), c.AcceptOrganizationUserInvitation)
		g.PATCH("/:user_id/refuse-invitation", organization_http_middlewares.BlockIfUserIsNotSameOrganizationUser(), c.RefuseOrganizationUserInvitation)
	}

	return g
}
