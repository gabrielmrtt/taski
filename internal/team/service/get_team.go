package teamservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
)

type GetTeamService struct {
	TeamRepository teamrepo.TeamRepository
}

func NewGetTeamService(teamRepository teamrepo.TeamRepository) *GetTeamService {
	return &GetTeamService{
		TeamRepository: teamRepository,
	}
}

type GetTeamInput struct {
	TeamIdentity         core.Identity
	OrganizationIdentity core.Identity
	RelationsInput       core.RelationsInput
}

func (i GetTeamInput) Validate() error {
	return nil
}

func (s *GetTeamService) Execute(input GetTeamInput) (*team.TeamDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tm, err := s.TeamRepository.GetTeamByIdentity(teamrepo.GetTeamByIdentityParams{
		TeamIdentity:         input.TeamIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
		RelationsInput:       input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	if tm == nil {
		return nil, core.NewNotFoundError("team not found")
	}

	return team.TeamToDto(tm), nil
}
