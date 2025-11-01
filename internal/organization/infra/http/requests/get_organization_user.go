package organizationhttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetOrganizationUserRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetOrganizationUserRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetOrganizationUserRequest) ToInput() organizationservice.GetOrganizationUserInput {
	return organizationservice.GetOrganizationUserInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
