package taskhttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/task"
	taskhttprequests "github.com/gabrielmrtt/taski/internal/task/infra/http/requests"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	ListTasksService        *taskservice.ListTasksService
	GetTaskService          *taskservice.GetTaskService
	CreateTaskService       *taskservice.CreateTaskService
	UpdateTaskService       *taskservice.UpdateTaskService
	DeleteTaskService       *taskservice.DeleteTaskService
	AddSubTaskService       *taskservice.AddSubTaskService
	UpdateSubTaskService    *taskservice.UpdateSubTaskService
	RemoveSubTaskService    *taskservice.RemoveSubTaskService
	ChangeTaskStatusService *taskservice.ChangeTaskStatusService
}

func NewTaskHandler(
	listTasksService *taskservice.ListTasksService,
	getTaskService *taskservice.GetTaskService,
	createTaskService *taskservice.CreateTaskService,
	updateTaskService *taskservice.UpdateTaskService,
	deleteTaskService *taskservice.DeleteTaskService,
	addSubTaskService *taskservice.AddSubTaskService,
	updateSubTaskService *taskservice.UpdateSubTaskService,
	removeSubTaskService *taskservice.RemoveSubTaskService,
	changeTaskStatusService *taskservice.ChangeTaskStatusService,
) *TaskHandler {
	return &TaskHandler{
		ListTasksService:        listTasksService,
		GetTaskService:          getTaskService,
		CreateTaskService:       createTaskService,
		UpdateTaskService:       updateTaskService,
		DeleteTaskService:       deleteTaskService,
		AddSubTaskService:       addSubTaskService,
		UpdateSubTaskService:    updateSubTaskService,
		RemoveSubTaskService:    removeSubTaskService,
		ChangeTaskStatusService: changeTaskStatusService,
	}
}

type ListTasksResponse = corehttp.HttpSuccessResponseWithData[core.PaginationOutput[task.TaskDto]]

// ListTasks godoc
// @Summary List tasks
// @Description Returns all accessible tasks by the authenticated user.
// @Tags Task
// @Accept json
// @Param request query taskhttprequests.ListTasksRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListTasksResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var request taskhttprequests.ListTasksRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(c)
	var input taskservice.ListTasksInput

	if err := request.FromQuery(c); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.Filters.OrganizationIdentity = organizationIdentity
	input.Filters.AuthenticatedUserIdentity = authenticatedUserIdentity
	result, err := h.ListTasksService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(c, http.StatusOK, result)
}

type GetTaskResponse = corehttp.HttpSuccessResponseWithData[task.TaskDto]

// GetTask godoc
// @Summary Get a task
// @Description Returns an accessible task by its ID.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Produce json
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	var request taskhttprequests.GetTaskRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.GetTaskInput

	if err := request.FromQuery(c); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.TaskIdentity = taskIdentity
	input.OrganizationIdentity = organizationIdentity
	result, err := h.GetTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(c, http.StatusOK, result)
}

type CreateTaskResponse = corehttp.HttpSuccessResponseWithData[task.TaskDto]

// CreateTask godoc
// @Summary Create a task
// @Description Creates a new task.
// @Tags Task
// @Accept json
// @Param request body taskhttprequests.CreateTaskRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var request taskhttprequests.CreateTaskRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(c)
	var input taskservice.CreateTaskInput

	if err := c.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

	response, err := h.CreateTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(c, http.StatusOK, response)
}

type UpdateTaskResponse = corehttp.EmptyHttpSuccessResponse

// UpdateTask godoc
// @Summary Update a task
// @Description Updates an accessible task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param request body taskhttprequests.UpdateTaskRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var request taskhttprequests.UpdateTaskRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.UpdateTaskInput

	if err := c.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.TaskIdentity = taskIdentity

	err := h.UpdateTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type DeleteTaskResponse = corehttp.EmptyHttpSuccessResponse

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes an accessible task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Produce json
// @Success 200 {object} DeleteTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.DeleteTaskInput = taskservice.DeleteTaskInput{
		OrganizationIdentity: organizationIdentity,
		TaskIdentity:         taskIdentity,
	}

	err := h.DeleteTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type AddSubTaskResponse = corehttp.EmptyHttpSuccessResponse

// AddSubTask godoc
// @Summary Add a sub task to a task
// @Description Adds a sub task to an accessible task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param request body taskhttprequests.AddSubTaskRequest true "Request body"
// @Produce json
// @Success 200 {object} AddSubTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/sub-task [post]
func (h *TaskHandler) AddSubTask(c *gin.Context) {
	var request taskhttprequests.AddSubTaskRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.AddSubTaskInput

	if err := c.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.TaskIdentity = taskIdentity

	err := h.AddSubTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type UpdateSubTaskResponse = corehttp.EmptyHttpSuccessResponse

// UpdateSubTask godoc
// @Summary Update a sub task
// @Description Updates an accessible sub task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param subTaskId path string true "Sub Task ID"
// @Param request body taskhttprequests.UpdateSubTaskRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateSubTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/sub-task/:subTaskId [put]
func (h *TaskHandler) UpdateSubTask(c *gin.Context) {
	var request taskhttprequests.UpdateSubTaskRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var subTaskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("subTaskId"))
	var input taskservice.UpdateSubTaskInput

	if err := c.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.TaskIdentity = taskIdentity
	input.SubTaskIdentity = subTaskIdentity

	err := h.UpdateSubTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type RemoveSubTaskResponse = corehttp.EmptyHttpSuccessResponse

// RemoveSubTask godoc
// @Summary Remove a sub task from a task
// @Description Removes a sub task from an accessible task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param subTaskId path string true "Sub Task ID"
// @Produce json
// @Success 200 {object} RemoveSubTaskResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/sub-task/:subTaskId [delete]
func (h *TaskHandler) RemoveSubTask(c *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var subTaskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("subTaskId"))
	var input taskservice.RemoveSubTaskInput = taskservice.RemoveSubTaskInput{
		OrganizationIdentity: organizationIdentity,
		TaskIdentity:         taskIdentity,
		SubTaskIdentity:      subTaskIdentity,
	}

	err := h.RemoveSubTaskService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type ChangeTaskStatusResponse = corehttp.EmptyHttpSuccessResponse

// ChangeTaskStatus godoc
// @Summary Change the status of a task
// @Description Changes the status of an accessible task.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param request body taskhttprequests.ChangeTaskStatusRequest true "Request body"
// @Produce json
// @Success 200 {object} ChangeTaskStatusResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/status [put]
func (h *TaskHandler) ChangeTaskStatus(c *gin.Context) {
	var request taskhttprequests.ChangeTaskStatusRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(c)
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.ChangeTaskStatusInput

	if err := c.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.TaskIdentity = taskIdentity

	err := h.ChangeTaskStatusService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}
}
func (h *TaskHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/task")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))
		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.ListTasks)
		g.GET("/:taskId", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.GetTask)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("tasks:create", middlewareOptions), h.CreateTask)
		g.PUT("/:taskId", organizationhttpmiddlewares.UserMustHavePermission("tasks:update", middlewareOptions), h.UpdateTask)
		g.DELETE("/:taskId", organizationhttpmiddlewares.UserMustHavePermission("tasks:delete", middlewareOptions), h.DeleteTask)
		g.PATCH("/:taskId/status", organizationhttpmiddlewares.UserMustHavePermission("tasks:update", middlewareOptions), h.ChangeTaskStatus)
		g.POST("/:taskId/sub-task", organizationhttpmiddlewares.UserMustHavePermission("tasks:create", middlewareOptions), h.AddSubTask)
		g.PUT("/:taskId/sub-task/:subTaskId", organizationhttpmiddlewares.UserMustHavePermission("tasks:update", middlewareOptions), h.UpdateSubTask)
		g.DELETE("/:taskId/sub-task/:subTaskId", organizationhttpmiddlewares.UserMustHavePermission("tasks:update", middlewareOptions), h.RemoveSubTask)
	}

	return g
}
