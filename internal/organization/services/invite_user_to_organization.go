package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type InviteUserToOrganizationService struct {
	OrganizationRepository     organization_repositories.OrganizationRepository
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	UserRepository             user_repositories.UserRepository
	RoleRepository             role_repositories.RoleRepository
	WorkspaceRepository        workspace_repositories.WorkspaceRepository
	WorkspaceUserRepository    workspace_repositories.WorkspaceUserRepository
	ProjectRepository          project_repositories.ProjectRepository
	ProjectUserRepository      project_repositories.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewInviteUserToOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	userRepository user_repositories.UserRepository,
	roleRepository role_repositories.RoleRepository,
	workspaceRepository workspace_repositories.WorkspaceRepository,
	workspaceUserRepository workspace_repositories.WorkspaceUserRepository,
	projectRepository project_repositories.ProjectRepository,
	projectUserRepository project_repositories.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *InviteUserToOrganizationService {
	return &InviteUserToOrganizationService{
		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		UserRepository:             userRepository,
		RoleRepository:             roleRepository,
		WorkspaceRepository:        workspaceRepository,
		WorkspaceUserRepository:    workspaceUserRepository,
		ProjectRepository:          projectRepository,
		ProjectUserRepository:      projectUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type InviteUserToOrganizationWorkspaceInput struct {
	WorkspaceIdentity core.Identity
	Projects          []core.Identity
}

type InviteUserToOrganizationInput struct {
	OrganizationIdentity core.Identity
	Email                string
	RoleIdentity         core.Identity
	Workspaces           []InviteUserToOrganizationWorkspaceInput
}

func (i InviteUserToOrganizationInput) Validate() error {
	return nil
}

func (s *InviteUserToOrganizationService) createWorkspaceUsers(organization *organization_core.Organization, user *user_core.User, workspaces []InviteUserToOrganizationWorkspaceInput) error {
	for _, w := range workspaces {
		workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{
			WorkspaceIdentity:    w.WorkspaceIdentity,
			OrganizationIdentity: &organization.Identity,
		})
		if err != nil {
			return err
		}

		if workspace == nil {
			return core.NewNotFoundError("workspace not found")
		}

		workspaceUser, err := workspace_core.NewWorkspaceUser(workspace_core.NewWorkspaceUserInput{
			WorkspaceIdentity: workspace.Identity,
			User:              *user,
		})
		if err != nil {
			return err
		}

		_, err = s.WorkspaceUserRepository.StoreWorkspaceUser(workspace_repositories.StoreWorkspaceUserParams{WorkspaceUser: workspaceUser})
		if err != nil {
			return err
		}

		for _, p := range w.Projects {
			project, err := s.ProjectRepository.GetProjectByIdentity(project_repositories.GetProjectByIdentityParams{
				ProjectIdentity:   p,
				WorkspaceIdentity: &workspace.Identity,
			})
			if err != nil {
				return err
			}

			if project == nil {
				return core.NewNotFoundError("project not found")
			}

			projectUser, err := project_core.NewProjectUser(project_core.NewProjectUserInput{
				ProjectIdentity: project.Identity,
				User:            *user,
			})
			if err != nil {
				return err
			}

			_, err = s.ProjectUserRepository.StoreProjectUser(project_repositories.StoreProjectUserParams{ProjectUser: projectUser})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *InviteUserToOrganizationService) Execute(input InviteUserToOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)
	s.WorkspaceRepository.SetTransaction(tx)
	s.WorkspaceUserRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)
	s.ProjectUserRepository.SetTransaction(tx)

	organization, err := s.OrganizationRepository.GetOrganizationByIdentity(organization_repositories.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		return err
	}

	if organization == nil {
		return core.NewNotFoundError("organization not found")
	}

	user, err := s.UserRepository.GetUserByEmail(user_repositories.GetUserByEmailParams{
		Email: input.Email,
	})
	if err != nil {
		return err
	}

	if user == nil {
		return core.NewNotFoundError("user not found")
	}

	role, err := s.RoleRepository.GetRoleByIdentity(role_repositories.GetRoleByIdentityParams{RoleIdentity: input.RoleIdentity})
	if err != nil {
		return err
	}

	if role == nil {
		return core.NewNotFoundError("role not found")
	}

	var organizationUser *organization_core.OrganizationUser = nil
	organizationUser, err = s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         user.Identity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		organizationUser, err = organization_core.NewOrganizationUser(organization_core.NewOrganizationUserInput{
			OrganizationIdentity: input.OrganizationIdentity,
			User:                 *user,
			Role:                 *role,
			Status:               organization_core.OrganizationUserStatusInvited,
		})
		if err != nil {
			return err
		}

		organizationUser, err = s.OrganizationUserRepository.StoreOrganizationUser(organization_repositories.StoreOrganizationUserParams{OrganizationUser: organizationUser})
		if err != nil {
			return err
		}

		err = s.createWorkspaceUsers(organization, user, input.Workspaces)
		if err != nil {
			return err
		}
	} else {
		s.WorkspaceUserRepository.DeleteAllByUserIdentity(workspace_repositories.DeleteAllByUserIdentityParams{UserIdentity: user.Identity})
		s.ProjectUserRepository.DeleteAllByUserIdentity(project_repositories.DeleteAllByUserIdentityParams{UserIdentity: user.Identity})
		err = s.createWorkspaceUsers(organization, user, input.Workspaces)
		if err != nil {
			return err
		}
	}

	organizationUser.Invite()

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
