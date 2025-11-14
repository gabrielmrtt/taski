package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
)

type ChangeTaskStatusRequest struct {
	StatusId     *string `json:"statusId"`
	AdvanceOrder bool    `json:"advanceOrder"`
}

func (r *ChangeTaskStatusRequest) ToInput() taskservice.ChangeTaskStatusInput {
	var projectTaskStatusIdentity *core.Identity = nil
	if r.StatusId != nil {
		identity := core.NewIdentityFromPublic(*r.StatusId)
		projectTaskStatusIdentity = &identity
	}

	return taskservice.ChangeTaskStatusInput{
		ProjectTaskStatusIdentity: projectTaskStatusIdentity,
		AdvanceOrder:              r.AdvanceOrder,
	}
}
