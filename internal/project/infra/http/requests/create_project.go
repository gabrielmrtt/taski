package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
)

type CreateProjectRequest struct {
	WorkspaceId   string  `json:"workspaceId"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Color         string  `json:"color"`
	PriorityLevel int8    `json:"priorityLevel"`
	StartAt       *string `json:"startAt"`
	EndAt         *string `json:"endAt"`
}

func (r *CreateProjectRequest) ToInput() projectservice.CreateProjectInput {
	var startAt *core.DateTime = nil
	if r.StartAt != nil {
		d, err := core.NewDateTimeFromRFC3339(*r.StartAt)
		if err != nil {
			return projectservice.CreateProjectInput{}
		}
		startAt = &d
	}

	var endAt *core.DateTime = nil
	if r.EndAt != nil {
		d, err := core.NewDateTimeFromRFC3339(*r.EndAt)
		if err != nil {
			return projectservice.CreateProjectInput{}
		}
		endAt = &d
	}

	return projectservice.CreateProjectInput{
		WorkspaceIdentity: core.NewIdentityFromPublic(r.WorkspaceId),
		Name:              r.Name,
		Description:       r.Description,
		Color:             r.Color,
		PriorityLevel:     project.ProjectPriorityLevels(r.PriorityLevel),
		StartAt:           startAt,
		EndAt:             endAt,
	}
}
