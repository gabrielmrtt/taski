package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type CreateProjectTaskStatusService struct {
	ProjectRepository           projectrepo.ProjectRepository
	ProjectTaskStatusRepository projectrepo.ProjectTaskStatusRepository
	TransactionRepository       core.TransactionRepository
}

func NewCreateProjectTaskStatusService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectTaskStatusService {
	return &CreateProjectTaskStatusService{
		ProjectRepository:           projectRepository,
		ProjectTaskStatusRepository: projectTaskStatusRepository,
		TransactionRepository:       transactionRepository,
	}
}

type CreateProjectTaskStatusInput struct {
	OrganizationIdentity     core.Identity
	ProjectIdentity          core.Identity
	Name                     string
	Color                    string
	ShouldSetTaskToCompleted bool
	IsDefault                bool
	ShouldUseOrder           bool
}

func (i CreateProjectTaskStatusInput) Validate() error {
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

	if i.IsDefault && i.ShouldSetTaskToCompleted {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "is_default",
			Error: "is default and should set task to completed cannot be true at the same time",
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateProjectTaskStatusService) Execute(input CreateProjectTaskStatusInput) (*project.ProjectTaskStatusDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.ProjectTaskStatusRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)

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

	var order *int8 = nil

	if input.ShouldUseOrder {
		lastOrder, err := s.ProjectTaskStatusRepository.GetLastTaskStatusOrder(projectrepo.GetLastTaskStatusOrderParams{
			ProjectIdentity: &prj.Identity,
		})
		if err != nil {
			return nil, err
		}

		lastOrder++

		order = &lastOrder
	}

	projectTaskStatus, err := project.NewProjectTaskStatus(project.NewProjectTaskStatusInput{
		ProjectIdentity:          prj.Identity,
		Name:                     input.Name,
		Color:                    input.Color,
		Order:                    order,
		ShouldSetTaskToCompleted: input.ShouldSetTaskToCompleted,
		IsDefault:                input.IsDefault,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	projectTaskStatus, err = s.ProjectTaskStatusRepository.StoreProjectTaskStatus(projectrepo.StoreProjectTaskStatusParams{ProjectTaskStatus: projectTaskStatus})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return project.ProjectTaskStatusToDto(projectTaskStatus), nil
}
