package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type RefuseOrganizationUserInvitationService struct {
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	WorkspaceUserRepository    workspacerepo.WorkspaceUserRepository
	ProjectUserRepository      projectrepo.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewRefuseOrganizationUserInvitationService(
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	workspaceUserRepository workspacerepo.WorkspaceUserRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *RefuseOrganizationUserInvitationService {
	return &RefuseOrganizationUserInvitationService{
		OrganizationUserRepository: organizationUserRepository,
		WorkspaceUserRepository:    workspaceUserRepository,
		ProjectUserRepository:      projectUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type RefuseOrganizationUserInvitationInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i RefuseOrganizationUserInvitationInput) Validate() error {
	return nil
}

func (s *RefuseOrganizationUserInvitationService) Execute(input RefuseOrganizationUserInvitationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationUserRepository.SetTransaction(tx)
	s.WorkspaceUserRepository.SetTransaction(tx)
	s.ProjectUserRepository.SetTransaction(tx)

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	organizationUser.RefuseInvitation()

	workspaceUsers, err := s.WorkspaceUserRepository.GetWorkspaceUsersByUserIdentity(workspacerepo.GetWorkspaceUsersByUserIdentityParams{
		UserIdentity: input.UserIdentity,
	})
	if err != nil {
		return err
	}

	for _, workspaceUser := range workspaceUsers {
		workspaceUser.RefuseInvitation()
		err = s.WorkspaceUserRepository.UpdateWorkspaceUser(workspacerepo.UpdateWorkspaceUserParams{WorkspaceUser: &workspaceUser})
		if err != nil {
			return err
		}
	}

	projectUsers, err := s.ProjectUserRepository.GetProjectUsersByUserIdentity(projectrepo.GetProjectUsersByUserIdentityParams{
		UserIdentity: input.UserIdentity,
	})
	if err != nil {
		return err
	}

	for _, projectUser := range projectUsers {
		projectUser.RefuseInvitation()
		err = s.ProjectUserRepository.UpdateProjectUser(projectrepo.UpdateProjectUserParams{ProjectUser: &projectUser})
		if err != nil {
			return err
		}
	}

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organizationrepo.UpdateOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
