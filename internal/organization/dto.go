package organization

import (
	"github.com/gabrielmrtt/taski/internal/role"
	"github.com/gabrielmrtt/taski/internal/user"
)

type OrganizationDto struct {
	Id            string        `json:"id"`
	Name          string        `json:"name"`
	Status        string        `json:"status"`
	UserCreatorId string        `json:"userCreatorId"`
	UserEditorId  *string       `json:"userEditorId"`
	Creator       *user.UserDto `json:"creator,omitempty"`
	Editor        *user.UserDto `json:"editor,omitempty"`
	CreatedAt     string        `json:"createdAt"`
	UpdatedAt     *string       `json:"updatedAt"`
}

func OrganizationToDto(organization *Organization) *OrganizationDto {
	var updatedAt *string = nil
	if organization.Timestamps.UpdatedAt != nil {
		updatedAtString := organization.Timestamps.UpdatedAt.ToRFC3339()
		updatedAt = &updatedAtString
	}

	var userCreatorId string
	if organization.UserCreatorIdentity != nil {
		userCreatorId = organization.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if organization.UserEditorIdentity != nil {
		userEditorId = &organization.UserEditorIdentity.Public
	}

	var creator *user.UserDto = nil
	if organization.Creator != nil {
		creator = user.UserToDto(organization.Creator)
	}

	var editor *user.UserDto = nil
	if organization.Editor != nil {
		editor = user.UserToDto(organization.Editor)
	}

	return &OrganizationDto{
		Id:            organization.Identity.Public,
		Name:          organization.Name,
		Status:        string(organization.Status),
		UserCreatorId: userCreatorId,
		UserEditorId:  userEditorId,
		Creator:       creator,
		Editor:        editor,
		CreatedAt:     organization.Timestamps.CreatedAt.ToRFC3339(),
		UpdatedAt:     updatedAt,
	}
}

type OrganizationUserDto struct {
	OrganizationId string        `json:"organizationId"`
	Role           *role.RoleDto `json:"role,omitempty"`
	User           *user.UserDto `json:"user,omitempty"`
	Status         string        `json:"status"`
	LastAccessAt   *string       `json:"lastAccessAt"`
}

func OrganizationUserToDto(organizationUser *OrganizationUser) *OrganizationUserDto {
	var lastAccessAt *string = nil
	if organizationUser.LastAccessAt != nil {
		lastAccessAtString := organizationUser.LastAccessAt.ToRFC3339()
		lastAccessAt = &lastAccessAtString
	}

	return &OrganizationUserDto{
		OrganizationId: organizationUser.OrganizationIdentity.Public,
		Role:           role.RoleToDto(&organizationUser.Role),
		User:           user.UserToDto(&organizationUser.User),
		Status:         string(organizationUser.Status),
		LastAccessAt:   lastAccessAt,
	}
}
