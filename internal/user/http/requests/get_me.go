package user_http_requests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
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

func (r *GetMeRequest) ToInput() user_services.GetMeInput {
	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != "" {
		relationsInput = strings.Split(stringutils.CamelCaseToPascalCase(r.Relations), ",")
	}
	return user_services.GetMeInput{
		RelationsInput: relationsInput,
	}
}
