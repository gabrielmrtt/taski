package team_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_core "github.com/gabrielmrtt/taski/internal/team"
)

type TeamFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Status               *core.ComparableFilter[team_core.TeamStatuses]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
}

type GetTeamByIdentityParams struct {
	TeamIdentity         core.Identity
	OrganizationIdentity *core.Identity
}

type PaginateTeamsParams struct {
	Filters    TeamFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type StoreTeamParams struct {
	Team *team_core.Team
}

type UpdateTeamParams struct {
	Team *team_core.Team
}

type DeleteTeamParams struct {
	TeamIdentity core.Identity
}

type TeamRepository interface {
	SetTransaction(tx core.Transaction) error

	GetTeamByIdentity(params GetTeamByIdentityParams) (*team_core.Team, error)
	PaginateTeamsBy(params PaginateTeamsParams) (*core.PaginationOutput[team_core.Team], error)

	StoreTeam(params StoreTeamParams) (*team_core.Team, error)
	UpdateTeam(params UpdateTeamParams) error
	DeleteTeam(params DeleteTeamParams) error
}
