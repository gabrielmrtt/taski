package storagehttp

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	GetFileContentByIdentityService *storageservice.GetFileContentByIdentityService
}

func NewFileHandler(
	getFileContentByIdentityService *storageservice.GetFileContentByIdentityService,
) *FileHandler {
	return &FileHandler{
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
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /file/:file_id [get]
func (c *FileHandler) GetFileContent(ctx *gin.Context) {
	fileId := ctx.Param("file_id")

	input := storageservice.GetFileContentByIdentityInput{
		FileIdentity: core.NewIdentityFromPublic(fileId),
	}

	fileContent, err := c.GetFileContentByIdentityService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	ctx.Data(http.StatusOK, fileContent.FileMimeType, fileContent.FileContent)
}

func (c *FileHandler) ConfigureRoutes(group *gin.RouterGroup) {
	g := group.Group("/file")
	{
		g.GET("/:file_id", c.GetFileContent)
	}
}
