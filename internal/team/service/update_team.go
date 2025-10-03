package teamservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
)

type UpdateTeamService struct {
	TeamRepository             teamrepo.TeamRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewUpdateTeamService(
	teamRepository teamrepo.TeamRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
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
	Status               *team.TeamStatuses
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
		teamStatus := team.TeamStatuses(*i.Status)
		if teamStatus != team.TeamStatusActive && teamStatus != team.TeamStatusInactive {
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

	tm, err := s.TeamRepository.GetTeamByIdentity(teamrepo.GetTeamByIdentityParams{
		TeamIdentity:         input.TeamIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if tm == nil {
		tx.Rollback()
		return core.NewNotFoundError("team not found")
	}

	if input.Name != nil {
		err = tm.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Description != nil {
		err = tm.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Status != nil {
		if *input.Status == team.TeamStatusActive {
			tm.Activate()
		} else {
			tm.Inactivate()
		}
	}

	if input.Members != nil {
		tm.RemoveAllUsers()
		for _, usrIdentity := range input.Members {
			usr, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
				OrganizationIdentity: input.OrganizationIdentity,
				UserIdentity:         usrIdentity,
			})
			if err != nil {
				tx.Rollback()
				return err
			}

			if usr == nil {
				tx.Rollback()
				return core.NewNotFoundError("user not found")
			}

			tm.AddUser(usr.User)
		}
	}

	err = s.TeamRepository.UpdateTeam(teamrepo.UpdateTeamParams{Team: tm})
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
