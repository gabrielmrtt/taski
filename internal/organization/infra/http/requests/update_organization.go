package organizationhttprequests

import organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

func (r *UpdateOrganizationRequest) ToInput() organizationservice.UpdateOrganizationInput {
	return organizationservice.UpdateOrganizationInput{
		Name: &r.Name,
	}
}
