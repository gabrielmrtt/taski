package team_core

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type Team struct {
	Identity             core.Identity
	OrganizationIdentity core.Identity
	Name                 string
	Description          string
	Status               TeamStatuses
	UserCreatorIdentity  *core.Identity
	UserEditorIdentity   *core.Identity
	Timestamps           core.Timestamps

	Members []TeamUser
}

type NewTeamInput struct {
	Name                 string
	Description          string
	OrganizationIdentity core.Identity
	UserCreatorIdentity  *core.Identity
}

func NewTeam(input NewTeamInput) (*Team, error) {
	now := datetimeutils.EpochNow()

	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewDescription(input.Description); err != nil {
		return nil, err
	}

	return &Team{
		Identity:             core.NewIdentity(TeamIdentityPrefix),
		OrganizationIdentity: input.OrganizationIdentity,
		Name:                 input.Name,
		Description:          input.Description,
		Status:               TeamStatusActive,
		UserCreatorIdentity:  input.UserCreatorIdentity,
		UserEditorIdentity:   nil,
		Timestamps:           core.Timestamps{CreatedAt: &now, UpdatedAt: nil},
	}, nil
}

func (t *Team) ChangeName(name string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	t.Name = name
	t.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Team) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewDescription(description); err != nil {
		return err
	}

	t.Description = description
	t.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Team) Activate() {
	t.Status = TeamStatusActive
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
}

func (t *Team) Inactivate() {
	t.Status = TeamStatusInactive
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
}

func (t *Team) IsActive() bool {
	return t.Status == TeamStatusActive
}

func (t *Team) IsInactive() bool {
	return t.Status == TeamStatusInactive
}

type TeamUser struct {
	TeamIdentity core.Identity
	User         user_core.User
}

func (t *Team) AddUser(user user_core.User) {
	teamUser := TeamUser{
		TeamIdentity: t.Identity,
		User:         user,
	}

	t.Members = append(t.Members, teamUser)
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
}

func (t *Team) RemoveUser(user user_core.User) {
	t.Members = slices.DeleteFunc(t.Members, func(tu TeamUser) bool {
		return tu.User.Identity == user.Identity
	})
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
}

func (t *Team) RemoveAllUsers() {
	t.Members = []TeamUser{}
	now := datetimeutils.EpochNow()
	t.Timestamps.UpdatedAt = &now
}
