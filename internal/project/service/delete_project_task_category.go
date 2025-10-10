package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type DeleteProjectTaskCategoryService struct {
	ProjectRepository             projectrepo.ProjectRepository
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
	TransactionRepository         core.TransactionRepository
}

func NewDeleteProjectTaskCategoryService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectTaskCategoryService {
	return &DeleteProjectTaskCategoryService{
		ProjectRepository:             projectRepository,
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
		TransactionRepository:         transactionRepository,
	}
}

type DeleteProjectTaskCategoryInput struct {
	OrganizationIdentity        core.Identity
	ProjectIdentity             core.Identity
	ProjectTaskCategoryIdentity core.Identity
}

func (i DeleteProjectTaskCategoryInput) Validate() error {
	return nil
}

func (s *DeleteProjectTaskCategoryService) Execute(input DeleteProjectTaskCategoryInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectTaskCategoryRepository.SetTransaction(tx)

	projectTaskCategory, err := s.ProjectTaskCategoryRepository.GetProjectTaskCategoryByIdentity(projectrepo.GetProjectTaskCategoryByIdentityParams{
		ProjectTaskCategoryIdentity: &input.ProjectTaskCategoryIdentity,
		ProjectIdentity:             &input.ProjectIdentity,
	})
	if err != nil {
		return err
	}

	if projectTaskCategory == nil || projectTaskCategory.IsDeleted() {
		return core.NewNotFoundError("project task category not found")
	}

	projectTaskCategory.Delete()

	err = s.ProjectTaskCategoryRepository.UpdateProjectTaskCategory(projectrepo.UpdateProjectTaskCategoryParams{ProjectTaskCategory: projectTaskCategory})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
