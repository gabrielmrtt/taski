package organization_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	organization_http_requests "github.com/gabrielmrtt/taski/internal/organization/http/requests"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gin-gonic/gin"
)

type OrganizationUserController struct {
	InviteUserToOrganizationService         *organization_services.InviteUserToOrganizationService
	RemoveUserFromOrganizationService       *organization_services.RemoveUserFromOrganizationService
	AcceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService
	RefuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService
}

func NewOrganizationUserController(
	inviteUserToOrganizationService *organization_services.InviteUserToOrganizationService,
	removeUserFromOrganizationService *organization_services.RemoveUserFromOrganizationService,
	acceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService,
	refuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService,
) *OrganizationUserController {
	return &OrganizationUserController{
		InviteUserToOrganizationService:         inviteUserToOrganizationService,
		RemoveUserFromOrganizationService:       removeUserFromOrganizationService,
		AcceptOrganizationUserInvitationService: acceptOrganizationUserInvitationService,
		RefuseOrganizationUserInvitationService: refuseOrganizationUserInvitationService,
	}
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
	group.POST("/organization/:organization_id/user", organization_http_middlewares.BlockIfUserIsNotPartOfOrganization(), c.InviteUserToOrganization)
	group.DELETE("/organization/:organization_id/user/:user_id", organization_http_middlewares.BlockIfUserIsNotPartOfOrganization(), c.RemoveUserFromOrganization)

	group.PATCH("/organization/:organization_id/user/:user_id/accept-invitation", organization_http_middlewares.BlockIfUserIsNotSameOrganizationUser(), c.AcceptOrganizationUserInvitation)
	group.PATCH("/organization/:organization_id/user/:user_id/refuse-invitation", organization_http_middlewares.BlockIfUserIsNotSameOrganizationUser(), c.RefuseOrganizationUserInvitation)

	return group
}
