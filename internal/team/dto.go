package team_core

import (
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type TeamDto struct {
	Id            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`
	UserCreatorId string        `json:"userCreatorId"`
	UserEditorId  *string       `json:"userEditorId"`
	CreatedAt     string        `json:"createdAt"`
	UpdatedAt     *string       `json:"updatedAt"`
	Members       []TeamUserDto `json:"members"`
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

	return &TeamDto{
		Id:            team.Identity.Public,
		Name:          team.Name,
		Description:   team.Description,
		Status:        string(team.Status),
		UserCreatorId: *userCreatorId,
		UserEditorId:  userEditorId,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Members:       members,
	}
}

type TeamUserDto struct {
	User *user_core.UserDto `json:"user"`
}

func TeamUserToDto(teamUser *TeamUser) *TeamUserDto {
	return &TeamUserDto{
		User: user_core.UserToDto(&teamUser.User),
	}
}
