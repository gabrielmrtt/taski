package organizationhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
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

func (r *InviteUserToOrganizationRequest) ToInput() organizationservice.InviteUserToOrganizationInput {
	var workspaces []organizationservice.InviteUserToOrganizationWorkspaceInput = make([]organizationservice.InviteUserToOrganizationWorkspaceInput, 0)
	for _, w := range r.Workspaces {
		var projects []core.Identity = make([]core.Identity, 0)
		for _, p := range w.Projects {
			projects = append(projects, core.NewIdentityFromPublic(p))
		}

		workspaces = append(workspaces, organizationservice.InviteUserToOrganizationWorkspaceInput{
			WorkspaceIdentity: core.NewIdentityFromPublic(w.WorkspaceId),
			Projects:          projects,
		})
	}

	return organizationservice.InviteUserToOrganizationInput{
		Email:        r.Email,
		RoleIdentity: core.NewIdentityFromPublic(r.RoleId),
		Workspaces:   workspaces,
	}
}
