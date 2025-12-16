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
	taskActionRepository := taskdatabase.NewTaskActionBunRepository(options.DbConnection)

	listTasksService := taskservice.NewListTasksService(taskRepository)
	getTaskService := taskservice.NewGetTaskService(taskRepository)
	createTaskService := taskservice.NewCreateTaskService(taskRepository, taskActionRepository, projectRepository, projectUserRepository, projectTaskStatusRepository, projectTaskCategoryRepository, transactionRepository)
	updateTaskService := taskservice.NewUpdateTaskService(taskRepository, taskActionRepository, projectTaskCategoryRepository, projectUserRepository, transactionRepository)
	deleteTaskService := taskservice.NewDeleteTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)
	addSubTaskService := taskservice.NewAddSubTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)
	updateSubTaskService := taskservice.NewUpdateSubTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)
	removeSubTaskService := taskservice.NewRemoveSubTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)
	changeTaskStatusService := taskservice.NewChangeTaskStatusService(taskRepository, projectTaskStatusRepository, taskActionRepository, projectUserRepository, transactionRepository)
	completeTaskService := taskservice.NewCompleteTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)
	completeSubTaskService := taskservice.NewCompleteSubTaskService(taskRepository, taskActionRepository, projectUserRepository, transactionRepository)

	listTaskCommentsService := taskservice.NewListTaskCommentsService(taskCommentRepository, taskRepository)
	createTaskCommentService := taskservice.NewCreateTaskCommentService(taskRepository, taskCommentRepository, uploadedFileRepository, storageRepository, projectUserRepository, taskActionRepository, transactionRepository)
	updateTaskCommentService := taskservice.NewUpdateTaskCommentService(taskCommentRepository, taskRepository, projectUserRepository, uploadedFileRepository, storageRepository, taskActionRepository, transactionRepository)
	deleteTaskCommentService := taskservice.NewDeleteTaskCommentService(taskCommentRepository, taskRepository, projectUserRepository, uploadedFileRepository, storageRepository, taskActionRepository, transactionRepository)

	taskHandler := taskhttp.NewTaskHandler(listTasksService, getTaskService, createTaskService, updateTaskService, deleteTaskService, addSubTaskService, updateSubTaskService, removeSubTaskService, changeTaskStatusService, completeTaskService, completeSubTaskService)
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
