package teamservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
)

type ListTeamsService struct {
	TeamRepository teamrepo.TeamRepository
}

func NewListTeamsService(teamRepository teamrepo.TeamRepository) *ListTeamsService {
	return &ListTeamsService{
		TeamRepository: teamRepository,
	}
}

type ListTeamsInput struct {
	Filters    teamrepo.TeamFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

func (i ListTeamsInput) Validate() error {
	return nil
}

func (s *ListTeamsService) Execute(input ListTeamsInput) (*core.PaginationOutput[team.TeamDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	teams, err := s.TeamRepository.PaginateTeamsBy(teamrepo.PaginateTeamsParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var teamsDto []team.TeamDto = make([]team.TeamDto, 0)
	for _, tm := range teams.Data {
		teamsDto = append(teamsDto, *team.TeamToDto(&tm))
	}

	return &core.PaginationOutput[team.TeamDto]{
		Data:    teamsDto,
		Page:    teams.Page,
		HasMore: teams.HasMore,
		Total:   teams.Total,
	}, nil
}
