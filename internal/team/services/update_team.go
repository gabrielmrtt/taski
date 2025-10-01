package team_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
)

type UpdateTeamService struct {
	TeamRepository             team_repositories.TeamRepository
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewUpdateTeamService(
	teamRepository team_repositories.TeamRepository,
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateTeamService {
	return &UpdateTeamService{
		TeamRepository:             teamRepository,
		OrganizationUserRepository: organizationUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type UpdateTeamInput struct {
	TeamIdentity         core.Identity
	OrganizationIdentity core.Identity
	UserEditorIdentity   core.Identity
	Name                 *string
	Description          *string
	Status               *team_core.TeamStatuses
	Members              []core.Identity
}

func (i UpdateTeamInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if i.Description != nil {
		_, err := core.NewDescription(*i.Description)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "description",
				Error: err.Error(),
			})
		}
	}

	if i.Status != nil {
		teamStatus := team_core.TeamStatuses(*i.Status)
		if teamStatus != team_core.TeamStatusActive && teamStatus != team_core.TeamStatusInactive {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "status",
				Error: "invalid status. valid statuses are: active, inactive",
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateTeamService) Execute(input UpdateTeamInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TeamRepository.SetTransaction(tx)
	s.OrganizationUserRepository.SetTransaction(tx)

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

	if input.Name != nil {
		err = team.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Description != nil {
		err = team.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Status != nil {
		if *input.Status == team_core.TeamStatusActive {
			team.Activate()
		} else {
			team.Inactivate()
		}
	}

	if input.Members != nil {
		team.RemoveAllUsers()
		for _, user := range input.Members {
			user, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
				OrganizationIdentity: input.OrganizationIdentity,
				UserIdentity:         user,
			})
			if err != nil {
				tx.Rollback()
				return err
			}

			if user == nil {
				tx.Rollback()
				return core.NewNotFoundError("user not found")
			}

			team.AddUser(user.User)
		}
	}

	err = s.TeamRepository.UpdateTeam(team_repositories.UpdateTeamParams{Team: team})
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
