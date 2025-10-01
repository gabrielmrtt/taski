package team_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
)

type DeleteTeamService struct {
	TeamRepository        team_repositories.TeamRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteTeamService(
	teamRepository team_repositories.TeamRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTeamService {
	return &DeleteTeamService{
		TeamRepository:        teamRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteTeamInput struct {
	TeamIdentity         core.Identity
	OrganizationIdentity core.Identity
}

func (i DeleteTeamInput) Validate() error {
	return nil
}

func (s *DeleteTeamService) Execute(input DeleteTeamInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TeamRepository.SetTransaction(tx)

	team, err := s.TeamRepository.GetTeamByIdentity(team_repositories.GetTeamByIdentityParams{
		TeamIdentity:         input.TeamIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if team == nil {
		tx.Rollback()
		return core.NewNotFoundError("team not found")
	}

	err = s.TeamRepository.DeleteTeam(team_repositories.DeleteTeamParams{
		TeamIdentity: team.Identity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
