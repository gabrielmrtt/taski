package projecthttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/project"
	projecthttpmiddlewares "github.com/gabrielmrtt/taski/internal/project/infra/http/middlewares"
	projecthttprequests "github.com/gabrielmrtt/taski/internal/project/infra/http/requests"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
)

type ProjectTaskStatusHandler struct {
	ListProjectTaskStatusesService *projectservice.ListProjectTaskStatusesService
	CreateProjectTaskStatusService *projectservice.CreateProjectTaskStatusService
	UpdateProjectTaskStatusService *projectservice.UpdateProjectTaskStatusService
	DeleteProjectTaskStatusService *projectservice.DeleteProjectTaskStatusService
}

func NewProjectTaskStatusHandler(
	listProjectTaskStatusesService *projectservice.ListProjectTaskStatusesService,
	createProjectTaskStatusService *projectservice.CreateProjectTaskStatusService,
	updateProjectTaskStatusService *projectservice.UpdateProjectTaskStatusService,
	deleteProjectTaskStatusService *projectservice.DeleteProjectTaskStatusService,
) *ProjectTaskStatusHandler {
	return &ProjectTaskStatusHandler{
		ListProjectTaskStatusesService: listProjectTaskStatusesService,
		CreateProjectTaskStatusService: createProjectTaskStatusService,
		UpdateProjectTaskStatusService: updateProjectTaskStatusService,
		DeleteProjectTaskStatusService: deleteProjectTaskStatusService,
	}
}

type ListProjectTaskStatusesResponse = corehttp.HttpSuccessResponseWithData[project.ProjectTaskStatusDto]

// ListProjectTaskStatuses godoc
// @Summary List project task statuses
// @Description Returns all project task statuses.
// @Tags Project Task Status
// @Accept json
// @Param request query projecthttprequests.ListProjectTaskStatusesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListProjectTaskStatusesResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-status [get]
func (c *ProjectTaskStatusHandler) ListProjectTaskStatuses(ctx *gin.Context) {
	var request projecthttprequests.ListProjectTaskStatusesRequest
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.ListProjectTaskStatusesInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.ProjectIdentity = &projectIdentity

	response, err := c.ListProjectTaskStatusesService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateProjectTaskStatusResponse = corehttp.HttpSuccessResponseWithData[project.ProjectTaskStatusDto]

// CreateProjectTaskStatus godoc
// @Summary Create a project task status
// @Description Creates a new project task status.
// @Tags Project Task Status
// @Accept json
// @Param request body projecthttprequests.CreateProjectTaskStatusRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateProjectTaskStatusResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-status [post]
func (c *ProjectTaskStatusHandler) CreateProjectTaskStatus(ctx *gin.Context) {
	var request projecthttprequests.CreateProjectTaskStatusRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.CreateProjectTaskStatusInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.ProjectIdentity = projectIdentity

	response, err := c.CreateProjectTaskStatusService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateProjectTaskStatusResponse = corehttp.EmptyHttpSuccessResponse

// UpdateProjectTaskStatus godoc
// @Summary Update a project task status
// @Description Updates an existing project task status.
// @Tags Project Task Status
// @Accept json
// @Param projectId path string true "Project ID"
// @Param taskStatusId path string true "Task Status ID"
// @Param request body projecthttprequests.UpdateProjectTaskStatusRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateProjectTaskStatusResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-status/:taskStatusId [put]
func (c *ProjectTaskStatusHandler) UpdateProjectTaskStatus(ctx *gin.Context) {
	var request projecthttprequests.UpdateProjectTaskStatusRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var taskStatusIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("taskStatusId"))
	var input projectservice.UpdateProjectTaskStatusInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.ProjectIdentity = projectIdentity
	input.ProjectTaskStatusIdentity = taskStatusIdentity

	err := c.UpdateProjectTaskStatusService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteProjectTaskStatusResponse = corehttp.EmptyHttpSuccessResponse

// DeleteProjectTaskStatus godoc
// @Summary Delete a project task status
// @Description Deletes an existing project task status.
// @Tags Project Task Status
// @Accept json
// @Param projectId path string true "Project ID"
// @Param taskStatusId path string true "Task Status ID"
// @Produce json
// @Success 200 {object} DeleteProjectTaskStatusResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-status/:taskStatusId [delete]
func (c *ProjectTaskStatusHandler) DeleteProjectTaskStatus(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var taskStatusIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("taskStatusId"))
	var input projectservice.DeleteProjectTaskStatusInput = projectservice.DeleteProjectTaskStatusInput{
		OrganizationIdentity:      *organizationIdentity,
		ProjectIdentity:           projectIdentity,
		ProjectTaskStatusIdentity: taskStatusIdentity,
	}

	err := c.DeleteProjectTaskStatusService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectTaskStatusHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/project/:projectId/task-status")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))
		g.Use(projecthttpmiddlewares.UserMustBeInProject(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.ListProjectTaskStatuses)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.CreateProjectTaskStatus)
		g.PUT("/:taskStatusId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.UpdateProjectTaskStatus)
		g.DELETE("/:taskStatusId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.DeleteProjectTaskStatus)
	}

	return g
}
