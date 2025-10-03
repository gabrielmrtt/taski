package organizationhttprequests

import organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

func (r *CreateOrganizationRequest) ToInput() organizationservice.CreateOrganizationInput {
	return organizationservice.CreateOrganizationInput{
		Name: r.Name,
	}
}
