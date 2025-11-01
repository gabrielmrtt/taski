package workspacehttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	workspaceservice "github.com/gabrielmrtt/taski/internal/workspace/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetWorkspaceRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetWorkspaceRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetWorkspaceRequest) ToInput() workspaceservice.GetWorkspaceInput {
	return workspaceservice.GetWorkspaceInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
