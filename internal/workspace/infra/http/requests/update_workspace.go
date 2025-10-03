package workspacehttprequests

import (
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspaceservice "github.com/gabrielmrtt/taski/internal/workspace/service"
)

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Status      *string `json:"status"`
}

func (r *UpdateWorkspaceRequest) ToInput() workspaceservice.UpdateWorkspaceInput {
	var status *workspace.WorkspaceStatuses = nil
	if r.Status != nil {
		workspaceStatus := workspace.WorkspaceStatuses(*r.Status)
		status = &workspaceStatus
	}

	return workspaceservice.UpdateWorkspaceInput{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
		Status:      status,
	}
}
