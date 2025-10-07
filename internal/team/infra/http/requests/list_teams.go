package teamhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
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

func (r *ListTeamsRequest) ToInput() teamservice.ListTeamsInput {
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

	var statusFilter *core.ComparableFilter[team.TeamStatuses] = nil
	if r.Status != nil {
		teamStatus := team.TeamStatuses(*r.Status)
		statusFilter = &core.ComparableFilter[team.TeamStatuses]{
			Equals: &teamStatus,
		}
	}

	return teamservice.ListTeamsInput{
		Filters: teamrepo.TeamFilters{
			Name:   nameFilter,
			Status: statusFilter,
		},
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
	}
}
