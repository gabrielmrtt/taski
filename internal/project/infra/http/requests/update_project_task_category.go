package projecthttprequests

import projectservice "github.com/gabrielmrtt/taski/internal/project/service"

type UpdateProjectTaskCategoryRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

func (r *UpdateProjectTaskCategoryRequest) ToInput() projectservice.UpdateProjectTaskCategoryInput {
	return projectservice.UpdateProjectTaskCategoryInput{
		Name:  r.Name,
		Color: r.Color,
	}
}
