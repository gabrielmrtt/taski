package workspace_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_http_middlewares "github.com/gabrielmrtt/taski/internal/workspace/http/middlewares"
	workspace_http_requests "github.com/gabrielmrtt/taski/internal/workspace/http/requests"
	workspace_services "github.com/gabrielmrtt/taski/internal/workspace/services"
	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	ListWorkspacesService  *workspace_services.ListWorkspacesService
	GetWorkspaceService    *workspace_services.GetWorkspaceService
	CreateWorkspaceService *workspace_services.CreateWorkspaceService
	UpdateWorkspaceService *workspace_services.UpdateWorkspaceService
	DeleteWorkspaceService *workspace_services.DeleteWorkspaceService
}

func NewWorkspaceController(
	listWorkspacesService *workspace_services.ListWorkspacesService,
	getWorkspaceService *workspace_services.GetWorkspaceService,
	createWorkspaceService *workspace_services.CreateWorkspaceService,
	updateWorkspaceService *workspace_services.UpdateWorkspaceService,
	deleteWorkspaceService *workspace_services.DeleteWorkspaceService,
) *WorkspaceController {
	return &WorkspaceController{
		ListWorkspacesService:  listWorkspacesService,
		GetWorkspaceService:    getWorkspaceService,
		CreateWorkspaceService: createWorkspaceService,
		UpdateWorkspaceService: updateWorkspaceService,
		DeleteWorkspaceService: deleteWorkspaceService,
	}
}

type ListWorkspacesResponse = core_http.HttpSuccessResponseWithData[workspace_core.WorkspaceDto]

// ListWorkspaces godoc
// @Summary List workspaces in an organization
// @Description Returns all workspaces in an organization based in the authenticated user accesses.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request query workspace_http_requests.ListWorkspacesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListWorkspacesResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/workspace [get]
func (c *WorkspaceController) ListWorkspaces(ctx *gin.Context) {
	var request workspace_http_requests.ListWorkspacesRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	response, err := c.ListWorkspacesService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

type GetWorkspaceResponse = core_http.HttpSuccessResponseWithData[workspace_core.WorkspaceDto]

// GetWorkspace godoc
// @Summary Get a workspace in an organization
// @Description Returns a workspace in an organization based in the authenticated user accesses.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Produce json
// @Success 200 {object} GetWorkspaceResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/workspace/:workspaceId [get]
func (c *WorkspaceController) GetWorkspace(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))

	input := workspace_services.GetWorkspaceInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
	}

	response, err := c.GetWorkspaceService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

type CreateWorkspaceResponse = core_http.HttpSuccessResponseWithData[workspace_core.WorkspaceDto]

// CreateWorkspace godoc
// @Summary Create a workspace in an organization
// @Description Creates a new workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request body workspace_http_requests.CreateWorkspaceRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateWorkspaceResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/workspace [post]
func (c *WorkspaceController) CreateWorkspace(ctx *gin.Context) {
	var request workspace_http_requests.CreateWorkspaceRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	response, err := c.CreateWorkspaceService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
	return
}

type UpdateWorkspaceResponse = core_http.EmptyHttpSuccessResponse

// UpdateWorkspace godoc
// @Summary Update a workspace in an organization
// @Description Updates an existing workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Param request body workspace_http_requests.UpdateWorkspaceRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateWorkspaceResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/workspace/:workspaceId [put]
func (c *WorkspaceController) UpdateWorkspace(ctx *gin.Context) {
	var request workspace_http_requests.UpdateWorkspaceRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.WorkspaceIdentity = workspaceIdentity

	err := c.UpdateWorkspaceService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

type DeleteWorkspaceResponse = core_http.EmptyHttpSuccessResponse

// DeleteWorkspace godoc
// @Summary Delete a workspace in an organization
// @Description Deletes an existing workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Produce json
// @Success 200 {object} DeleteWorkspaceResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/workspace/:workspaceId [delete]
func (c *WorkspaceController) DeleteWorkspace(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))

	input := workspace_services.DeleteWorkspaceInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
	}

	err := c.DeleteWorkspaceService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
	return
}

func (c *WorkspaceController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organizationId/workspace")
	{
		g.Use(user_http_middlewares.AuthMiddleware())

		g.GET("", organization_http_middlewares.UserMustHavePermission("workspaces:view"), c.ListWorkspaces)
		g.GET("/:workspaceId", organization_http_middlewares.UserMustHavePermission("workspaces:view"), workspace_http_middlewares.UserMustBeInWorkspace(), c.GetWorkspace)
		g.POST("", organization_http_middlewares.UserMustHavePermission("workspaces:create"), c.CreateWorkspace)
		g.PUT("/:workspaceId", organization_http_middlewares.UserMustHavePermission("workspaces:update"), workspace_http_middlewares.UserMustBeInWorkspace(), c.UpdateWorkspace)
		g.DELETE("/:workspaceId", organization_http_middlewares.UserMustHavePermission("workspaces:delete"), workspace_http_middlewares.UserMustBeInWorkspace(), c.DeleteWorkspace)
	}

	return g
}
