package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type GetTaskHistoryRequest struct {
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
	Page          *int    `json:"page"`
	PerPage       *int    `json:"perPage"`
	Relations     *string `json:"relations"`
}

func (r *GetTaskHistoryRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *GetTaskHistoryRequest) ToInput() taskservice.GetTaskHistoryInput {
	var sortDirection *core.SortDirection = nil
	if r.SortDirection != nil {
		s := core.SortDirection(*r.SortDirection)
		sortDirection = &s
	}

	return taskservice.GetTaskHistoryInput{
		SortInput:       core.SortInput{By: r.SortBy, Direction: sortDirection},
		PaginationInput: core.PaginationInput{Page: r.Page, PerPage: r.PerPage},
		RelationsInput:  corehttp.GetRelationsInput(r.Relations),
	}
}
