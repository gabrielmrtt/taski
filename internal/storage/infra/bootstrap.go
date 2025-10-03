package storageinfra

import (
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	storagehttp "github.com/gabrielmrtt/taski/internal/storage/infra/http"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BoostrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BoostrapInfraOptions) {
	fileRepository := storagedatabase.NewUploadedFileBunRepository(options.DbConnection)
	storageRepository := storagedatabase.NewLocalStorageRepository()

	getFileContentByIdentityService := storageservice.NewGetFileContentByIdentityService(fileRepository, storageRepository)

	fileController := storagehttp.NewFileHandler(getFileContentByIdentityService)

	fileController.ConfigureRoutes(options.RouterGroup)
}
