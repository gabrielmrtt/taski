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
	return
}

type GetOrganizationResponse = core_http.HttpSuccessResponseWithData[organization_core.OrganizationDto]

// GetOrganization godoc
// @Summary Get an organization
// @Description Returns an accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Produce json
// @Success 200 {object} GetOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id [get]
func (c *OrganizationController) GetOrganization(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	input := organization_services.GetOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	response, err := c.GetOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
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
	return
}

type UpdateOrganizationResponse = core_http.EmptyHttpSuccessResponse

// UpdateOrganization godoc
// @Summary Update an organization
// @Description Updates an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Param request body organization_http_requests.UpdateOrganizationRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id [put]
func (c *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	var request organization_http_requests.UpdateOrganizationRequest
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
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
	return
}

type DeleteOrganizationResponse = core_http.EmptyHttpSuccessResponse

// DeleteOrganization godoc
// @Summary Delete an organization
// @Description Deletes an existing and accessible organization by the authenticated user and the organization ID.
// @Tags Organization
// @Accept json
// @Param organization_id path string true "Organization ID"
// @Produce json
// @Success 200 {object} DeleteOrganizationResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organization_id [delete]
func (c *OrganizationController) DeleteOrganization(ctx *gin.Context) {
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))

	input := organization_services.DeleteOrganizationInput{
		OrganizationIdentity: organizationIdentity,
	}

	err := c.DeleteOrganizationService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
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
	return
}

func (c *OrganizationController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", c.ListOrganizations)
		g.GET("/:organization_id", organization_http_middlewares.UserMustHavePermission("organizations:view"), c.GetOrganization)
		g.POST("", c.CreateOrganization)
		g.PUT("/:organization_id", organization_http_middlewares.UserMustHavePermission("organizations:update"), c.UpdateOrganization)
		g.DELETE("/:organization_id", organization_http_middlewares.UserMustHavePermission("organizations:delete"), c.DeleteOrganization)
	}

	return g
}
