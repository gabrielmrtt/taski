package workspacehttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacehttpmiddlewares "github.com/gabrielmrtt/taski/internal/workspace/infra/http/middlewares"
	workspacehttprequests "github.com/gabrielmrtt/taski/internal/workspace/infra/http/requests"
	workspaceservice "github.com/gabrielmrtt/taski/internal/workspace/service"
	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	ListWorkspacesService  *workspaceservice.ListWorkspacesService
	GetWorkspaceService    *workspaceservice.GetWorkspaceService
	CreateWorkspaceService *workspaceservice.CreateWorkspaceService
	UpdateWorkspaceService *workspaceservice.UpdateWorkspaceService
	DeleteWorkspaceService *workspaceservice.DeleteWorkspaceService
}

func NewWorkspaceHandler(
	listWorkspacesService *workspaceservice.ListWorkspacesService,
	getWorkspaceService *workspaceservice.GetWorkspaceService,
	createWorkspaceService *workspaceservice.CreateWorkspaceService,
	updateWorkspaceService *workspaceservice.UpdateWorkspaceService,
	deleteWorkspaceService *workspaceservice.DeleteWorkspaceService,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		ListWorkspacesService:  listWorkspacesService,
		GetWorkspaceService:    getWorkspaceService,
		CreateWorkspaceService: createWorkspaceService,
		UpdateWorkspaceService: updateWorkspaceService,
		DeleteWorkspaceService: deleteWorkspaceService,
	}
}

type ListWorkspacesResponse = corehttp.HttpSuccessResponseWithData[workspace.WorkspaceDto]

// ListWorkspaces godoc
// @Summary List workspaces in an organization
// @Description Returns all workspaces in an organization based in the authenticated user accesses.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request query workspacehttprequests.ListWorkspacesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListWorkspacesResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace [get]
func (c *WorkspaceHandler) ListWorkspaces(ctx *gin.Context) {
	var request workspacehttprequests.ListWorkspacesRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input workspaceservice.ListWorkspacesInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity
	input.Filters.AuthenticatedUserIdentity = authenticatedUserIdentity

	response, err := c.ListWorkspacesService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetWorkspaceResponse = corehttp.HttpSuccessResponseWithData[workspace.WorkspaceDto]

// GetWorkspace godoc
// @Summary Get a workspace in an organization
// @Description Returns a workspace in an organization based in the authenticated user accesses.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Param request query workspacehttprequests.GetWorkspaceRequest true "Query parameters"
// @Produce json
// @Success 200 {object} GetWorkspaceResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId [get]
func (c *WorkspaceHandler) GetWorkspace(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var workspaceIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	var request workspacehttprequests.GetWorkspaceRequest
	var input workspaceservice.GetWorkspaceInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.WorkspaceIdentity = workspaceIdentity

	response, err := c.GetWorkspaceService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateWorkspaceResponse = corehttp.HttpSuccessResponseWithData[workspace.WorkspaceDto]

// CreateWorkspace godoc
// @Summary Create a workspace in an organization
// @Description Creates a new workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param request body workspacehttprequests.CreateWorkspaceRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateWorkspaceResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace [post]
func (c *WorkspaceHandler) CreateWorkspace(ctx *gin.Context) {
	var request workspacehttprequests.CreateWorkspaceRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input workspaceservice.CreateWorkspaceInput = workspaceservice.CreateWorkspaceInput{
		OrganizationIdentity: *organizationIdentity,
		UserCreatorIdentity:  *authenticatedUserIdentity,
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

	response, err := c.CreateWorkspaceService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateWorkspaceResponse = corehttp.EmptyHttpSuccessResponse

// UpdateWorkspace godoc
// @Summary Update a workspace in an organization
// @Description Updates an existing workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Param request body workspacehttprequests.UpdateWorkspaceRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateWorkspaceResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId [put]
func (c *WorkspaceHandler) UpdateWorkspace(ctx *gin.Context) {
	var request workspacehttprequests.UpdateWorkspaceRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var workspaceIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input workspaceservice.UpdateWorkspaceInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.WorkspaceIdentity = workspaceIdentity
	input.UserEditorIdentity = *authenticatedUserIdentity

	err := c.UpdateWorkspaceService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteWorkspaceResponse = corehttp.EmptyHttpSuccessResponse

// DeleteWorkspace godoc
// @Summary Delete a workspace in an organization
// @Description Deletes an existing workspace in an organization.
// @Tags Workspace
// @Accept json
// @Param organizationId path string true "Organization ID"
// @Param workspaceId path string true "Workspace ID"
// @Produce json
// @Success 200 {object} DeleteWorkspaceResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId [delete]
func (c *WorkspaceHandler) DeleteWorkspace(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var workspaceIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	var input workspaceservice.DeleteWorkspaceInput = workspaceservice.DeleteWorkspaceInput{
		OrganizationIdentity: *organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
	}

	err := c.DeleteWorkspaceService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *WorkspaceHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/workspace")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("workspaces:view", middlewareOptions), c.ListWorkspaces)
		g.GET("/:workspaceId", organizationhttpmiddlewares.UserMustHavePermission("workspaces:view", middlewareOptions), workspacehttpmiddlewares.UserMustBeInWorkspace(middlewareOptions), c.GetWorkspace)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("workspaces:create", middlewareOptions), c.CreateWorkspace)
		g.PUT("/:workspaceId", organizationhttpmiddlewares.UserMustHavePermission("workspaces:update", middlewareOptions), workspacehttpmiddlewares.UserMustBeInWorkspace(middlewareOptions), c.UpdateWorkspace)
		g.DELETE("/:workspaceId", organizationhttpmiddlewares.UserMustHavePermission("workspaces:delete", middlewareOptions), workspacehttpmiddlewares.UserMustBeInWorkspace(middlewareOptions), c.DeleteWorkspace)
	}

	return g
}
