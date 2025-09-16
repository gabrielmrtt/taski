package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
)

type InviteUserToOrganizationWorkspaceInput struct {
	WorkspaceId string   `json:"workspaceId"`
	Projects    []string `json:"projects"`
}

type InviteUserToOrganizationRequest struct {
	Email      string                                   `json:"email"`
	RoleId     string                                   `json:"roleId"`
	Workspaces []InviteUserToOrganizationWorkspaceInput `json:"workspaces"`
}

func (r *InviteUserToOrganizationRequest) ToInput() organization_services.InviteUserToOrganizationInput {
	var workspaces []organization_services.InviteUserToOrganizationWorkspaceInput = make([]organization_services.InviteUserToOrganizationWorkspaceInput, 0)
	for _, w := range r.Workspaces {
		var projects []core.Identity = make([]core.Identity, 0)
		for _, p := range w.Projects {
			projects = append(projects, core.NewIdentityFromPublic(p))
		}

		workspaces = append(workspaces, organization_services.InviteUserToOrganizationWorkspaceInput{
			WorkspaceIdentity: core.NewIdentityFromPublic(w.WorkspaceId),
			Projects:          projects,
		})
	}

	return organization_services.InviteUserToOrganizationInput{
		Email:        r.Email,
		RoleIdentity: core.NewIdentityFromPublic(r.RoleId),
		Workspaces:   workspaces,
	}
}
