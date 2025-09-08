package organization_http_requests

import organization_services "github.com/gabrielmrtt/taski/internal/organization/services"

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

func (r *UpdateOrganizationRequest) ToInput() organization_services.UpdateOrganizationInput {
	return organization_services.UpdateOrganizationInput{
		Name: &r.Name,
	}
}
