package organizationhttprequests

import (
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetOrganizationRequest struct {
	Relations *string `json:"relations" schema:"relations"`
}

func (r *GetOrganizationRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetOrganizationRequest) ToInput() organizationservice.GetOrganizationInput {
	return organizationservice.GetOrganizationInput{
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
