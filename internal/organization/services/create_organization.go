package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type CreateOrganizationService struct {
	OrganizationRepository     organization_repositories.OrganizationRepository
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	RoleRepository             role_repositories.RoleRepository
	UserRepository             user_repositories.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewCreateOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	roleRepository role_repositories.RoleRepository,
	userRepository user_repositories.UserRepository,
	transactionRepository core.TransactionRepository,
) *CreateOrganizationService {
	return &CreateOrganizationService{
		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		RoleRepository:             roleRepository,
		UserRepository:             userRepository,
		TransactionRepository:      transactionRepository,
	}
}

type CreateOrganizationInput struct {
	Name                string
	UserCreatorIdentity core.Identity
}

func (i CreateOrganizationInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := core.NewName(i.Name)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateOrganizationService) Execute(input CreateOrganizationInput) (*organization_core.OrganizationDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.OrganizationRepository.SetTransaction(tx)
	s.OrganizationUserRepository.SetTransaction(tx)
	s.RoleRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)

	organization, err := organization_core.NewOrganization(organization_core.NewOrganizationInput{
		Name:                input.Name,
		UserCreatorIdentity: &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	adminRole, err := s.RoleRepository.GetSystemDefaultRole(role_repositories.GetDefaultRoleParams{
		Slug: "admin",
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if adminRole == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("admin role not found")
	}

	user, err := s.UserRepository.GetUserByIdentity(user_repositories.GetUserByIdentityParams{
		UserIdentity: input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if user == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("user not found")
	}

	organizationUser, err := organization_core.NewOrganizationUser(organization_core.NewOrganizationUserInput{
		OrganizationIdentity: organization.Identity,
		User:                 user,
		Role:                 adminRole,
		Status:               organization_core.OrganizationUserStatusActive,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	organization, err = s.OrganizationRepository.StoreOrganization(organization_repositories.StoreOrganizationParams{Organization: organization})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	organizationUser, err = s.OrganizationUserRepository.StoreOrganizationUser(organization_repositories.StoreOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return organization_core.OrganizationToDto(organization), nil
}
