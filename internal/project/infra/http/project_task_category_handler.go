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

type ProjectTaskCategoryHandler struct {
	ListProjectTaskCategoriesService *projectservice.ListProjectTaskCategoriesService
	CreateProjectTaskCategoryService *projectservice.CreateProjectTaskCategoryService
	UpdateProjectTaskCategoryService *projectservice.UpdateProjectTaskCategoryService
	DeleteProjectTaskCategoryService *projectservice.DeleteProjectTaskCategoryService
}

func NewProjectTaskCategoryHandler(
	listProjectTaskCategoriesService *projectservice.ListProjectTaskCategoriesService,
	createProjectTaskCategoryService *projectservice.CreateProjectTaskCategoryService,
	updateProjectTaskCategoryService *projectservice.UpdateProjectTaskCategoryService,
	deleteProjectTaskCategoryService *projectservice.DeleteProjectTaskCategoryService,
) *ProjectTaskCategoryHandler {
	return &ProjectTaskCategoryHandler{
		ListProjectTaskCategoriesService: listProjectTaskCategoriesService,
		CreateProjectTaskCategoryService: createProjectTaskCategoryService,
		UpdateProjectTaskCategoryService: updateProjectTaskCategoryService,
		DeleteProjectTaskCategoryService: deleteProjectTaskCategoryService,
	}
}

type ListProjectTaskCategoriesResponse = corehttp.HttpSuccessResponseWithData[project.ProjectTaskCategoryDto]

// ListProjectTaskCategories godoc
// @Summary List project task categories
// @Description Returns all project task categories.
// @Tags Project Task Category
// @Accept json
// @Param request query projecthttprequests.ListProjectTaskCategoriesRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListProjectTaskCategoriesResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-category [get]
func (c *ProjectTaskCategoryHandler) ListProjectTaskCategories(ctx *gin.Context) {
	var request projecthttprequests.ListProjectTaskCategoriesRequest
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.ListProjectTaskCategoriesInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.ProjectIdentity = &projectIdentity

	response, err := c.ListProjectTaskCategoriesService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateProjectTaskCategoryResponse = corehttp.HttpSuccessResponseWithData[project.ProjectTaskCategoryDto]

// CreateProjectTaskCategory godoc
// @Summary Create a project task category
// @Description Creates a new project task category.
// @Tags Project Task Category
// @Accept json
// @Param request body projecthttprequests.CreateProjectTaskCategoryRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateProjectTaskCategoryResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-category [post]
func (c *ProjectTaskCategoryHandler) CreateProjectTaskCategory(ctx *gin.Context) {
	var request projecthttprequests.CreateProjectTaskCategoryRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.CreateProjectTaskCategoryInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.ProjectIdentity = projectIdentity

	response, err := c.CreateProjectTaskCategoryService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateProjectTaskCategoryResponse = corehttp.EmptyHttpSuccessResponse

// UpdateProjectTaskCategory godoc
// @Summary Update a project task category
// @Description Updates an existing project task category.
// @Tags Project Task Category
// @Accept json
// @Param projectId path string true "Project ID"
// @Param taskCategoryId path string true "Task Category ID"
// @Param request body projecthttprequests.UpdateProjectTaskCategoryRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateProjectTaskCategoryResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-category/:taskCategoryId [put]
func (c *ProjectTaskCategoryHandler) UpdateProjectTaskCategory(ctx *gin.Context) {
	var request projecthttprequests.UpdateProjectTaskCategoryRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var taskCategoryIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("taskCategoryId"))
	var input projectservice.UpdateProjectTaskCategoryInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.ProjectIdentity = projectIdentity
	input.ProjectTaskCategoryIdentity = taskCategoryIdentity

	err := c.UpdateProjectTaskCategoryService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteProjectTaskCategoryResponse = corehttp.EmptyHttpSuccessResponse

// DeleteProjectTaskCategory godoc
// @Summary Delete a project task category
// @Description Deletes an existing project task category.
// @Tags Project Task Category
// @Accept json
// @Param projectId path string true "Project ID"
// @Param taskCategoryId path string true "Task Category ID"
// @Produce json
// @Success 200 {object} DeleteProjectTaskCategoryResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/task-category/:taskCategoryId [delete]
func (c *ProjectTaskCategoryHandler) DeleteProjectTaskCategory(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var taskCategoryIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("taskCategoryId"))
	var input projectservice.DeleteProjectTaskCategoryInput = projectservice.DeleteProjectTaskCategoryInput{
		OrganizationIdentity:        *organizationIdentity,
		ProjectIdentity:             projectIdentity,
		ProjectTaskCategoryIdentity: taskCategoryIdentity,
	}

	err := c.DeleteProjectTaskCategoryService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectTaskCategoryHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/project/:projectId/task-category")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))
		g.Use(projecthttpmiddlewares.UserMustBeInProject(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.ListProjectTaskCategories)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.CreateProjectTaskCategory)
		g.PUT("/:taskCategoryId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.UpdateProjectTaskCategory)
		g.DELETE("/:taskCategoryId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.DeleteProjectTaskCategory)
	}

	return g
}
