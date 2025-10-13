package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListProjectDocumentsRequest struct {
	Title         *string `json:"title"`
	Page          *int    `json:"page"`
	PerPage       *int    `json:"perPage"`
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
}

func (r *ListProjectDocumentsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListProjectDocumentsRequest) ToInput() projectservice.ListProjectDocumentsInput {
	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var titleFilter *core.ComparableFilter[string] = nil
	if r.Title != nil {
		titleFilter = &core.ComparableFilter[string]{
			Like: r.Title,
		}
	}

	return projectservice.ListProjectDocumentsInput{
		Filters: projectrepo.ProjectDocumentVersionManagerFilters{
			Title: titleFilter,
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
