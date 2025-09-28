package project_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_http_middlewares "github.com/gabrielmrtt/taski/internal/project/http/middlewares"
	project_http_requests "github.com/gabrielmrtt/taski/internal/project/http/requests"
	project_services "github.com/gabrielmrtt/taski/internal/project/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	workspace_http_middlewares "github.com/gabrielmrtt/taski/internal/workspace/http/middlewares"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	ListProjectsService  *project_services.ListProjectsService
	GetProjectService    *project_services.GetProjectService
	CreateProjectService *project_services.CreateProjectService
	UpdateProjectService *project_services.UpdateProjectService
	DeleteProjectService *project_services.DeleteProjectService
}

func NewProjectController(
	listProjectsService *project_services.ListProjectsService,
	getProjectService *project_services.GetProjectService,
	createProjectService *project_services.CreateProjectService,
	updateProjectService *project_services.UpdateProjectService,
	deleteProjectService *project_services.DeleteProjectService,
) *ProjectController {
	return &ProjectController{
		ListProjectsService:  listProjectsService,
		GetProjectService:    getProjectService,
		CreateProjectService: createProjectService,
		UpdateProjectService: updateProjectService,
		DeleteProjectService: deleteProjectService,
	}
}

type ListProjectsResponse = core_http.HttpSuccessResponseWithData[project_core.ProjectDto]

// ListProjects godoc
// @Summary List projects in a workspace
// @Description Returns all projects in a workspace accessible by the authenticated user.
// @Tags Project
// @Accept json
// @Param request query project_http_requests.ListProjectsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListProjectsResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /workspace/:workspaceId/project [get]
func (c *ProjectController) ListProjects(ctx *gin.Context) {
	var request project_http_requests.ListProjectsRequest

	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.Filters.LoggedUserIdentity = &authenticatedUserIdentity

	response, err := c.ListProjectsService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetProjectResponse = core_http.HttpSuccessResponseWithData[project_core.ProjectDto]

// GetProject godoc
// @Summary Get a project in a workspace
// @Description Returns an accessible project in a workspace by its ID.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Produce json
// @Success 200 {object} GetProjectResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [get]
func (c *ProjectController) GetProject(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))

	input := project_services.GetProjectInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
		ProjectIdentity:      projectIdentity,
	}

	response, err := c.GetProjectService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateProjectResponse = core_http.HttpSuccessResponseWithData[project_core.ProjectDto]

// CreateProject godoc
// @Summary Create a project in a workspace
// @Description Creates a new project in a workspace.
// @Tags Project
// @Accept json
// @Param request body project_http_requests.CreateProjectRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateProjectResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /workspace/:workspaceId/project [post]
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var request project_http_requests.CreateProjectRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateProjectService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateProjectResponse = core_http.EmptyHttpSuccessResponse

// UpdateProject godoc
// @Summary Update a project in a workspace
// @Description Updates an accessible project in a workspace.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Param request body project_http_requests.UpdateProjectRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateProjectResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [put]
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	var request project_http_requests.UpdateProjectRequest
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.WorkspaceIdentity = workspaceIdentity
	input.OrganizationIdentity = organizationIdentity
	input.UserEditorIdentity = authenticatedUserIdentity
	input.ProjectIdentity = projectIdentity

	err := c.UpdateProjectService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteProjectResponse = core_http.EmptyHttpSuccessResponse

// DeleteProject godoc
// @Summary Delete a project in a workspace
// @Description Deletes an accessible project in a workspace.
// @Tags Project
// @Accept json
// @Param projectId path string true "Project ID"
// @Produce json
// @Success 200 {object} DeleteProjectResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /workspace/:workspaceId/project/:projectId [delete]
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	workspaceIdentity := core.NewIdentityFromPublic(ctx.Param("workspaceId"))
	projectIdentity := core.NewIdentityFromPublic(ctx.Param("projectId"))

	input := project_services.DeleteProjectInput{
		OrganizationIdentity: organizationIdentity,
		WorkspaceIdentity:    workspaceIdentity,
		ProjectIdentity:      projectIdentity,
	}

	err := c.DeleteProjectService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectController) ConfigureRoutes(group *gin.RouterGroup) *gin.RouterGroup {
	g := group.Group("/organization/:organizationId/workspace/:workspaceId/project")
	{
		g.Use(user_http_middlewares.AuthMiddleware())
		g.Use(workspace_http_middlewares.UserMustBeInWorkspace())

		g.GET("", organization_http_middlewares.UserMustHavePermission("projects:view"), c.ListProjects)
		g.GET("/:projectId", organization_http_middlewares.UserMustHavePermission("projects:view"), project_http_middlewares.UserMustBeInProject(), c.GetProject)
		g.POST("", organization_http_middlewares.UserMustHavePermission("projects:create"), c.CreateProject)
		g.PUT("/:projectId", organization_http_middlewares.UserMustHavePermission("projects:update"), project_http_middlewares.UserMustBeInProject(), c.UpdateProject)
		g.DELETE("/:projectId", organization_http_middlewares.UserMustHavePermission("projects:delete"), project_http_middlewares.UserMustBeInProject(), c.DeleteProject)
	}

	return g
}
