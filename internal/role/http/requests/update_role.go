package role_http_requests

import role_services "github.com/gabrielmrtt/taski/internal/role/services"

type UpdateRoleRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

func (r *UpdateRoleRequest) ToInput() role_services.UpdateRoleInput {
	return role_services.UpdateRoleInput{
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
	}
}
