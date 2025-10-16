package projectinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	projecthttp "github.com/gabrielmrtt/taski/internal/project/infra/http"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	projectRepository := projectdatabase.NewProjectBunRepository(options.DbConnection)
	projectUserRepository := projectdatabase.NewProjectUserBunRepository(options.DbConnection)
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	workspaceRepository := workspacedatabase.NewWorkspaceBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)
	projectTaskCategoryRepository := projectdatabase.NewProjectTaskCategoryBunRepository(options.DbConnection)
	projectTaskStatusRepository := projectdatabase.NewProjectTaskStatusBunRepository(options.DbConnection)
	projectDocumentRepository := projectdatabase.NewProjectDocumentBunRepository(options.DbConnection)
	uploadedFileRepository := storagedatabase.NewUploadedFileBunRepository(options.DbConnection)
	storageRepository := storagedatabase.NewLocalStorageRepository()

	listProjectsService := projectservice.NewListProjectsService(projectRepository)
	getProjectService := projectservice.NewGetProjectService(projectRepository)
	createProjectService := projectservice.NewCreateProjectService(projectRepository, projectUserRepository, projectTaskStatusRepository, userRepository, workspaceRepository, transactionRepository)
	updateProjectService := projectservice.NewUpdateProjectService(projectRepository, transactionRepository)
	deleteProjectService := projectservice.NewDeleteProjectService(projectRepository, transactionRepository)

	listProjectTaskCategoriesService := projectservice.NewListProjectTaskCategoriesService(projectTaskCategoryRepository)
	createProjectTaskCategoryService := projectservice.NewCreateProjectTaskCategoryService(projectRepository, projectTaskCategoryRepository, transactionRepository)
	updateProjectTaskCategoryService := projectservice.NewUpdateProjectTaskCategoryService(projectRepository, projectTaskCategoryRepository, transactionRepository)
	deleteProjectTaskCategoryService := projectservice.NewDeleteProjectTaskCategoryService(projectRepository, projectTaskCategoryRepository, transactionRepository)

	listProjectTaskStatusesService := projectservice.NewListProjectTaskStatusesService(projectTaskStatusRepository)
	createProjectTaskStatusService := projectservice.NewCreateProjectTaskStatusService(projectRepository, projectTaskStatusRepository, transactionRepository)
	updateProjectTaskStatusService := projectservice.NewUpdateProjectTaskStatusService(projectRepository, projectTaskStatusRepository, transactionRepository)
	deleteProjectTaskStatusService := projectservice.NewDeleteProjectTaskStatusService(projectRepository, projectTaskStatusRepository, transactionRepository)

	listProjectDocumentsService := projectservice.NewListProjectDocumentsService(projectRepository, projectDocumentRepository)
	listProjectDocumentVersionsService := projectservice.NewListProjectDocumentVersionsService(projectRepository, projectDocumentRepository)
	getProjectDocumentVersionService := projectservice.NewGetProjectDocumentVersionService(projectRepository, projectDocumentRepository)
	createProjectDocumentService := projectservice.NewCreateProjectDocumentService(projectRepository, projectDocumentRepository, uploadedFileRepository, storageRepository, transactionRepository)
	updateProjectDocumentService := projectservice.NewUpdateProjectDocumentService(projectRepository, projectDocumentRepository, uploadedFileRepository, storageRepository, transactionRepository)
	deleteProjectDocumentService := projectservice.NewDeleteProjectDocumentService(projectRepository, projectDocumentRepository, uploadedFileRepository, storageRepository, transactionRepository)
	deleteProjectDocumentVersionService := projectservice.NewDeleteProjectDocumentVersionService(projectRepository, projectDocumentRepository, uploadedFileRepository, storageRepository, transactionRepository)

	configureRoutesOptions := corehttp.ConfigureRoutesOptions{
		DbConnection: options.DbConnection,
		RouterGroup:  options.RouterGroup,
	}

	projectController := projecthttp.NewProjectHandler(listProjectsService, getProjectService, createProjectService, updateProjectService, deleteProjectService)
	projectController.ConfigureRoutes(configureRoutesOptions)

	projectTaskCategoryController := projecthttp.NewProjectTaskCategoryHandler(listProjectTaskCategoriesService, createProjectTaskCategoryService, updateProjectTaskCategoryService, deleteProjectTaskCategoryService)
	projectTaskCategoryController.ConfigureRoutes(configureRoutesOptions)

	projectTaskStatusController := projecthttp.NewProjectTaskStatusHandler(listProjectTaskStatusesService, createProjectTaskStatusService, updateProjectTaskStatusService, deleteProjectTaskStatusService)
	projectTaskStatusController.ConfigureRoutes(configureRoutesOptions)

	projectDocumentController := projecthttp.NewProjectDocumentHandler(getProjectDocumentVersionService, listProjectDocumentsService, listProjectDocumentVersionsService, createProjectDocumentService, updateProjectDocumentService, deleteProjectDocumentService, deleteProjectDocumentVersionService)
	projectDocumentController.ConfigureRoutes(configureRoutesOptions)
}
