package organization

import (
	"github.com/gabrielmrtt/taski/internal/role"
	"github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
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
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type OrganizationUserDto struct {
	OrganizationId string        `json:"organizationId"`
	Role           *role.RoleDto `json:"role,omitempty"`
	User           *user.UserDto `json:"user,omitempty"`
	Status         string        `json:"status"`
}

func OrganizationUserToDto(organizationUser *OrganizationUser) *OrganizationUserDto {
	return &OrganizationUserDto{
		OrganizationId: organizationUser.OrganizationIdentity.Public,
		Role:           role.RoleToDto(&organizationUser.Role),
		User:           user.UserToDto(&organizationUser.User),
		Status:         string(organizationUser.Status),
	}
}
