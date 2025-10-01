package team_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
)

type GetTeamService struct {
	TeamRepository team_repositories.TeamRepository
}

func NewGetTeamService(teamRepository team_repositories.TeamRepository) *GetTeamService {
	return &GetTeamService{
		TeamRepository: teamRepository,
	}
}

type GetTeamInput struct {
	TeamIdentity         core.Identity
	OrganizationIdentity core.Identity
}

func (i GetTeamInput) Validate() error {
	return nil
}

func (s *GetTeamService) Execute(input GetTeamInput) (*team_core.TeamDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	team, err := s.TeamRepository.GetTeamByIdentity(team_repositories.GetTeamByIdentityParams{
		TeamIdentity:         input.TeamIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, core.NewNotFoundError("team not found")
	}

	return team_core.TeamToDto(team), nil
}
