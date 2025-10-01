package team_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type CreateTeamService struct {
	TeamRepository             team_repositories.TeamRepository
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewCreateTeamService(
	teamRepository team_repositories.TeamRepository,
	organizationUserRepository organization_repositories.OrganizationUserRepository,
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

func (s *CreateTeamService) Execute(input CreateTeamInput) (*team_core.TeamDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TeamRepository.SetTransaction(tx)
	s.OrganizationUserRepository.SetTransaction(tx)

	var users []user_core.User = make([]user_core.User, 0)
	for _, user := range input.Members {
		user, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
			UserIdentity:         user,
			OrganizationIdentity: input.OrganizationIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if user == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("user not found")
		}

		users = append(users, user.User)
	}

	team, err := team_core.NewTeam(team_core.NewTeamInput{
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
		team.AddUser(user)
	}

	_, err = s.TeamRepository.StoreTeam(team_repositories.StoreTeamParams{Team: team})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return team_core.TeamToDto(team), nil
}
