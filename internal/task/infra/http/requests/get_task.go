package taskhttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetTaskRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetTaskRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetTaskRequest) ToInput() taskservice.GetTaskInput {
	return taskservice.GetTaskInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
