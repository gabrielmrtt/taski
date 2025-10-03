package workspacehttprequests

import workspaceservice "github.com/gabrielmrtt/taski/internal/workspace/service"

type CreateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (r *CreateWorkspaceRequest) ToInput() workspaceservice.CreateWorkspaceInput {
	return workspaceservice.CreateWorkspaceInput{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
	}
}
