package projecthttp

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	project "github.com/gabrielmrtt/taski/internal/project"
	projecthttpmiddlewares "github.com/gabrielmrtt/taski/internal/project/infra/http/middlewares"
	projecthttprequests "github.com/gabrielmrtt/taski/internal/project/infra/http/requests"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	userhttpmiddlewares "github.com/gabrielmrtt/taski/internal/user/infra/http/middlewares"
	workspacehttpmiddlewares "github.com/gabrielmrtt/taski/internal/workspace/infra/http/middlewares"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	ListProjectsService  *projectservice.ListProjectsService
	GetProjectService    *projectservice.GetProjectService
	CreateProjectService *projectservice.CreateProjectService
	UpdateProjectService *projectservice.UpdateProjectService
	DeleteProjectService *projectservice.DeleteProjectService
}

func NewProjectHandler(
	listProjectsService *projectservice.ListProjectsService,
	getProjectService *projectservice.GetProjectService,
	createProjectService *projectservice.CreateProjectService,
	updateProjectService *projectservice.UpdateProjectService,
	deleteProjectService *projectservice.DeleteProjectService,
) *ProjectHandler {
	return &ProjectHandler{
		ListProjectsService:  listProjectsService,
		GetProjectService:    getProjectService,
		CreateProjectService: createProjectService,
		UpdateProjectService: updateProjectService,
		DeleteProjectService: deleteProjectService,
	}
}

type ListProjectsResponse = corehttp.HttpSuccessResponseWithData[project.ProjectDto]

// ListProjects godoc
// @Summary List projects in a workspace
// @Description Returns all projects in a workspace accessible by the authenticated user.
// @Tags Project
// @Accept json
// @Param request query projecthttprequests.ListProjectsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListProjectsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId/project [get]
func (c *ProjectHandler) ListProjects(ctx *gin.Context) {
	var request projecthttprequests.ListProjectsRequest

	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.Filters.LoggedUserIdentity = &authenticatedUserIdentity

	response, err := c.ListProjectsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetProjectResponse = corehttp.HttpSuccessResponseWithData[project.ProjectDto]

// GetProject godoc
// @Summary Get a project in a workspace
// @Description Returns an accessible project in a workspace by its ID.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Produce json
// @Success 200 {object} GetProjectResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [get]
func (c *ProjectHandler) GetProject(ctx *gin.Context) {
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))

	input := projectservice.GetProjectInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
		ProjectIdentity:      projectIdentity,
	}

	response, err := c.GetProjectService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateProjectResponse = corehttp.HttpSuccessResponseWithData[project.ProjectDto]

// CreateProject godoc
// @Summary Create a project in a workspace
// @Description Creates a new project in a workspace.
// @Tags Project
// @Accept json
// @Param request body projecthttprequests.CreateProjectRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateProjectResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId/project [post]
func (c *ProjectHandler) CreateProject(ctx *gin.Context) {
	var request projecthttprequests.CreateProjectRequest
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateProjectService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateProjectResponse = corehttp.EmptyHttpSuccessResponse

// UpdateProject godoc
// @Summary Update a project in a workspace
// @Description Updates an accessible project in a workspace.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Param request body projecthttprequests.UpdateProjectRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateProjectResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [put]
func (c *ProjectHandler) UpdateProject(ctx *gin.Context) {
	var request projecthttprequests.UpdateProjectRequest
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))
	authenticatedUserIdentity := userhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.UserEditorIdentity = authenticatedUserIdentity
	input.ProjectIdentity = projectIdentity

	err := c.UpdateProjectService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteProjectResponse = corehttp.EmptyHttpSuccessResponse

// DeleteProject godoc
// @Summary Delete a project in a workspace
// @Description Deletes an accessible project in a workspace.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Produce json
// @Success 200 {object} DeleteProjectResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [delete]
func (c *ProjectHandler) DeleteProject(ctx *gin.Context) {
	organizationIdentity := organizationhttpmiddlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))

	input := projectservice.DeleteProjectInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
		ProjectIdentity:      projectIdentity,
	}

	err := c.DeleteProjectService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectHandler) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organizationId/workspace/:workspaceId/project")
	{
		g.Use(userhttpmiddlewares.AuthMiddleware())
		g.Use(workspacehttpmiddlewares.UserMustBeInWorkspace())

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("projects:view"), c.ListProjects)
		g.GET("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:view"), projecthttpmiddlewares.UserMustBeInProject(), c.GetProject)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("projects:create"), c.CreateProject)
		g.PUT("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:update"), projecthttpmiddlewares.UserMustBeInProject(), c.UpdateProject)
		g.DELETE("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:delete"), projecthttpmiddlewares.UserMustBeInProject(), c.DeleteProject)
	}

	return g
}
