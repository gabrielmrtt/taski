package workspace_http_requests

import (
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_services "github.com/gabrielmrtt/taski/internal/workspace/services"
)

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Status      *string `json:"status"`
}

func (r *UpdateWorkspaceRequest) ToInput() workspace_services.UpdateWorkspaceInput {
	var status *workspace_core.WorkspaceStatuses = nil
	if r.Status != nil {
		workspaceStatus := workspace_core.WorkspaceStatuses(*r.Status)
		status = &workspaceStatus
	}

	return workspace_services.UpdateWorkspaceInput{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
		Status:      status,
	}
}
