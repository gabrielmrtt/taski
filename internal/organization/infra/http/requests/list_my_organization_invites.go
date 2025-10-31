package organizationhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListMyOrganizationInvitesRequest struct {
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"perPage" schema:"perPage"`
	SortBy        *string `json:"sortBy" schema:"sortBy"`
	SortDirection *string `json:"sortDirection" schema:"sortDirection"`
	Relations     *string `json:"relations" schema:"relations"`
}

func (r *ListMyOrganizationInvitesRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListMyOrganizationInvitesRequest) ToInput() organizationservice.ListMyOrganizationInvitesInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return organizationservice.ListMyOrganizationInvitesInput{
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
		RelationsInput: corehttp.GetRelationsInput(*r.Relations),
	}
}
