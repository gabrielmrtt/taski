package projecthttprequests

import projectservice "github.com/gabrielmrtt/taski/internal/project/service"

type UpdateProjectTaskStatusRequest struct {
	Name                     *string `json:"name"`
	Color                    *string `json:"color"`
	ShouldSetTaskToCompleted *bool   `json:"shouldSetTaskToCompleted"`
	IsDefault                *bool   `json:"isDefault"`
	Order                    *int8   `json:"order"`
}

func (r *UpdateProjectTaskStatusRequest) ToInput() projectservice.UpdateProjectTaskStatusInput {
	return projectservice.UpdateProjectTaskStatusInput{
		Name:                     r.Name,
		Color:                    r.Color,
		ShouldSetTaskToCompleted: r.ShouldSetTaskToCompleted,
		IsDefault:                r.IsDefault,
		Order:                    r.Order,
	}
}
