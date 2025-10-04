package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
)

type UpdateProjectRequest struct {
	WorkspaceId   *string `json:"workspaceId"`
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Color         *string `json:"color"`
	Status        *string `json:"status"`
	PriorityLevel *int8   `json:"priorityLevel"`
	StartAt       *int64  `json:"startAt"`
	EndAt         *int64  `json:"endAt"`
}

func (r *UpdateProjectRequest) ToInput() projectservice.UpdateProjectInput {
	var status *project.ProjectStatuses = nil
	if r.Status != nil {
		projectStatus := project.ProjectStatuses(*r.Status)
		status = &projectStatus
	}

	var priorityLevel *project.ProjectPriorityLevels = nil
	if r.PriorityLevel != nil {
		projectPriorityLevel := project.ProjectPriorityLevels(*r.PriorityLevel)
		priorityLevel = &projectPriorityLevel
	}

	var workspaceIdentity *core.Identity = nil
	if r.WorkspaceId != nil {
		identity := core.NewIdentityFromPublic(*r.WorkspaceId)
		workspaceIdentity = &identity
	}

	return projectservice.UpdateProjectInput{
		WorkspaceIdentity: workspaceIdentity,
		Name:              r.Name,
		Description:       r.Description,
		Color:             r.Color,
		Status:            status,
		PriorityLevel:     priorityLevel,
		StartAt:           r.StartAt,
		EndAt:             r.EndAt,
	}
}
