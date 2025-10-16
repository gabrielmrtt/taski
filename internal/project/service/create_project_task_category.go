package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type CreateProjectTaskCategoryService struct {
	ProjectRepository             projectrepo.ProjectRepository
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
	TransactionRepository         core.TransactionRepository
}

func NewCreateProjectTaskCategoryService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectTaskCategoryService {
	return &CreateProjectTaskCategoryService{
		ProjectRepository:             projectRepository,
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
		TransactionRepository:         transactionRepository,
	}
}

type CreateProjectTaskCategoryInput struct {
	OrganizationIdentity core.Identity
	ProjectIdentity      core.Identity
	Name                 string
	Color                string
}

func (i CreateProjectTaskCategoryInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewName(i.Name); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if _, err := core.NewColor(i.Color); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "color",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateProjectTaskCategoryService) Execute(input CreateProjectTaskCategoryInput) (*project.ProjectTaskCategoryDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectTaskCategoryRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if prj == nil {
		return nil, core.NewNotFoundError("project not found")
	}

	projectTaskCategory, err := project.NewProjectTaskCategory(project.NewProjectTaskCategoryInput{
		ProjectIdentity: input.ProjectIdentity,
		Name:            input.Name,
		Color:           input.Color,
	})
	if err != nil {
		return nil, err
	}

	projectTaskCategory, err = s.ProjectTaskCategoryRepository.StoreProjectTaskCategory(projectrepo.StoreProjectTaskCategoryParams{ProjectTaskCategory: projectTaskCategory})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project.ProjectTaskCategoryToDto(projectTaskCategory), nil
}
