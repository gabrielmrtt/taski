package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type UpdateOrganizationUserService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	RoleRepository             role_repositories.RoleRepository
	WorkspaceRepository        workspace_repositories.WorkspaceRepository
	WorkspaceUserRepository    workspace_repositories.WorkspaceUserRepository
	ProjectRepository          project_repositories.ProjectRepository
	ProjectUserRepository      project_repositories.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewUpdateOrganizationUserService(
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	roleRepository role_repositories.RoleRepository,
	workspaceRepository workspace_repositories.WorkspaceRepository,
	projectRepository project_repositories.ProjectRepository,
	workspaceUserRepository workspace_repositories.WorkspaceUserRepository,
	projectUserRepository project_repositories.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateOrganizationUserService {
	return &UpdateOrganizationUserService{
		OrganizationUserRepository: organizationUserRepository,
		RoleRepository:             roleRepository,
		WorkspaceRepository:        workspaceRepository,
		ProjectRepository:          projectRepository,
		WorkspaceUserRepository:    workspaceUserRepository,
		ProjectUserRepository:      projectUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type UpdateOrganizationUserWorkspaceInput struct {
	WorkspaceIdentity core.Identity
	Projects          []core.Identity
}

type UpdateOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
	RoleIdentity         *core.Identity
	Status               *organization_core.OrganizationUserStatuses
	Workspaces           []UpdateOrganizationUserWorkspaceInput
}

func (i UpdateOrganizationUserInput) Validate() error {
	return nil
}

func (s *UpdateOrganizationUserService) Execute(input UpdateOrganizationUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.OrganizationUserRepository.SetTransaction(tx)
	s.RoleRepository.SetTransaction(tx)

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	if input.RoleIdentity != nil {
		role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(role_repositories.GetRoleByIdentityAndOrganizationIdentityParams{
			RoleIdentity:         *input.RoleIdentity,
			OrganizationIdentity: input.OrganizationIdentity,
		})
		if err != nil {
			return core.NewInternalError(err.Error())
		}

		if role == nil {
			return core.NewNotFoundError("role not found")
		}

		organizationUser.ChangeRole(role)
	}

	if input.Status != nil {
		if *input.Status == organization_core.OrganizationUserStatusActive && organizationUser.IsInactive() {
			organizationUser.Activate()
		} else if *input.Status == organization_core.OrganizationUserStatusInactive && organizationUser.IsActive() {
			organizationUser.Deactivate()
		} else {
			return core.NewInvalidInputError("invalid input", []core.InvalidInputErrorField{
				{
					Field: "status",
					Error: "valid statuses are: active, inactive",
				},
			})
		}
	}

	if input.Workspaces != nil {
		s.WorkspaceUserRepository.DeleteAllByUserIdentity(workspace_repositories.DeleteAllByUserIdentityParams{UserIdentity: input.UserIdentity})
		s.ProjectUserRepository.DeleteAllByUserIdentity(project_repositories.DeleteAllByUserIdentityParams{UserIdentity: input.UserIdentity})

		for _, w := range input.Workspaces {
			workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{
				WorkspaceIdentity:    w.WorkspaceIdentity,
				OrganizationIdentity: &input.OrganizationIdentity,
			})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			if workspace == nil {
				return core.NewNotFoundError("workspace not found")
			}

			workspaceUser, err := workspace_core.NewWorkspaceUser(workspace_core.NewWorkspaceUserInput{
				WorkspaceIdentity: workspace.Identity,
				UserIdentity:      input.UserIdentity,
				Status:            workspace_core.WorkspaceUserStatuses(organizationUser.Status),
			})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			_, err = s.WorkspaceUserRepository.StoreWorkspaceUser(workspace_repositories.StoreWorkspaceUserParams{WorkspaceUser: workspaceUser})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			for _, p := range w.Projects {
				project, err := s.ProjectRepository.GetProjectByIdentity(project_repositories.GetProjectByIdentityParams{
					ProjectIdentity:   p,
					WorkspaceIdentity: &workspace.Identity,
				})
				if err != nil {
					return core.NewInternalError(err.Error())
				}

				if project == nil {
					return core.NewNotFoundError("project not found")
				}

				projectUser, err := project_core.NewProjectUser(project_core.NewProjectUserInput{
					ProjectIdentity: project.Identity,
					UserIdentity:    input.UserIdentity,
					Status:          project_core.ProjectUserStatuses(organizationUser.Status),
				})
				if err != nil {
					return core.NewInternalError(err.Error())
				}

				_, err = s.ProjectUserRepository.StoreProjectUser(project_repositories.StoreProjectUserParams{ProjectUser: projectUser})
				if err != nil {
					return core.NewInternalError(err.Error())
				}
			}
		}
	}

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organization_repositories.UpdateOrganizationUserParams{
		OrganizationUser: organizationUser,
	})
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	return nil
}
