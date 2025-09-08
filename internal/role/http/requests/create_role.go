package role_http_requests

import (
	role_services "github.com/gabrielmrtt/taski/internal/role/services"
)

type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (r *CreateRoleRequest) ToInput() role_services.CreateRoleInput {
	return role_services.CreateRoleInput{
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
	}
}
