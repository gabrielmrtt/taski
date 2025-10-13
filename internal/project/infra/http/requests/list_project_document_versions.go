package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListProjectDocumentVersionsRequest struct {
	Version       *string `json:"version"`
	Page          *int    `json:"page"`
	PerPage       *int    `json:"perPage"`
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
}

func (r *ListProjectDocumentVersionsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListProjectDocumentVersionsRequest) ToInput() projectservice.ListProjectDocumentVersionsInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var versionFilter *core.ComparableFilter[string] = nil
	if r.Version != nil {
		versionFilter = &core.ComparableFilter[string]{
			Like: r.Version,
		}
	}

	return projectservice.ListProjectDocumentVersionsInput{
		Filters: projectrepo.ProjectDocumentVersionFilters{
			Version: versionFilter,
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
