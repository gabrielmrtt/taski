package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type UpdateOrganizationUserService struct {
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	OrganizationRepository     organizationrepo.OrganizationRepository
	RoleRepository             rolerepo.RoleRepository
	WorkspaceRepository        workspacerepo.WorkspaceRepository
	WorkspaceUserRepository    workspacerepo.WorkspaceUserRepository
	ProjectRepository          projectrepo.ProjectRepository
	ProjectUserRepository      projectrepo.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewUpdateOrganizationUserService(
	organizationRepository organizationrepo.OrganizationRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	roleRepository rolerepo.RoleRepository,
	workspaceRepository workspacerepo.WorkspaceRepository,
	projectRepository projectrepo.ProjectRepository,
	workspaceUserRepository workspacerepo.WorkspaceUserRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateOrganizationUserService {
	return &UpdateOrganizationUserService{
		OrganizationUserRepository: organizationUserRepository,
		OrganizationRepository:     organizationRepository,
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
	Status               *organization.OrganizationUserStatuses
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
	s.OrganizationRepository.SetTransaction(tx)
	s.RoleRepository.SetTransaction(tx)

	org, err := s.OrganizationRepository.GetOrganizationByIdentity(organizationrepo.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	if org == nil {
		return core.NewNotFoundError("organization not found")
	}

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
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
		role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(rolerepo.GetRoleByIdentityAndOrganizationIdentityParams{
			RoleIdentity:         *input.RoleIdentity,
			OrganizationIdentity: input.OrganizationIdentity,
		})
		if err != nil {
			return core.NewInternalError(err.Error())
		}

		if role == nil {
			return core.NewNotFoundError("role not found")
		}

		organizationUser.ChangeRole(*role)
	}

	if input.Status != nil {
		if *input.Status == organization.OrganizationUserStatusActive && organizationUser.IsInactive() {
			organizationUser.Activate()
		} else if *input.Status == organization.OrganizationUserStatusInactive && organizationUser.IsActive() {
			if organizationUser.User.Identity.Equals(*org.UserCreatorIdentity) {
				return core.NewConflictError("cannot deactivate the creator of the organization")
			}

			organizationUser.Deactivate()
		} else {
			return core.NewInvalidInputError("invalid input", []core.InvalidInputErrorField{
				{
					Field: "status",
					Error: "valid statuses are: active, inactive",
				},
			})
		}

		workspaceUsers, err := s.WorkspaceUserRepository.GetWorkspaceUsersByUserIdentity(workspacerepo.GetWorkspaceUsersByUserIdentityParams{
			UserIdentity: input.UserIdentity,
		})
		if err != nil {
			return core.NewInternalError(err.Error())
		}

		for _, workspaceUser := range workspaceUsers {
			workspaceUser.Status = workspace.WorkspaceUserStatuses(organizationUser.Status)
			err = s.WorkspaceUserRepository.UpdateWorkspaceUser(workspacerepo.UpdateWorkspaceUserParams{WorkspaceUser: &workspaceUser})
			if err != nil {
				return core.NewInternalError(err.Error())
			}
		}

		projectUsers, err := s.ProjectUserRepository.GetProjectUsersByUserIdentity(projectrepo.GetProjectUsersByUserIdentityParams{
			UserIdentity: input.UserIdentity,
		})
		if err != nil {
			return core.NewInternalError(err.Error())
		}

		for _, projectUser := range projectUsers {
			projectUser.Status = project.ProjectUserStatuses(organizationUser.Status)
			err = s.ProjectUserRepository.UpdateProjectUser(projectrepo.UpdateProjectUserParams{ProjectUser: &projectUser})
			if err != nil {
				return core.NewInternalError(err.Error())
			}
		}
	}

	if input.Workspaces != nil {
		s.WorkspaceUserRepository.DeleteAllByUserIdentity(workspacerepo.DeleteAllByUserIdentityParams{UserIdentity: input.UserIdentity})
		s.ProjectUserRepository.DeleteAllByUserIdentity(projectrepo.DeleteAllByUserIdentityParams{UserIdentity: input.UserIdentity})

		for _, w := range input.Workspaces {
			wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{
				WorkspaceIdentity:    w.WorkspaceIdentity,
				OrganizationIdentity: &input.OrganizationIdentity,
			})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			if wrk == nil {
				return core.NewNotFoundError("workspace not found")
			}

			workspaceUser, err := workspace.NewWorkspaceUser(workspace.NewWorkspaceUserInput{
				WorkspaceIdentity: wrk.Identity,
				User:              organizationUser.User,
				Status:            workspace.WorkspaceUserStatuses(organizationUser.Status),
			})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			_, err = s.WorkspaceUserRepository.StoreWorkspaceUser(workspacerepo.StoreWorkspaceUserParams{WorkspaceUser: workspaceUser})
			if err != nil {
				return core.NewInternalError(err.Error())
			}

			for _, p := range w.Projects {
				prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
					ProjectIdentity:   p,
					WorkspaceIdentity: &wrk.Identity,
				})
				if err != nil {
					return core.NewInternalError(err.Error())
				}

				if prj == nil {
					return core.NewNotFoundError("project not found")
				}

				projectUser, err := project.NewProjectUser(project.NewProjectUserInput{
					ProjectIdentity: prj.Identity,
					User:            organizationUser.User,
					Status:          project.ProjectUserStatuses(organizationUser.Status),
				})
				if err != nil {
					return core.NewInternalError(err.Error())
				}

				_, err = s.ProjectUserRepository.StoreProjectUser(projectrepo.StoreProjectUserParams{ProjectUser: projectUser})
				if err != nil {
					return core.NewInternalError(err.Error())
				}
			}
		}
	}

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organizationrepo.UpdateOrganizationUserParams{
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
