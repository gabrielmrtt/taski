package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type AcceptOrganizationUserInvitationService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
	UserRepository         user_repositories.UserRepository
	TransactionRepository  core.TransactionRepository
}

func NewAcceptOrganizationUserInvitationService(
	organizationRepository organization_repositories.OrganizationRepository,
	userRepository user_repositories.UserRepository,
	transactionRepository core.TransactionRepository,
) *AcceptOrganizationUserInvitationService {
	return &AcceptOrganizationUserInvitationService{
		OrganizationRepository: organizationRepository,
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
	}
}

type AcceptOrganizationUserInvitationInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i AcceptOrganizationUserInvitationInput) Validate() error {
	return nil
}

func (s *AcceptOrganizationUserInvitationService) Execute(input AcceptOrganizationUserInvitationInput) error {
	return nil
}
