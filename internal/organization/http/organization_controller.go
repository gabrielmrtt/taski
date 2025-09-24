package organization_http

import (
	"net/http"

	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	organization_http_requests "github.com/gabrielmrtt/taski/internal/organization/http/requests"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	ListOrganizationsService         *organization_services.ListOrganizationsService
	GetOrganizationService           *organization_services.GetOrganizationService
	CreateOrganizationService        *organization_services.CreateOrganizationService
	UpdateOrganizationService        *organization_services.UpdateOrganizationService
	DeleteOrganizationService        *organization_services.DeleteOrganizationService
	ListMyOrganizationInvitesService *organization_services.ListMyOrganizationInvitesService
}

func NewOrganizationController(
	listOrganizationsService *organization_services.ListOrganizationsService,
	getOrganizationService *organization_services.GetOrganizationService,
	createOrganizationService *organization_services.CreateOrganizationService,
	updateOrganizationService *organization_services.UpdateOrganizationService,
	deleteOrganizationService *organization_services.DeleteOrganizationService,
) *OrganizationController {
	return &OrganizationController{
		ListOrganizationsService:  listOrganizationsService,
		GetOrganizationService:    getOrganizationService,
		CreateOrganizationService: createOrganizationService,
		UpdateOrganizationService: updateOrganizationService,
		DeleteOrganizationService: deleteOrganizationService,
	}
}

type ListOrganizationsResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationDto]

// ListOrganizations godoc
// @Summary List organizations
// @Description List accessible organizations by the authenticated user (created by them or organizations they are part of).
// @Tags Organization
// @Accept json
// @Param request query organization_http_requests.ListOrganizationsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListOrganizationsResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization [get]
func (c *OrganizationController) ListOrganizations(ctx *gin.Context) {
	var request organization_http_requests.ListOrganizationsRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.Filters.LoggedUserIdentity = &authenticatedUserIdentity

	response, err := c.ListOrganizationsService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetOrganizationResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationDto]

// GetOrganization godoc
// @Summary Get an organization
// @Description Returns an accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Produce json
// @Success 200 {object} GetOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId [get]
func (c *OrganizationController) GetOrganization(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	input := organization_services.GetOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	response, err := c.GetOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateOrganizationResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationDto]

// CreateOrganization godoc
// @Summary Create an organization
// @Description Creates a new organization.
// @Tags Organization
// @Accept json
// @Param request body organization_http_requests.CreateOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization [post]
func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	var request organization_http_requests.CreateOrganizationRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateOrganizationResponse = core_http.EmptyHttpSuccessResponse

// UpdateOrganization godoc
// @Summary Update an organization
// @Description Updates an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request body organization_http_requests.UpdateOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId [put]
func (c *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	var request organization_http_requests.UpdateOrganizationRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserEditorIdentity = authenticatedUserIdentity

	err := c.UpdateOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteOrganizationResponse = core_http.EmptyHttpSuccessResponse

// DeleteOrganization godoc
// @Summary Delete an organization
// @Description Deletes an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Produce json
// @Success 200 {object} DeleteOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId [delete]
func (c *OrganizationController) DeleteOrganization(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	input := organization_services.DeleteOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	err := c.DeleteOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *OrganizationController) ListMyOrganizationInvites(ctx *gin.Context) {
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

func (c *OrganizationController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", c.ListOrganizations)
		g.POST("", c.CreateOrganization)
		g.GET("/:organizationId", organization_http_middlewares.UserMustHavePermission("organizations:view"), c.GetOrganization)
		g.PUT("/:organizationId", organization_http_middlewares.UserMustHavePermission("organizations:update"), c.UpdateOrganization)
		g.DELETE("/:organizationId", organization_http_middlewares.UserMustHavePermission("organizations:delete"), c.DeleteOrganization)
	}

	return g
}
