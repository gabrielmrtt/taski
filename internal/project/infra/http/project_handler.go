package projecthttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	project "github.com/gabrielmrtt/taski/internal/project"
	projecthttpmiddlewares "github.com/gabrielmrtt/taski/internal/project/infra/http/middlewares"
	projecthttprequests "github.com/gabrielmrtt/taski/internal/project/infra/http/requests"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
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
// @Router /project [get]
func (c *ProjectHandler) ListProjects(ctx *gin.Context) {
	var request projecthttprequests.ListProjectsRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input projectservice.ListProjectsInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity
	input.Filters.AuthenticatedUserIdentity = authenticatedUserIdentity

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
// @Router /project/:projectId [get]
func (c *ProjectHandler) GetProject(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.GetProjectInput = projectservice.GetProjectInput{
		OrganizationIdentity: *organizationIdentity,
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
// @Router /project [post]
func (c *ProjectHandler) CreateProject(ctx *gin.Context) {
	var request projecthttprequests.CreateProjectRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input projectservice.CreateProjectInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

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
// @Router /project/:projectId [put]
func (c *ProjectHandler) UpdateProject(ctx *gin.Context) {
	var request projecthttprequests.UpdateProjectRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input projectservice.UpdateProjectInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserEditorIdentity = *authenticatedUserIdentity
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
// @Router /project/:projectId [delete]
func (c *ProjectHandler) DeleteProject(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.DeleteProjectInput = projectservice.DeleteProjectInput{
		OrganizationIdentity: *organizationIdentity,
		ProjectIdentity:      projectIdentity,
	}

	err := c.DeleteProjectService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/project")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.ListProjects)
		g.GET("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), projecthttpmiddlewares.UserMustBeInProject(middlewareOptions), c.GetProject)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("projects:create", middlewareOptions), c.CreateProject)
		g.PUT("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), projecthttpmiddlewares.UserMustBeInProject(middlewareOptions), c.UpdateProject)
		g.DELETE("/:projectId", organizationhttpmiddlewares.UserMustHavePermission("projects:delete", middlewareOptions), projecthttpmiddlewares.UserMustBeInProject(middlewareOptions), c.DeleteProject)
	}

	return g
}
