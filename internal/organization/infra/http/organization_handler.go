package organizationhttp

import (
	"net/http"

	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	organizationhttprequests "github.com/gabrielmrtt/taski/internal/organization/infra/http/requests"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	ListOrganizationsService         *organizationservice.ListOrganizationsService
	GetOrganizationService           *organizationservice.GetOrganizationService
	CreateOrganizationService        *organizationservice.CreateOrganizationService
	UpdateOrganizationService        *organizationservice.UpdateOrganizationService
	DeleteOrganizationService        *organizationservice.DeleteOrganizationService
	ListMyOrganizationInvitesService *organizationservice.ListMyOrganizationInvitesService
}

func NewOrganizationHandler(
	listOrganizationsService *organizationservice.ListOrganizationsService,
	getOrganizationService *organizationservice.GetOrganizationService,
	createOrganizationService *organizationservice.CreateOrganizationService,
	updateOrganizationService *organizationservice.UpdateOrganizationService,
	deleteOrganizationService *organizationservice.DeleteOrganizationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		ListOrganizationsService:  listOrganizationsService,
		GetOrganizationService:    getOrganizationService,
		CreateOrganizationService: createOrganizationService,
		UpdateOrganizationService: updateOrganizationService,
		DeleteOrganizationService: deleteOrganizationService,
	}
}

type ListOrganizationsResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationDto]

// ListOrganizations godoc
// @Summary List organizations
// @Description List accessible organizations by the authenticated user (created by them or organizations they are part of).
// @Tags Organization
// @Accept json
// @Param request query organizationhttprequests.ListOrganizationsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListOrganizationsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization [get]
func (c *OrganizationHandler) ListOrganizations(ctx *gin.Context) {
	var request organizationhttprequests.ListOrganizationsRequest
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.Filters.LoggedUserIdentity = &authenticatedUserIdentity

	response, err := c.ListOrganizationsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetOrganizationResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationDto]

// GetOrganization godoc
// @Summary Get an organization
// @Description Returns an accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Produce json
// @Success 200 {object} GetOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId [get]
func (c *OrganizationHandler) GetOrganization(ctx *gin.Context) {
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	input := organizationservice.GetOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	response, err := c.GetOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateOrganizationResponse = corehttp.HttpSuccessResponseWithData[organization.OrganizationDto]

// CreateOrganization godoc
// @Summary Create an organization
// @Description Creates a new organization.
// @Tags Organization
// @Accept json
// @Param request body organizationhttprequests.CreateOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization [post]
func (c *OrganizationHandler) CreateOrganization(ctx *gin.Context) {
	var request organizationhttprequests.CreateOrganizationRequest
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateOrganizationResponse = corehttp.EmptyHttpSuccessResponse

// UpdateOrganization godoc
// @Summary Update an organization
// @Description Updates an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request body organizationhttprequests.UpdateOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId [put]
func (c *OrganizationHandler) UpdateOrganization(ctx *gin.Context) {
	var request organizationhttprequests.UpdateOrganizationRequest
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserEditorIdentity = authenticatedUserIdentity

	err := c.UpdateOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteOrganizationResponse = corehttp.EmptyHttpSuccessResponse

// DeleteOrganization godoc
// @Summary Delete an organization
// @Description Deletes an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Produce json
// @Success 200 {object} DeleteOrganizationResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId [delete]
func (c *OrganizationHandler) DeleteOrganization(ctx *gin.Context) {
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)

	input := organizationservice.DeleteOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	err := c.DeleteOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *OrganizationHandler) ListMyOrganizationInvites(ctx *gin.Context) {
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

func (c *OrganizationHandler) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization")
	{
		g.Use(userhttpmiddlewares.AuthMiddleware())

		g.GET("", c.ListOrganizations)
		g.POST("", c.CreateOrganization)
		g.GET("/:organizationId", organizationhttpmiddlewares.UserMustHavePermission("organizations:view"), c.GetOrganization)
		g.PUT("/:organizationId", organizationhttpmiddlewares.UserMustHavePermission("organizations:update"), c.UpdateOrganization)
		g.DELETE("/:organizationId", organizationhttpmiddlewares.UserMustHavePermission("organizations:delete"), c.DeleteOrganization)
	}

	return g
}
