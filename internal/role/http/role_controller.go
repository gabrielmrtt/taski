package role_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
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

func (c *RoleController) CreateRole(ctx *gin.Context) {
	var request role_http_requests.CreateRoleRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))

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

func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var request role_http_requests.UpdateRoleRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	roleIdentity := core.NewIdentityFromPublic(ctx.Param("role_id"))

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

func (c *RoleController) DeleteRole(ctx *gin.Context) {
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))
	roleIdentity := core.NewIdentityFromPublic(ctx.Param("role_id"))

	input := role_services.DeleteRoleInput{
		RoleIdentity:         roleIdentity,
		OrganizationIdentity: organizationIdentity,
		UserDeleterIdentity:  authenticatedUserIdentity,
	}

	err := c.DeleteRoleService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *RoleController) ListRoles(ctx *gin.Context) {
	var request role_http_requests.ListRolesRequest
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := core.NewIdentityFromPublic(ctx.Param("organization_id"))

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.LoggedUserIdentity = authenticatedUserIdentity

	response, err := c.ListRolesService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

func (c *RoleController) ConfigureRoutes(group *gin.RouterGroup) {
	orgGroup := group.Group("/organization/:organization_id")
	{
		orgGroup.Use(user_http_middlewares.AuthMiddleware())
		orgGroup.Use(organization_http_middlewares.BlockIfUserIsNotPartOfOrganization())
		g := orgGroup.Group("/role")
		{
			g.GET("", c.ListRoles)
			g.POST("", c.CreateRole)
			g.PUT("/:role_id", c.UpdateRole)
			g.DELETE("/:role_id", c.DeleteRole)
		}
	}
}
