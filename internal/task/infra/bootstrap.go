package taskinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	taskdatabase "github.com/gabrielmrtt/taski/internal/task/infra/database"
	taskhttp "github.com/gabrielmrtt/taski/internal/task/infra/http"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	taskRepository := taskdatabase.NewTaskBunRepository(options.DbConnection)
	taskCommentRepository := taskdatabase.NewTaskCommentBunRepository(options.DbConnection)
	projectRepository := projectdatabase.NewProjectBunRepository(options.DbConnection)
	projectTaskStatusRepository := projectdatabase.NewProjectTaskStatusBunRepository(options.DbConnection)
	projectTaskCategoryRepository := projectdatabase.NewProjectTaskCategoryBunRepository(options.DbConnection)
	projectUserRepository := projectdatabase.NewProjectUserBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)
	uploadedFileRepository := storagedatabase.NewUploadedFileBunRepository(options.DbConnection)
	storageRepository := storagedatabase.NewLocalStorageRepository()

	listTasksService := taskservice.NewListTasksService(taskRepository)
	getTaskService := taskservice.NewGetTaskService(taskRepository)
	createTaskService := taskservice.NewCreateTaskService(taskRepository, projectRepository, projectUserRepository, projectTaskStatusRepository, projectTaskCategoryRepository, transactionRepository)
	updateTaskService := taskservice.NewUpdateTaskService(taskRepository, projectTaskCategoryRepository, projectUserRepository, transactionRepository)
	deleteTaskService := taskservice.NewDeleteTaskService(taskRepository, transactionRepository)
	addSubTaskService := taskservice.NewAddSubTaskService(taskRepository, transactionRepository)
	updateSubTaskService := taskservice.NewUpdateSubTaskService(taskRepository, transactionRepository)
	removeSubTaskService := taskservice.NewRemoveSubTaskService(taskRepository, transactionRepository)
	changeTaskStatusService := taskservice.NewChangeTaskStatusService(taskRepository, projectTaskStatusRepository, transactionRepository)

	listTaskCommentsService := taskservice.NewListTaskCommentsService(taskCommentRepository, taskRepository)
	createTaskCommentService := taskservice.NewCreateTaskCommentService(taskRepository, taskCommentRepository, uploadedFileRepository, storageRepository, projectUserRepository, transactionRepository)
	updateTaskCommentService := taskservice.NewUpdateTaskCommentService(taskCommentRepository, taskRepository, uploadedFileRepository, storageRepository, transactionRepository)
	deleteTaskCommentService := taskservice.NewDeleteTaskCommentService(taskCommentRepository, taskRepository, uploadedFileRepository, storageRepository, transactionRepository)

	taskHandler := taskhttp.NewTaskHandler(listTasksService, getTaskService, createTaskService, updateTaskService, deleteTaskService, addSubTaskService, updateSubTaskService, removeSubTaskService, changeTaskStatusService)
	taskCommentHandler := taskhttp.NewTaskCommentHandler(listTaskCommentsService, createTaskCommentService, updateTaskCommentService, deleteTaskCommentService)

	taskHandler.ConfigureRoutes(corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	})

	taskCommentHandler.ConfigureRoutes(corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	})
}
