package teamhttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetTeamRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetTeamRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetTeamRequest) ToInput() teamservice.GetTeamInput {
	return teamservice.GetTeamInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
