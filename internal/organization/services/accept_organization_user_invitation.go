package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type AcceptOrganizationUserInvitationService struct {
	OrganizationRepository organization_core.OrganizationRepository
	UserRepository         user_core.UserRepository
	TransactionRepository  core.TransactionRepository
}

func NewAcceptOrganizationUserInvitationService(
	organizationRepository organization_core.OrganizationRepository,
	userRepository user_core.UserRepository,
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
