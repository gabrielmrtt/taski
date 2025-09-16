package organization_core

import (
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type OrganizationDto struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	UserCreatorId string  `json:"userCreatorId"`
	UserEditorId  *string `json:"userEditorId"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     *string `json:"updatedAt"`
}

func OrganizationToDto(organization *Organization) *OrganizationDto {
	createdAt := datetimeutils.EpochToRFC3339(*organization.Timestamps.CreatedAt)

	var updatedAt *string = nil
	if organization.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*organization.Timestamps.UpdatedAt)
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
	OrganizationId string             `json:"organizationId"`
	Role           *role_core.RoleDto `json:"role,omitempty"`
	User           *user_core.UserDto `json:"user,omitempty"`
	Status         string             `json:"status"`
}

func OrganizationUserToDto(organizationUser *OrganizationUser) *OrganizationUserDto {
	return &OrganizationUserDto{
		OrganizationId: organizationUser.OrganizationIdentity.Public,
		Role:           role_core.RoleToDto(&organizationUser.Role),
		User:           user_core.UserToDto(&organizationUser.User),
		Status:         string(organizationUser.Status),
	}
}
