package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
)

type RemoveSubTaskRequest struct {
	SubTaskIdentity string `json:"subTaskIdentity"`
}

func (r *RemoveSubTaskRequest) ToInput() taskservice.RemoveSubTaskInput {
	return taskservice.RemoveSubTaskInput{
		SubTaskIdentity: core.NewIdentity(r.SubTaskIdentity),
	}
}
