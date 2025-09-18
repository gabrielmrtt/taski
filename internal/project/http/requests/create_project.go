package project_http_requests

import (
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_services "github.com/gabrielmrtt/taski/internal/project/services"
)

type CreateProjectRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Color         string `json:"color"`
	PriorityLevel int8   `json:"priorityLevel"`
	StartAt       *int64 `json:"startAt"`
	EndAt         *int64 `json:"endAt"`
}

func (r *CreateProjectRequest) ToInput() project_services.CreateProjectInput {
	return project_services.CreateProjectInput{
		Name:          r.Name,
		Description:   r.Description,
		Color:         r.Color,
		PriorityLevel: project_core.ProjectPriorityLevels(r.PriorityLevel),
		StartAt:       r.StartAt,
		EndAt:         r.EndAt,
	}
}
