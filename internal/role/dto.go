package role_core

import (
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type RoleDto struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Permissions     []string `json:"permissions"`
	OrganizationId  *string  `json:"organization_id"`
	IsSystemDefault bool     `json:"is_system_default"`
	UserCreatorId   *string  `json:"user_creator_id"`
	UserEditorId    *string  `json:"user_editor_id"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       *string  `json:"updated_at"`
}

func RoleToDto(role *Role) *RoleDto {
	var permissions []string = make([]string, 0)
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Slug)
	}

	createdAt := datetimeutils.EpochToRFC3339(*role.Timestamps.CreatedAt)
	var updatedAt *string = nil
	if role.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*role.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
	}

	var organizationId *string = nil

	if role.OrganizationIdentity != nil {
		organizationId = &role.OrganizationIdentity.Public
	}

	var userCreatorId *string = nil
	if role.UserCreatorIdentity != nil {
		userCreatorId = &role.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if role.UserEditorIdentity != nil {
		userEditorId = &role.UserEditorIdentity.Public
	}

	return &RoleDto{
		Id:              role.Identity.Public,
		Name:            role.Name,
		Description:     role.Description,
		Permissions:     permissions,
		OrganizationId:  organizationId,
		IsSystemDefault: role.IsSystemDefault,
		UserCreatorId:   userCreatorId,
		UserEditorId:    userEditorId,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

type PermissionDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func PermissionToDto(permission *Permission) *PermissionDto {
	return &PermissionDto{
		Name:        permission.Name,
		Description: permission.Description,
		Slug:        permission.Slug,
	}
}
