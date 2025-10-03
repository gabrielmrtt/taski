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

type OrganizationInvitesHandler struct {
	ListMyOrganizationInvitesService        *organizationservice.ListMyOrganizationInvitesService
	AcceptOrganizationUserInvitationService *organizationservice.AcceptOrganizationUserInvitationService
	RefuseOrganizationUserInvitationService *organizationservice.RefuseOrganizationUserInvitationService
}

func NewOrganizationInvitesHandler(
	listMyOrganizationInvitesService *organizationservice.ListMyOrganizationInvitesService,
	acceptOrganizationUserInvitationService *organizationservice.AcceptOrganizationUserInvitationService,
	refuseOrganizationUserInvitationService *organizationservice.RefuseOrganizationUserInvitationService,
) *OrganizationInvitesHandler {
	return &OrganizationInvitesHandler{
		ListMyOrganizationInvitesService:        listMyOrganizationInvitesService,
		AcceptOrganizationUserInvitationService: acceptOrganizationUserInvitationService,
		RefuseOrganizationUserInvitationService: refuseOrganizationUserInvitationService,
	}
}

type ListMyOrganizationInvitesResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationDto]

// ListMyOrganizationInvites godoc
// @Summary List my organization invites
// @Description Returns organizations the authenticated user has been invited to.
// @Tags Organization Invites
// @Accept json
// @Param request query organizationhttprequests.ListMyOrganizationInvitesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListMyOrganizationInvitesResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization-invites [get]
func (c *OrganizationInvitesHandler) ListMyOrganizationInvites(ctx *gin.Context) {
	var request organizationhttprequests.ListMyOrganizationInvitesRequest

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.AuthenticatedUserIdentity = userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	response, err := c.ListMyOrganizationInvitesService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type AcceptOrganizationUserInvitationResponse = corehttp.EmptyHttpSuccessResponse

// AcceptOrganizationUserInvitation godoc
// @Summary Accept organization user invitation
// @Description Accept organization user invitation
// @Tags Organization Invites
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} AcceptOrganizationUserInvitationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization-invites/:organizationId/user/:userId/accept-invitation [patch]
func (c *OrganizationInvitesHandler) AcceptOrganizationUserInvitation(ctx *gin.Context) {
	var organizationIdentity core.Identity
	var userIdentity core.Identity

	organizationId := ctx.Param("organizationId")
	userId := ctx.Param("userId")

	organizationIdentity = core.NewIdentityFromPublic(organizationId)
	userIdentity = core.NewIdentityFromPublic(userId)

	input := organizationservice.AcceptOrganizationUserInvitationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.AcceptOrganizationUserInvitationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type RefuseOrganizationUserInvitationResponse = corehttp.EmptyHttpSuccessResponse

// RefuseOrganizationUserInvitation godoc
// @Summary Refuse organization user invitation
// @Description Refuse organization user invitation
// @Tags Organization Invites
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} RefuseOrganizationUserInvitationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization-invites/:organizationId/user/:userId/refuse-invitation [patch]
func (c *OrganizationInvitesHandler) RefuseOrganizationUserInvitation(ctx *gin.Context) {
	var organizationIdentity core.Identity
	var userIdentity core.Identity

	organizationId := ctx.Param("organizationId")
	userId := ctx.Param("userId")

	organizationIdentity = core.NewIdentityFromPublic(organizationId)
	userIdentity = core.NewIdentityFromPublic(userId)

	input := organizationservice.RefuseOrganizationUserInvitationInput{
		OrganizationIdentity: organizationIdentity,
		UserIdentity:         userIdentity,
	}

	err := c.RefuseOrganizationUserInvitationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *OrganizationInvitesHandler) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization-invites")
	{
		g.Use(userhttpmiddlewares.AuthMiddleware())
		g.Use(organizationhttpmiddlewares.UserMustBeSame())

		g.GET("", c.ListMyOrganizationInvites)
		g.PATCH("/:organizationId/user/:userId/accept-invitation", c.AcceptOrganizationUserInvitation)
		g.PATCH("/:organizationId/user/:userId/refuse-invitation", c.RefuseOrganizationUserInvitation)
	}

	return g
}
