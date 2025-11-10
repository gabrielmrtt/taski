package taskhttprequests

import taskservice "github.com/gabrielmrtt/taski/internal/task/services"

type AddSubTaskRequest struct {
	Name string `json:"name"`
}

func (r *AddSubTaskRequest) ToInput() taskservice.AddSubTaskInput {
	return taskservice.AddSubTaskInput{
		Name: r.Name,
	}
}
