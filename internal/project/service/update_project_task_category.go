package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type UpdateProjectTaskCategoryService struct {
	ProjectRepository             projectrepo.ProjectRepository
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
	TransactionRepository         core.TransactionRepository
}

func NewUpdateProjectTaskCategoryService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
	transactionRepository core.TransactionRepository,
) *UpdateProjectTaskCategoryService {
	return &UpdateProjectTaskCategoryService{
		ProjectRepository:             projectRepository,
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
		TransactionRepository:         transactionRepository,
	}
}

type UpdateProjectTaskCategoryInput struct {
	OrganizationIdentity        core.Identity
	ProjectIdentity             core.Identity
	ProjectTaskCategoryIdentity core.Identity
	Name                        *string
	Color                       *string
}

func (i UpdateProjectTaskCategoryInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if i.Color != nil {
		_, err := core.NewColor(*i.Color)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "color",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateProjectTaskCategoryService) Execute(input UpdateProjectTaskCategoryInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectTaskCategoryRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if prj == nil {
		return core.NewNotFoundError("project not found")
	}

	projectTaskCategory, err := s.ProjectTaskCategoryRepository.GetProjectTaskCategoryByIdentity(projectrepo.GetProjectTaskCategoryByIdentityParams{
		ProjectTaskCategoryIdentity: &input.ProjectTaskCategoryIdentity,
		ProjectIdentity:             &input.ProjectIdentity,
	})
	if err != nil {
		return err
	}

	if projectTaskCategory == nil {
		return core.NewNotFoundError("project task category not found")
	}

	if input.Name != nil {
		err = projectTaskCategory.ChangeName(*input.Name)
		if err != nil {
			return err
		}
	}

	if input.Color != nil {
		err = projectTaskCategory.ChangeColor(*input.Color)
		if err != nil {
			return err
		}
	}

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
