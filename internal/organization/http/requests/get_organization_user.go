package organization_http_requests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetOrganizationUserRequest struct {
	Relations string `json:"relations" schema:"relations"`
}

func (r *GetOrganizationUserRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetOrganizationUserRequest) ToInput() organization_services.GetOrganizationUserInput {
	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != "" {
		relationsInput = strings.Split(r.Relations, ",")

		for i, relation := range relationsInput {
			relationsInput[i] = stringutils.CamelCaseToPascalCase(relation)
		}
	}

	return organization_services.GetOrganizationUserInput{
		RelationsInput: relationsInput,
	}
}
