package project_http_requests

import (
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_services "github.com/gabrielmrtt/taski/internal/project/services"
)

type UpdateProjectRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Color         *string `json:"color"`
	Status        *string `json:"status"`
	PriorityLevel *int8   `json:"priorityLevel"`
	StartAt       *int64  `json:"startAt"`
	EndAt         *int64  `json:"endAt"`
}

func (r *UpdateProjectRequest) ToInput() project_services.UpdateProjectInput {
	var status *project_core.ProjectStatuses = nil
	if r.Status != nil {
		projectStatus := project_core.ProjectStatuses(*r.Status)
		status = &projectStatus
	}

	var priorityLevel *project_core.ProjectPriorityLevels = nil
	if r.PriorityLevel != nil {
		projectPriorityLevel := project_core.ProjectPriorityLevels(*r.PriorityLevel)
		priorityLevel = &projectPriorityLevel
	}

	return project_services.UpdateProjectInput{
		Name:          r.Name,
		Description:   r.Description,
		Color:         r.Color,
		Status:        status,
		PriorityLevel: priorityLevel,
		StartAt:       r.StartAt,
		EndAt:         r.EndAt,
	}
}
