package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type AcceptOrganizationUserInvitationService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	WorkspaceUserRepository    workspace_repositories.WorkspaceUserRepository
	ProjectUserRepository      project_repositories.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewAcceptOrganizationUserInvitationService(
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	workspaceUserRepository workspace_repositories.WorkspaceUserRepository,
	projectUserRepository project_repositories.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *AcceptOrganizationUserInvitationService {
	return &AcceptOrganizationUserInvitationService{
		OrganizationUserRepository: organizationUserRepository,
		WorkspaceUserRepository:    workspaceUserRepository,
		ProjectUserRepository:      projectUserRepository,
		TransactionRepository:      transactionRepository,
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

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	organizationUser.AcceptInvitation()

	workspaceUsers, err := s.WorkspaceUserRepository.GetWorkspaceUsersByUserIdentity(workspace_repositories.GetWorkspaceUsersByUserIdentityParams{
		UserIdentity: input.UserIdentity,
	})
	if err != nil {
		return err
	}

	for _, workspaceUser := range workspaceUsers {
		workspaceUser.AcceptInvitation()
		err = s.WorkspaceUserRepository.UpdateWorkspaceUser(workspace_repositories.UpdateWorkspaceUserParams{WorkspaceUser: &workspaceUser})
		if err != nil {
			return err
		}
	}

	projectUsers, err := s.ProjectUserRepository.GetProjectUsersByUserIdentity(project_repositories.GetProjectUsersByUserIdentityParams{
		UserIdentity: input.UserIdentity,
	})
	if err != nil {
		return err
	}

	for _, projectUser := range projectUsers {
		projectUser.AcceptInvitation()
		err = s.ProjectUserRepository.UpdateProjectUser(project_repositories.UpdateProjectUserParams{ProjectUser: &projectUser})
		if err != nil {
			return err
		}
	}

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organization_repositories.UpdateOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
