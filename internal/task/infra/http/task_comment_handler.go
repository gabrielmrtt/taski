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

type TaskCommentHandler struct {
	ListTaskCommentsService  *taskservice.ListTaskCommentsService
	CreateTaskCommentService *taskservice.CreateTaskCommentService
	UpdateTaskCommentService *taskservice.UpdateTaskCommentService
	DeleteTaskCommentService *taskservice.DeleteTaskCommentService
}

func NewTaskCommentHandler(
	listTaskCommentsService *taskservice.ListTaskCommentsService,
	createTaskCommentService *taskservice.CreateTaskCommentService,
	updateTaskCommentService *taskservice.UpdateTaskCommentService,
	deleteTaskCommentService *taskservice.DeleteTaskCommentService,
) *TaskCommentHandler {
	return &TaskCommentHandler{
		ListTaskCommentsService:  listTaskCommentsService,
		CreateTaskCommentService: createTaskCommentService,
		UpdateTaskCommentService: updateTaskCommentService,
		DeleteTaskCommentService: deleteTaskCommentService,
	}
}

type ListTaskCommentsResponse = corehttp.HttpSuccessResponseWithData[core.PaginationOutput[task.TaskCommentDto]]

// ListTaskComments godoc
// @Summary List task comments
// @Description Returns all accessible task comments by the authenticated user.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param request body taskhttprequests.ListTaskCommentsRequest true "Request body"
// @Produce json
// @Success 200 {object} ListTaskCommentsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/comment [get]
func (h *TaskCommentHandler) ListTaskComments(c *gin.Context) {
	var request taskhttprequests.ListTaskCommentsRequest
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.ListTaskCommentsInput

	if err := request.FromQuery(c); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.TaskIdentity = taskIdentity
	response1, err := h.ListTaskCommentsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(c, http.StatusOK, response1)
}

type CreateTaskCommentResponse = corehttp.HttpSuccessResponseWithData[task.TaskCommentDto]

// CreateTaskComment godoc
// @Summary Create a task comment
// @Description Creates a new task comment.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param request body taskhttprequests.CreateTaskCommentRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateTaskCommentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/comment [post]
func (h *TaskCommentHandler) CreateTaskComment(c *gin.Context) {
	var request taskhttprequests.CreateTaskCommentRequest
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var input taskservice.CreateTaskCommentInput

	if err := c.ShouldBind(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.TaskIdentity = taskIdentity
	response, err := h.CreateTaskCommentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(c, http.StatusOK, response)

}

type UpdateTaskCommentResponse = corehttp.HttpSuccessResponseWithData[task.TaskCommentDto]

// UpdateTaskComment godoc
// @Summary Update a task comment
// @Description Updates an accessible task comment.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param commentId path string true "Comment ID"
// @Param request body taskhttprequests.UpdateTaskCommentRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateTaskCommentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/comment/:commentId [put]
func (h *TaskCommentHandler) UpdateTaskComment(c *gin.Context) {
	var request taskhttprequests.UpdateTaskCommentRequest
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var commentIdentity core.Identity = core.NewIdentityFromPublic(c.Param("commentId"))
	var input taskservice.UpdateTaskCommentInput

	if err := c.ShouldBind(&request); err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	input = request.ToInput()
	input.TaskIdentity = taskIdentity
	input.TaskCommentIdentity = commentIdentity
	err := h.UpdateTaskCommentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

type DeleteTaskCommentResponse = corehttp.EmptyHttpSuccessResponse

// DeleteTaskComment godoc
// @Summary Delete a task comment
// @Description Deletes an accessible task comment.
// @Tags Task
// @Accept json
// @Param taskId path string true "Task ID"
// @Param commentId path string true "Comment ID"
// @Produce json
// @Success 200 {object} DeleteTaskCommentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /task/:taskId/comment/:commentId [delete]
func (h *TaskCommentHandler) DeleteTaskComment(c *gin.Context) {
	var taskIdentity core.Identity = core.NewIdentityFromPublic(c.Param("taskId"))
	var commentIdentity core.Identity = core.NewIdentityFromPublic(c.Param("commentId"))
	var input taskservice.DeleteTaskCommentInput = taskservice.DeleteTaskCommentInput{
		TaskIdentity:        taskIdentity,
		TaskCommentIdentity: commentIdentity,
	}

	err := h.DeleteTaskCommentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(c, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(c, http.StatusOK)
}

func (h *TaskCommentHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/task/:taskId/comment")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))
		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.ListTaskComments)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.CreateTaskComment)
		g.PUT("/:commentId", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.UpdateTaskComment)
		g.DELETE("/:commentId", organizationhttpmiddlewares.UserMustHavePermission("tasks:view", middlewareOptions), h.DeleteTaskComment)
	}

	return g
}
