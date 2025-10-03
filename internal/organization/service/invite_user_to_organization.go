package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	"github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type InviteUserToOrganizationService struct {
	OrganizationRepository     organizationrepo.OrganizationRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	UserRepository             userrepo.UserRepository
	RoleRepository             rolerepo.RoleRepository
	WorkspaceRepository        workspacerepo.WorkspaceRepository
	WorkspaceUserRepository    workspacerepo.WorkspaceUserRepository
	ProjectRepository          projectrepo.ProjectRepository
	ProjectUserRepository      projectrepo.ProjectUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewInviteUserToOrganizationService(
	organizationRepository organizationrepo.OrganizationRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	userRepository userrepo.UserRepository,
	roleRepository rolerepo.RoleRepository,
	workspaceRepository workspacerepo.WorkspaceRepository,
	workspaceUserRepository workspacerepo.WorkspaceUserRepository,
	projectRepository projectrepo.ProjectRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
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

func (s *InviteUserToOrganizationService) createWorkspaceUsers(organization *organization.Organization, user *user.User, workspaces []InviteUserToOrganizationWorkspaceInput) error {
	for _, w := range workspaces {
		wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{
			WorkspaceIdentity:    w.WorkspaceIdentity,
			OrganizationIdentity: &organization.Identity,
		})
		if err != nil {
			return err
		}

		if wrk == nil {
			return core.NewNotFoundError("workspace not found")
		}

		workspaceUser, err := workspace.NewWorkspaceUser(workspace.NewWorkspaceUserInput{
			WorkspaceIdentity: wrk.Identity,
			User:              *user,
		})
		if err != nil {
			return err
		}

		_, err = s.WorkspaceUserRepository.StoreWorkspaceUser(workspacerepo.StoreWorkspaceUserParams{WorkspaceUser: workspaceUser})
		if err != nil {
			return err
		}

		for _, p := range w.Projects {
			prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
				ProjectIdentity:   p,
				WorkspaceIdentity: &wrk.Identity,
			})
			if err != nil {
				return err
			}

			if prj == nil {
				return core.NewNotFoundError("project not found")
			}

			projectUser, err := project.NewProjectUser(project.NewProjectUserInput{
				ProjectIdentity: prj.Identity,
				User:            *user,
			})
			if err != nil {
				return err
			}

			_, err = s.ProjectUserRepository.StoreProjectUser(projectrepo.StoreProjectUserParams{ProjectUser: projectUser})
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

	org, err := s.OrganizationRepository.GetOrganizationByIdentity(organizationrepo.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		return err
	}

	if org == nil {
		return core.NewNotFoundError("organization not found")
	}

	user, err := s.UserRepository.GetUserByEmail(userrepo.GetUserByEmailParams{
		Email: input.Email,
	})
	if err != nil {
		return err
	}

	if user == nil {
		return core.NewNotFoundError("user not found")
	}

	role, err := s.RoleRepository.GetRoleByIdentity(rolerepo.GetRoleByIdentityParams{RoleIdentity: input.RoleIdentity})
	if err != nil {
		return err
	}

	if role == nil {
		return core.NewNotFoundError("role not found")
	}

	var organizationUser *organization.OrganizationUser = nil
	organizationUser, err = s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         user.Identity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		organizationUser, err = organization.NewOrganizationUser(organization.NewOrganizationUserInput{
			OrganizationIdentity: input.OrganizationIdentity,
			User:                 *user,
			Role:                 *role,
			Status:               organization.OrganizationUserStatusInvited,
		})
		if err != nil {
			return err
		}

		organizationUser, err = s.OrganizationUserRepository.StoreOrganizationUser(organizationrepo.StoreOrganizationUserParams{OrganizationUser: organizationUser})
		if err != nil {
			return err
		}

		err = s.createWorkspaceUsers(org, user, input.Workspaces)
		if err != nil {
			return err
		}
	} else {
		s.WorkspaceUserRepository.DeleteAllByUserIdentity(workspacerepo.DeleteAllByUserIdentityParams{UserIdentity: user.Identity})
		s.ProjectUserRepository.DeleteAllByUserIdentity(projectrepo.DeleteAllByUserIdentityParams{UserIdentity: user.Identity})
		err = s.createWorkspaceUsers(org, user, input.Workspaces)
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
