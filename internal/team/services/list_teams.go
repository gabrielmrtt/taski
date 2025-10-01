package team_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
)

type ListTeamsService struct {
	TeamRepository team_repositories.TeamRepository
}

func NewListTeamsService(teamRepository team_repositories.TeamRepository) *ListTeamsService {
	return &ListTeamsService{
		TeamRepository: teamRepository,
	}
}

type ListTeamsInput struct {
	OrganizationIdentity core.Identity
	Filters              team_repositories.TeamFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
}

func (i ListTeamsInput) Validate() error {
	return nil
}

func (s *ListTeamsService) Execute(input ListTeamsInput) (*core.PaginationOutput[team_core.TeamDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Filters.OrganizationIdentity = input.OrganizationIdentity

	teams, err := s.TeamRepository.PaginateTeamsBy(team_repositories.PaginateTeamsParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var teamsDto []team_core.TeamDto = make([]team_core.TeamDto, 0)
	for _, team := range teams.Data {
		teamsDto = append(teamsDto, *team_core.TeamToDto(&team))
	}

	return &core.PaginationOutput[team_core.TeamDto]{
		Data:    teamsDto,
		Page:    teams.Page,
		HasMore: teams.HasMore,
		Total:   teams.Total,
	}, nil
}
