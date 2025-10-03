package rolehttprequests

import roleservice "github.com/gabrielmrtt/taski/internal/role/service"

type UpdateRoleRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

func (r *UpdateRoleRequest) ToInput() roleservice.UpdateRoleInput {
	return roleservice.UpdateRoleInput{
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
	}
}
