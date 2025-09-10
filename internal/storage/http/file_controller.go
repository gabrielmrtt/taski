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

type GetFileContentByIdentityResponse = []byte

// GetFileContent godoc
// @Summary Get file content
// @Description Returns the file contents.
// @Tags File
// @Accept json
// @Param file_id path string true "File ID"
// @Produce json
// @Success 200 {object} GetFileContentByIdentityResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /file/:file_id [get]
func (c *FileController) GetFileContent(ctx *gin.Context) {
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
		g.GET("/:file_id", c.GetFileContent)
	}
}
