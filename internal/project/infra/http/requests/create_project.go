package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
)

type CreateProjectRequest struct {
	WorkspaceId   string `json:"workspaceId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Color         string `json:"color"`
	PriorityLevel int8   `json:"priorityLevel"`
	StartAt       *int64 `json:"startAt"`
	EndAt         *int64 `json:"endAt"`
}

func (r *CreateProjectRequest) ToInput() projectservice.CreateProjectInput {
	return projectservice.CreateProjectInput{
		WorkspaceIdentity: core.NewIdentityFromPublic(r.WorkspaceId),
		Name:              r.Name,
		Description:       r.Description,
		Color:             r.Color,
		PriorityLevel:     project.ProjectPriorityLevels(r.PriorityLevel),
		StartAt:           r.StartAt,
		EndAt:             r.EndAt,
	}
}
