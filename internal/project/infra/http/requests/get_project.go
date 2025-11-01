package projecthttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetProjectRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetProjectRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetProjectRequest) ToInput() projectservice.GetProjectInput {
	return projectservice.GetProjectInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
