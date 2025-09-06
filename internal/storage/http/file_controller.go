package storage_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	storage_services "github.com/gabrielmrtt/taski/internal/storage/services"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	GetFileContentByIdentityService *storage_services.GetFileContentByIdentityService
}

func NewFileController(
	getFileContentByIdentityService *storage_services.GetFileContentByIdentityService,
) *FileController {
	return &FileController{
		GetFileContentByIdentityService: getFileContentByIdentityService,
	}
}

func (c *FileController) GetFileContentByIdentity(ctx *gin.Context) {
	fileId := ctx.Param("file_id")

	input := storage_services.GetFileContentByIdentityInput{
		FileIdentity: core.NewIdentityFromPublic(fileId),
	}

	fileContent, err := c.GetFileContentByIdentityService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	ctx.Data(http.StatusOK, fileContent.FileMimeType, fileContent.FileContent)
}

func (c *FileController) ConfigureRoutes(group *gin.RouterGroup) {
	g := group.Group("/file")
	{
		g.GET("/:file_id", c.GetFileContentByIdentity)
	}
}
