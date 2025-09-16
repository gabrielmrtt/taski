package workspace_http_requests

import workspace_services "github.com/gabrielmrtt/taski/internal/workspace/services"

type CreateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (r *CreateWorkspaceRequest) ToInput() workspace_services.CreateWorkspaceInput {
	return workspace_services.CreateWorkspaceInput{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
	}
}
