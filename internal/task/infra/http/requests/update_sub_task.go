package taskhttprequests

import taskservice "github.com/gabrielmrtt/taski/internal/task/services"

type UpdateSubTaskRequest struct {
	Name      *string `json:"name"`
	Completed *bool   `json:"completed"`
}

func (r *UpdateSubTaskRequest) ToInput() taskservice.UpdateSubTaskInput {
	return taskservice.UpdateSubTaskInput{
		Name:      r.Name,
		Completed: r.Completed,
	}
}
