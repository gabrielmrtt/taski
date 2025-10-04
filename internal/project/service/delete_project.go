package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type DeleteProjectService struct {
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteProjectService(
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectService {
	return &DeleteProjectService{
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteProjectInput struct {
	OrganizationIdentity core.Identity
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

	project, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if project == nil {
		return core.NewNotFoundError("project not found")
	}

	project.Delete()

	err = s.ProjectRepository.DeleteProject(projectrepo.DeleteProjectParams{ProjectIdentity: input.ProjectIdentity})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
