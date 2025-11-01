package userhttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetMeRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetMeRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetMeRequest) ToInput() userservice.GetMeInput {
	return userservice.GetMeInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
