package project_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
)

type DeleteProjectService struct {
	ProjectRepository     project_repositories.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteProjectService(
	projectRepository project_repositories.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectService {
	return &DeleteProjectService{
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteProjectInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
	ProjectIdentity      core.Identity
}

func (i DeleteProjectInput) Validate() error {
	return nil
}

func (s *DeleteProjectService) Execute(input DeleteProjectInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectRepository.SetTransaction(tx)

	project, err := s.ProjectRepository.GetProjectByIdentity(project_repositories.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		WorkspaceIdentity:    &input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if project == nil {
		return core.NewNotFoundError("project not found")
	}

	project.Delete()

	err = s.ProjectRepository.DeleteProject(project_repositories.DeleteProjectParams{ProjectIdentity: input.ProjectIdentity})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
