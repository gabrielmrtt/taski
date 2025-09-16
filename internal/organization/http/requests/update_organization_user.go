package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
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

func (r *UpdateOrganizationUserRequest) ToInput() organization_services.UpdateOrganizationUserInput {
	var roleIdentity *core.Identity = nil
	if r.RoleId != nil {
		identity := core.NewIdentityFromPublic(*r.RoleId)
		roleIdentity = &identity
	}

	var status *organization_core.OrganizationUserStatuses = nil
	if r.Status != nil {
		organizationStatus := organization_core.OrganizationUserStatuses(*r.Status)
		status = &organizationStatus
	}

	var workspaces []organization_services.UpdateOrganizationUserWorkspaceInput = make([]organization_services.UpdateOrganizationUserWorkspaceInput, 0)
	for _, w := range r.Workspaces {
		var projects []core.Identity = make([]core.Identity, 0)
		for _, p := range w.Projects {
			projects = append(projects, core.NewIdentityFromPublic(p))
		}

		workspaces = append(workspaces, organization_services.UpdateOrganizationUserWorkspaceInput{
			WorkspaceIdentity: core.NewIdentityFromPublic(w.WorkspaceId),
			Projects:          projects,
		})
	}

	return organization_services.UpdateOrganizationUserInput{
		RoleIdentity: roleIdentity,
		Status:       status,
		Workspaces:   workspaces,
	}
}
