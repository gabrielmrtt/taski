package storage_http

import (
	storage_database_local "github.com/gabrielmrtt/taski/internal/storage/database/local"
	storage_database_postgres "github.com/gabrielmrtt/taski/internal/storage/database/postgres"
	storage_services "github.com/gabrielmrtt/taski/internal/storage/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	fileRepository := storage_database_postgres.NewUploadedFilePostgresRepository()
	storageRepository := storage_database_local.NewLocalStorageRepository()

	getFileContentByIdentityService := storage_services.NewGetFileContentByIdentityService(fileRepository, storageRepository)

	fileController := NewFileController(getFileContentByIdentityService)

	fileController.ConfigureRoutes(g)
}
