package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListTaskCommentsRequest struct {
	AuthorId      *string `json:"authorId"`
	PerPage       *int    `json:"perPage"`
	Page          *int    `json:"page"`
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
	Relations     *string `json:"relations"`
}

func (r *ListTaskCommentsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListTaskCommentsRequest) ToInput() taskservice.ListTaskCommentsInput {
	var authorIdentity *core.Identity = nil
	if r.AuthorId != nil {
		identity := core.NewIdentity(*r.AuthorId)
		authorIdentity = &identity
	}

	var sortDirection *core.SortDirection = nil
	if r.SortDirection != nil {
		s := core.SortDirection(*r.SortDirection)
		sortDirection = &s
	}

	return taskservice.ListTaskCommentsInput{
		Filters: taskrepo.TaskCommentFilters{
			AuthorIdentity: authorIdentity,
		},
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: sortDirection,
		},
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
