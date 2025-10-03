package userhttprequests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetMeRequest struct {
	Relations string `json:"relations" schema:"relations"`
}

func (r *GetMeRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetMeRequest) ToInput() userservice.GetMeInput {
	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != "" {
		relationsInput = strings.Split(stringutils.CamelCaseToPascalCase(r.Relations), ",")
	}
	return userservice.GetMeInput{
		RelationsInput: relationsInput,
	}
}
