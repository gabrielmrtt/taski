package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type CreateOrganizationService struct {
	OrganizationRepository organization_core.OrganizationRepository
	RoleRepository         role_core.RoleRepository
	UserRepository         user_core.UserRepository
	TransactionRepository  core.TransactionRepository
}

func NewCreateOrganizationService(
	organizationRepository organization_core.OrganizationRepository,
	roleRepository role_core.RoleRepository,
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *CreateOrganizationService {
	return &CreateOrganizationService{
		OrganizationRepository: organizationRepository,
		RoleRepository:         roleRepository,
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
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

	organization, err := organization_core.NewOrganization(organization_core.NewOrganizationInput{
		Name:                input.Name,
		UserCreatorIdentity: &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	adminRole, err := s.RoleRepository.GetSystemDefaultRole(role_core.GetDefaultRoleParams{
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

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: input.UserCreatorIdentity,
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

	organization, err = s.OrganizationRepository.StoreOrganization(organization)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	organizationUser, err = s.OrganizationRepository.CreateOrganizationUser(organizationUser)
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
