package teamservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
	"github.com/gabrielmrtt/taski/internal/user"
)

type CreateTeamService struct {
	TeamRepository             teamrepo.TeamRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewCreateTeamService(
	teamRepository teamrepo.TeamRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	transactionRepository core.TransactionRepository,
) *CreateTeamService {
	return &CreateTeamService{
		TeamRepository:             teamRepository,
		OrganizationUserRepository: organizationUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type CreateTeamInput struct {
	Name                 string
	Description          string
	OrganizationIdentity core.Identity
	UserCreatorIdentity  core.Identity
	Members              []core.Identity
}

func (i CreateTeamInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewName(i.Name); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if _, err := core.NewDescription(i.Description); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "description",
			Error: err.Error(),
		})
	}

	if len(i.Members) == 0 {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "members",
			Error: "members cannot be empty",
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateTeamService) Execute(input CreateTeamInput) (*team.TeamDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TeamRepository.SetTransaction(tx)
	s.OrganizationUserRepository.SetTransaction(tx)

	var users []user.User = make([]user.User, 0)
	for _, usrIdentity := range input.Members {
		usr, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
			UserIdentity:         usrIdentity,
			OrganizationIdentity: input.OrganizationIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if usr == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("user not found")
		}

		users = append(users, usr.User)
	}

	tm, err := team.NewTeam(team.NewTeamInput{
		Name:                 input.Name,
		Description:          input.Description,
		OrganizationIdentity: input.OrganizationIdentity,
		UserCreatorIdentity:  &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, user := range users {
		tm.AddUser(user)
	}

	_, err = s.TeamRepository.StoreTeam(teamrepo.StoreTeamParams{Team: tm})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return team.TeamToDto(tm), nil
}
