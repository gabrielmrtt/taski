package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListMyOrganizationInvitesRequest struct {
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"per_page" schema:"per_page"`
	SortBy        *string `json:"sort_by" schema:"sort_by"`
	SortDirection *string `json:"sort_direction" schema:"sort_direction"`
}

func (r *ListMyOrganizationInvitesRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListMyOrganizationInvitesRequest) ToInput() organization_services.ListMyOrganizationInvitesInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return organization_services.ListMyOrganizationInvitesInput{
		Pagination: &core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: &core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
	}
}
