package team

import (
	"github.com/gabrielmrtt/taski/internal/organization"
	"github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type TeamDto struct {
	Id            string                        `json:"id"`
	Name          string                        `json:"name"`
	Description   string                        `json:"description"`
	Status        string                        `json:"status"`
	UserCreatorId string                        `json:"userCreatorId"`
	UserEditorId  *string                       `json:"userEditorId"`
	Creator       *user.UserDto                 `json:"creator,omitempty"`
	Editor        *user.UserDto                 `json:"editor,omitempty"`
	Organization  *organization.OrganizationDto `json:"organization,omitempty"`
	CreatedAt     string                        `json:"createdAt"`
	UpdatedAt     *string                       `json:"updatedAt"`
	Members       []TeamUserDto                 `json:"members"`
}

func TeamToDto(team *Team) *TeamDto {
	createdAt := datetimeutils.EpochToRFC3339(*team.Timestamps.CreatedAt)
	var updatedAt *string = nil
	if team.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*team.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
	}

	var userCreatorId *string = nil
	if team.UserCreatorIdentity != nil {
		userCreatorId = &team.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if team.UserEditorIdentity != nil {
		userEditorId = &team.UserEditorIdentity.Public
	}

	var members []TeamUserDto = make([]TeamUserDto, 0)
	for _, user := range team.Members {
		members = append(members, *TeamUserToDto(&user))
	}

	var creator *user.UserDto = nil
	if team.Creator != nil {
		creator = user.UserToDto(team.Creator)
	}

	var editor *user.UserDto = nil
	if team.Editor != nil {
		editor = user.UserToDto(team.Editor)
	}

	var org *organization.OrganizationDto = nil
	if team.Organization != nil {
		org = organization.OrganizationToDto(team.Organization)
	}

	return &TeamDto{
		Id:            team.Identity.Public,
		Name:          team.Name,
		Description:   team.Description,
		Status:        string(team.Status),
		UserCreatorId: *userCreatorId,
		UserEditorId:  userEditorId,
		Creator:       creator,
		Editor:        editor,
		Organization:  org,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Members:       members,
	}
}

type TeamUserDto struct {
	User *user.UserDto `json:"user"`
}

func TeamUserToDto(teamUser *TeamUser) *TeamUserDto {
	return &TeamUserDto{
		User: user.UserToDto(&teamUser.User),
	}
}
