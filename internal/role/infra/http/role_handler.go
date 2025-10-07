package rolehttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/role"
	rolehttprequests "github.com/gabrielmrtt/taski/internal/role/infra/http/requests"
	roleservice "github.com/gabrielmrtt/taski/internal/role/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	CreateRoleService *roleservice.CreateRoleService
	UpdateRoleService *roleservice.UpdateRoleService
	DeleteRoleService *roleservice.DeleteRoleService
	ListRolesService  *roleservice.ListRolesService
}

func NewRoleHandler(
	createRoleService *roleservice.CreateRoleService,
	updateRoleService *roleservice.UpdateRoleService,
	deleteRoleService *roleservice.DeleteRoleService,
	listRolesService *roleservice.ListRolesService,
) *RoleHandler {
	return &RoleHandler{
		CreateRoleService: createRoleService,
		UpdateRoleService: updateRoleService,
		DeleteRoleService: deleteRoleService,
		ListRolesService:  listRolesService,
	}
}

type CreateRoleResponse = corehttp.HttpSuccessResponseWithData[role.RoleDto]

// CreateRole godoc
// @Summary Create a role in an organization
// @Description Creates a new role in an organization.
// @Tags Role
// @Accept json
// @Param request body rolehttprequests.CreateRoleRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateRoleResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /role [post]
func (c *RoleHandler) CreateRole(ctx *gin.Context) {
	var request rolehttprequests.CreateRoleRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var input roleservice.CreateRoleInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

	response, err := c.CreateRoleService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateRoleResponse = corehttp.EmptyHttpSuccessResponse

// UpdateRole godoc
// @Summary Update a role in an organization
// @Description Updates an existing role in an organization.
// @Tags Role
// @Accept json
// @Produce json
// @Param roleId path string true "Role ID"
// @Param request body rolehttprequests.UpdateRoleRequest true "Request body"
// @Success 200 {object} UpdateRoleResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /role/:roleId [put]
func (c *RoleHandler) UpdateRole(ctx *gin.Context) {
	var request rolehttprequests.UpdateRoleRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var roleIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("roleId"))
	var input roleservice.UpdateRoleInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserEditorIdentity = *authenticatedUserIdentity
	input.RoleIdentity = roleIdentity

	err := c.UpdateRoleService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteRoleResponse = corehttp.EmptyHttpSuccessResponse

// DeleteRole godoc
// @Summary Delete a role in an organization
// @Description Deletes an existing role in an organization.
// @Tags Role
// @Accept json
// @Param roleId path string true "Role ID"
// @Produce json
// @Success 200 {object} DeleteRoleResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /role/:roleId [delete]
func (c *RoleHandler) DeleteRole(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var roleIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("roleId"))
	var input roleservice.DeleteRoleInput = roleservice.DeleteRoleInput{
		RoleIdentity:         roleIdentity,
		OrganizationIdentity: *organizationIdentity,
	}

	err := c.DeleteRoleService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type ListRolesResponse = corehttp.HttpSuccessResponseWithData[role.RoleDto]

// ListRoles godoc
// @Summary List roles in an organization
// @Description Lists all existing roles in an organization.
// @Tags Role
// @Accept json
// @Param request query rolehttprequests.ListRolesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListRolesResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/role [get]
func (c *RoleHandler) ListRoles(ctx *gin.Context) {
	var request rolehttprequests.ListRolesRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var input roleservice.ListRolesInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity

	response, err := c.ListRolesService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *RoleHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/role")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("roles:view", middlewareOptions), c.ListRoles)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("roles:create", middlewareOptions), c.CreateRole)
		g.PUT("/:roleId", organizationhttpmiddlewares.UserMustHavePermission("roles:update", middlewareOptions), c.UpdateRole)
		g.DELETE("/:roleId", organizationhttpmiddlewares.UserMustHavePermission("roles:delete", middlewareOptions), c.DeleteRole)
	}

	return g
}
