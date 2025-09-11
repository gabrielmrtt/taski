package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
)

type UpdateOrganizationUserRequest struct {
	RoleId *string `json:"role_id"`
	Status *string `json:"status"`
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

	return organization_services.UpdateOrganizationUserInput{
		RoleIdentity: roleIdentity,
		Status:       status,
	}
}
