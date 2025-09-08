package organization_http_requests

import organization_services "github.com/gabrielmrtt/taski/internal/organization/services"

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

func (r *CreateOrganizationRequest) ToInput() organization_services.CreateOrganizationInput {
	return organization_services.CreateOrganizationInput{
		Name: r.Name,
	}
}
