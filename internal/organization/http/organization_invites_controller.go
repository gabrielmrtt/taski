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

type OrganizationInvitesController struct {
	ListMyOrganizationInvitesService        *organization_services.ListMyOrganizationInvitesService
	AcceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService
	RefuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService
}

func NewOrganizationInvitesController(
	listMyOrganizationInvitesService *organization_services.ListMyOrganizationInvitesService,
	acceptOrganizationUserInvitationService *organization_services.AcceptOrganizationUserInvitationService,
	refuseOrganizationUserInvitationService *organization_services.RefuseOrganizationUserInvitationService,
) *OrganizationInvitesController {
	return &OrganizationInvitesController{
		ListMyOrganizationInvitesService:        listMyOrganizationInvitesService,
		AcceptOrganizationUserInvitationService: acceptOrganizationUserInvitationService,
		RefuseOrganizationUserInvitationService: refuseOrganizationUserInvitationService,
	}
}

type ListMyOrganizationInvitesResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationDto]

// ListMyOrganizationInvites godoc
// @Summary List my organization invites
// @Description Returns organizations the authenticated user has been invited to.
// @Tags Organization Invites
// @Accept json
// @Param request query organization_http_requests.ListMyOrganizationInvitesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListMyOrganizationInvitesResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization-invites [get]
func (c *OrganizationInvitesController) ListMyOrganizationInvites(ctx *gin.Context) {
	var request organization_http_requests.ListMyOrganizationInvitesRequest

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.AuthenticatedUserIdentity = user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	response, err := c.ListMyOrganizationInvitesService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type AcceptOrganizationUserInvitationResponse = core_http.EmptyHttpSuccessResponse

// AcceptOrganizationUserInvitation godoc
// @Summary Accept organization user invitation
// @Description Accept organization user invitation
// @Tags Organization Invites
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} AcceptOrganizationUserInvitationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization-invites/:organizationId/user/:userId/accept-invitation [patch]
func (c *OrganizationInvitesController) AcceptOrganizationUserInvitation(ctx *gin.Context) {
	var organizationIdentity core.Identity
	var userIdentity core.Identity

	organizationId := ctx.Param("organizationId")
	userId := ctx.Param("userId")

	organizationIdentity = core.NewIdentityFromPublic(organizationId)
	userIdentity = core.NewIdentityFromPublic(userId)

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
}

type RefuseOrganizationUserInvitationResponse = core_http.EmptyHttpSuccessResponse

// RefuseOrganizationUserInvitation godoc
// @Summary Refuse organization user invitation
// @Description Refuse organization user invitation
// @Tags Organization Invites
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} RefuseOrganizationUserInvitationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization-invites/:organizationId/user/:userId/refuse-invitation [patch]
func (c *OrganizationInvitesController) RefuseOrganizationUserInvitation(ctx *gin.Context) {
	var organizationIdentity core.Identity
	var userIdentity core.Identity

	organizationId := ctx.Param("organizationId")
	userId := ctx.Param("userId")

	organizationIdentity = core.NewIdentityFromPublic(organizationId)
	userIdentity = core.NewIdentityFromPublic(userId)

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
}

func (c *OrganizationInvitesController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization-invites")
	{
		g.Use(user_http_middlewares.AuthMiddleware())
		g.Use(organization_http_middlewares.UserMustBeSame())

		g.GET("", c.ListMyOrganizationInvites)
		g.PATCH("/:organizationId/user/:userId/accept-invitation", c.AcceptOrganizationUserInvitation)
		g.PATCH("/:organizationId/user/:userId/refuse-invitation", c.RefuseOrganizationUserInvitation)
	}

	return g
}
