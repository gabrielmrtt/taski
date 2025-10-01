package team_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
	team_services "github.com/gabrielmrtt/taski/internal/team/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListTeamsRequest struct {
	Name          *string `json:"name"`
	Status        *string `json:"status"`
	Page          *int    `json:"page"`
	PerPage       *int    `json:"perPage"`
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
	Relations     *string `json:"relations"`
}

func (r *ListTeamsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListTeamsRequest) ToInput() team_services.ListTeamsInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var statusFilter *core.ComparableFilter[team_core.TeamStatuses] = nil
	if r.Status != nil {
		teamStatus := team_core.TeamStatuses(*r.Status)
		statusFilter = &core.ComparableFilter[team_core.TeamStatuses]{
			Equals: &teamStatus,
		}
	}

	return team_services.ListTeamsInput{
		Filters: team_repositories.TeamFilters{
			Name:   nameFilter,
			Status: statusFilter,
		},
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
