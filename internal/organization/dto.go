package organization_core

import (
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type OrganizationDto struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	UserCreatorId string `json:"user_creator_id"`
	UserEditorId  string `json:"user_editor_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func OrganizationToDto(organization *Organization) *OrganizationDto {
	createdAt := datetimeutils.EpochToRFC3339(*organization.Timestamps.CreatedAt)

	var updatedAt string
	if organization.Timestamps.UpdatedAt != nil {
		updatedAt = datetimeutils.EpochToRFC3339(*organization.Timestamps.UpdatedAt)
	}

	var userCreatorId string
	if organization.UserCreatorIdentity != nil {
		userCreatorId = organization.UserCreatorIdentity.Public
	}

	var userEditorId string
	if organization.UserEditorIdentity != nil {
		userEditorId = organization.UserEditorIdentity.Public
	}

	return &OrganizationDto{
		Id:            organization.Identity.Public,
		Name:          organization.Name,
		Status:        string(organization.Status),
		UserCreatorId: userCreatorId,
		UserEditorId:  userEditorId,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type OrganizationUserDto struct {
	OrganizationId string             `json:"organization_id"`
	User           *user_core.UserDto `json:"user"`
	Role           *role_core.RoleDto `json:"role"`
	Status         string             `json:"status"`
}

func OrganizationUserToDto(organizationUser *OrganizationUser) *OrganizationUserDto {
	return &OrganizationUserDto{
		OrganizationId: organizationUser.OrganizationIdentity.Public,
		User:           user_core.UserToDto(organizationUser.User),
		Role:           role_core.RoleToDto(organizationUser.Role),
		Status:         string(organizationUser.Status),
	}
}
