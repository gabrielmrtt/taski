package teamrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/team"
)

type TeamFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Status               *core.ComparableFilter[team.TeamStatuses]
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
	Team *team.Team
}

type UpdateTeamParams struct {
	Team *team.Team
}

type DeleteTeamParams struct {
	TeamIdentity core.Identity
}

type TeamRepository interface {
	SetTransaction(tx core.Transaction) error

	GetTeamByIdentity(params GetTeamByIdentityParams) (*team.Team, error)
	PaginateTeamsBy(params PaginateTeamsParams) (*core.PaginationOutput[team.Team], error)

	StoreTeam(params StoreTeamParams) (*team.Team, error)
	UpdateTeam(params UpdateTeamParams) error
	DeleteTeam(params DeleteTeamParams) error
}
