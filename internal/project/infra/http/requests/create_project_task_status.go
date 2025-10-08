package projecthttprequests

import projectservice "github.com/gabrielmrtt/taski/internal/project/service"

type CreateProjectTaskStatusRequest struct {
	Name                     string `json:"name"`
	Color                    string `json:"color"`
	ShouldSetTaskToCompleted bool   `json:"shouldSetTaskToCompleted"`
	IsDefault                bool   `json:"isDefault"`
	ShouldUseOrder           bool   `json:"shouldUseOrder"`
}

func (r *CreateProjectTaskStatusRequest) ToInput() projectservice.CreateProjectTaskStatusInput {
	return projectservice.CreateProjectTaskStatusInput{
		Name:                     r.Name,
		Color:                    r.Color,
		ShouldSetTaskToCompleted: r.ShouldSetTaskToCompleted,
		IsDefault:                r.IsDefault,
		ShouldUseOrder:           r.ShouldUseOrder,
	}
}
