package role

import (
	"github.com/gabrielmrtt/taski/internal/user"
)

type RoleDto struct {
	Id              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Permissions     []string      `json:"permissions"`
	OrganizationId  *string       `json:"organizationId"`
	IsSystemDefault bool          `json:"isSystemDefault"`
	UserCreatorId   *string       `json:"userCreatorId"`
	UserEditorId    *string       `json:"userEditorId"`
	Creator         *user.UserDto `json:"creator,omitempty"`
	Editor          *user.UserDto `json:"editor,omitempty"`
	CreatedAt       string        `json:"createdAt"`
	UpdatedAt       *string       `json:"updatedAt"`
}

func RoleToDto(role *Role) *RoleDto {
	var permissions []string = make([]string, 0)
	for _, permission := range role.Permissions {
		permissions = append(permissions, string(permission.Slug))
	}

	createdAt := role.Timestamps.CreatedAt.ToRFC3339()
	var updatedAt *string = nil
	if role.Timestamps.UpdatedAt != nil {
		updatedAtString := role.Timestamps.UpdatedAt.ToRFC3339()
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

	var creator *user.UserDto = nil
	if role.Creator != nil {
		creator = user.UserToDto(role.Creator)
	}

	var editor *user.UserDto = nil
	if role.Editor != nil {
		editor = user.UserToDto(role.Editor)
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
		Creator:         creator,
		Editor:          editor,
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
		Slug:        string(permission.Slug),
	}
}
