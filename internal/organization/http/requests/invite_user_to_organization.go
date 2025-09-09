package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
)

type InviteUserToOrganizationRequest struct {
	Email  string `json:"email"`
	RoleId string `json:"role_id"`
}

func (r *InviteUserToOrganizationRequest) ToInput() organization_services.InviteUserToOrganizationInput {
	return organization_services.InviteUserToOrganizationInput{
		Email:        r.Email,
		RoleIdentity: core.NewIdentityFromPublic(r.RoleId),
	}
}
