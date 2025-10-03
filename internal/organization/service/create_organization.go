package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type CreateOrganizationService struct {
	OrganizationRepository     organizationrepo.OrganizationRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	RoleRepository             rolerepo.RoleRepository
	UserRepository             userrepo.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewCreateOrganizationService(
	organizationRepository organizationrepo.OrganizationRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	roleRepository rolerepo.RoleRepository,
	userRepository userrepo.UserRepository,
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

func (s *CreateOrganizationService) Execute(input CreateOrganizationInput) (*organization.OrganizationDto, error) {
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

	org, err := organization.NewOrganization(organization.NewOrganizationInput{
		Name:                input.Name,
		UserCreatorIdentity: &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	adminRole, err := s.RoleRepository.GetSystemDefaultRole(rolerepo.GetDefaultRoleParams{
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

	user, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{
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

	organizationUser, err := organization.NewOrganizationUser(organization.NewOrganizationUserInput{
		OrganizationIdentity: org.Identity,
		User:                 *user,
		Role:                 *adminRole,
		Status:               organization.OrganizationUserStatusActive,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	org, err = s.OrganizationRepository.StoreOrganization(organizationrepo.StoreOrganizationParams{Organization: org})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = s.OrganizationUserRepository.StoreOrganizationUser(organizationrepo.StoreOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return organization.OrganizationToDto(org), nil
}
