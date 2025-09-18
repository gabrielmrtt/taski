package role_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_http_requests "github.com/gabrielmrtt/taski/internal/role/http/requests"
	role_services "github.com/gabrielmrtt/taski/internal/role/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	CreateRoleService *role_services.CreateRoleService
	UpdateRoleService *role_services.UpdateRoleService
	DeleteRoleService *role_services.DeleteRoleService
	ListRolesService  *role_services.ListRolesService
}

func NewRoleController(
	createRoleService *role_services.CreateRoleService,
	updateRoleService *role_services.UpdateRoleService,
	deleteRoleService *role_services.DeleteRoleService,
	listRolesService *role_services.ListRolesService,
) *RoleController {
	return &RoleController{
		CreateRoleService: createRoleService,
		UpdateRoleService: updateRoleService,
		DeleteRoleService: deleteRoleService,
		ListRolesService:  listRolesService,
	}
}

type CreateRoleResponse = core_http.HttpSuccessResponseWithData[role_core.RoleDto]

// CreateRole godoc
// @Summary Create a role in an organization
// @Description Creates a new role in an organization.
// @Tags Role
// @Accept json
// @Param request body role_http_requests.CreateRoleRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateRoleResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/role [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var request role_http_requests.CreateRoleRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateRoleService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateRoleResponse = core_http.EmptyHttpSuccessResponse

// UpdateRole godoc
// @Summary Update a role in an organization
// @Description Updates an existing role in an organization.
// @Tags Role
// @Accept json
// @Produce json
// @Param roleId path string true "Role ID"
// @Param request body role_http_requests.UpdateRoleRequest true "Request body"
// @Success 200 {object} UpdateRoleResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/role/:roleId [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var request role_http_requests.UpdateRoleRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	roleIdentity := core.NewIdentityFromPublic(ctx.Param("roleId"))

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserEditorIdentity = authenticatedUserIdentity
	input.RoleIdentity = roleIdentity

	err := c.UpdateRoleService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteRoleResponse = core_http.EmptyHttpSuccessResponse

// DeleteRole godoc
// @Summary Delete a role in an organization
// @Description Deletes an existing role in an organization.
// @Tags Role
// @Accept json
// @Param roleId path string true "Role ID"
// @Produce json
// @Success 200 {object} DeleteRoleResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/role/:roleId [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	roleIdentity := core.NewIdentityFromPublic(ctx.Param("roleId"))

	input := role_services.DeleteRoleInput{
		RoleIdentity:         roleIdentity,
		OrganizationIdentity: organizationIdentity,
	}

	err := c.DeleteRoleService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type ListRolesResponse = core_http.HttpSuccessResponseWithData[role_core.RoleDto]

// ListRoles godoc
// @Summary List roles in an organization
// @Description Lists all existing roles in an organization.
// @Tags Role
// @Accept json
// @Param request query role_http_requests.ListRolesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListRolesResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/role [get]
func (c *RoleController) ListRoles(ctx *gin.Context) {
	var request role_http_requests.ListRolesRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	response, err := c.ListRolesService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *RoleController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organizationId/role")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", organization_http_middlewares.UserMustHavePermission("roles:view"), c.ListRoles)
		g.POST("", organization_http_middlewares.UserMustHavePermission("roles:create"), c.CreateRole)
		g.PUT("/:roleId", organization_http_middlewares.UserMustHavePermission("roles:update"), c.UpdateRole)
		g.DELETE("/:roleId", organization_http_middlewares.UserMustHavePermission("roles:delete"), c.DeleteRole)
	}

	return g
}
