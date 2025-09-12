package organization_http_requests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
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

func (r *ListMyOrganizationInvitesRequest) ToInput() organization_services.ListMyOrganizationInvitesInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != nil {
		relationsInput = strings.Split(stringutils.CamelCaseToPascalCase(*r.Relations), ",")
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
		RelationsInput: relationsInput,
	}
}
