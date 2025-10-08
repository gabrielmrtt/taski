package projecthttprequests

import projectservice "github.com/gabrielmrtt/taski/internal/project/service"

type CreateProjectTaskCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (r *CreateProjectTaskCategoryRequest) ToInput() projectservice.CreateProjectTaskCategoryInput {
	return projectservice.CreateProjectTaskCategoryInput{
		Name:  r.Name,
		Color: r.Color,
	}
}
