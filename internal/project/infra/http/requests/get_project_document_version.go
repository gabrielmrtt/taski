package projecthttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetProjectDocumentVersionRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetProjectDocumentVersionRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetProjectDocumentVersionRequest) ToInput() projectservice.GetProjectDocumentVersionInput {
	return projectservice.GetProjectDocumentVersionInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
