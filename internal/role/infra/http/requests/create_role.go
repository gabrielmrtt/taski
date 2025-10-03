package rolehttprequests

import (
	roleservice "github.com/gabrielmrtt/taski/internal/role/service"
)

type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (r *CreateRoleRequest) ToInput() roleservice.CreateRoleInput {
	return roleservice.CreateRoleInput{
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
	}
}
