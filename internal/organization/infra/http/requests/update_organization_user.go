package organizationhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
)

type UpdateOrganizationUserWorkspaceInput struct {
	WorkspaceId string   `json:"workspaceId"`
	Projects    []string `json:"projects"`
}

type UpdateOrganizationUserRequest struct {
	RoleId     *string                                `json:"roleId"`
	Status     *string                                `json:"status"`
	Workspaces []UpdateOrganizationUserWorkspaceInput `json:"workspaces"`
}

func (r *UpdateOrganizationUserRequest) ToInput() organizationservice.UpdateOrganizationUserInput {
	var roleIdentity *core.Identity = nil
	if r.RoleId != nil {
		identity := core.NewIdentityFromPublic(*r.RoleId)
		roleIdentity = &identity
	}

	var status *organization.OrganizationUserStatuses = nil
	if r.Status != nil {
		organizationStatus := organization.OrganizationUserStatuses(*r.Status)
		status = &organizationStatus
	}

	var workspaces []organizationservice.UpdateOrganizationUserWorkspaceInput = make([]organizationservice.UpdateOrganizationUserWorkspaceInput, 0)
	for _, w := range r.Workspaces {
		var projects []core.Identity = make([]core.Identity, 0)
		for _, p := range w.Projects {
			projects = append(projects, core.NewIdentityFromPublic(p))
		}

		workspaces = append(workspaces, organizationservice.UpdateOrganizationUserWorkspaceInput{
			WorkspaceIdentity: core.NewIdentityFromPublic(w.WorkspaceId),
			Projects:          projects,
		})
	}

	return organizationservice.UpdateOrganizationUserInput{
		RoleIdentity: roleIdentity,
		Status:       status,
		Workspaces:   workspaces,
	}
}
