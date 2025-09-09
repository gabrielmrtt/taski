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

type OrganizationController struct {
	ListOrganizationsService  *organization_services.ListOrganizationsService
	GetOrganizationService    *organization_services.GetOrganizationService
	CreateOrganizationService *organization_services.CreateOrganizationService
	UpdateOrganizationService *organization_services.UpdateOrganizationService
	DeleteOrganizationService *organization_services.DeleteOrganizationService
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

func (c *OrganizationController) ConfigureRoutes(group *gin.RouterGroup) {
	g := group.Group("/organization")
	{
		g.Use(user_http_middlewares.AuthMiddleware())
		g.Use(organization_http_middlewares.BlockIfUserIsNotPartOfOrganization())

		g.GET("", c.ListOrganizations)
		g.GET("/:organization_id", c.GetOrganization)
		g.POST("", c.CreateOrganization)
		g.PUT("/:organization_id", c.UpdateOrganization)
		g.DELETE("/:organization_id", c.DeleteOrganization)
	}
}
