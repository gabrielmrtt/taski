package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type UpdateProjectTaskStatusService struct {
	ProjectRepository           projectrepo.ProjectRepository
	ProjectTaskStatusRepository projectrepo.ProjectTaskStatusRepository
	TransactionRepository       core.TransactionRepository
}

func NewUpdateProjectTaskStatusService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	transactionRepository core.TransactionRepository,
) *UpdateProjectTaskStatusService {
	return &UpdateProjectTaskStatusService{
		ProjectRepository:           projectRepository,
		ProjectTaskStatusRepository: projectTaskStatusRepository,
		TransactionRepository:       transactionRepository,
	}
}

type UpdateProjectTaskStatusInput struct {
	OrganizationIdentity      core.Identity
	ProjectIdentity           core.Identity
	ProjectTaskStatusIdentity core.Identity
	Name                      *string
	Color                     *string
	Order                     *int8
	ShouldSetTaskToCompleted  *bool
	IsDefault                 *bool
}

func (i UpdateProjectTaskStatusInput) Validate() error {
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

	if i.IsDefault != nil {
		if i.ShouldSetTaskToCompleted != nil {
			if *i.IsDefault && *i.ShouldSetTaskToCompleted {
				fields = append(fields, core.InvalidInputErrorField{
					Field: "is_default",
					Error: "is default and should set task to completed cannot be true at the same time",
				})
			}
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateProjectTaskStatusService) Execute(input UpdateProjectTaskStatusInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectTaskStatusRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if prj == nil {
		tx.Rollback()
		return core.NewNotFoundError("project not found")
	}

	projectTaskStatus, err := s.ProjectTaskStatusRepository.GetProjectTaskStatusByIdentity(projectrepo.GetProjectTaskStatusByIdentityParams{
		ProjectTaskStatusIdentity: &input.ProjectTaskStatusIdentity,
		ProjectIdentity:           &input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if projectTaskStatus == nil {
		tx.Rollback()
		return core.NewNotFoundError("project task status not found")
	}

	if input.Name != nil {
		err = projectTaskStatus.ChangeName(*input.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Color != nil {
		err = projectTaskStatus.ChangeColor(*input.Color)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Order != nil {
		err = projectTaskStatus.ChangeOrder(*input.Order)
		if err != nil {
			tx.Rollback()
			return err
		}

		statuses, err := s.ProjectTaskStatusRepository.ListProjectTaskStatusesBy(projectrepo.ListProjectTaskStatusesByParams{
			Filters: projectrepo.ProjectTaskStatusFilters{
				ProjectIdentity: &prj.Identity,
				Order: &core.ComparableFilter[int8]{
					GreaterThanOrEqual: input.Order,
				},
			},
		})
		if err != nil {
			tx.Rollback()
			return err
		}

		isNext := func(current int8, other int8) bool {
			next := current + 1

			return next == other
		}

		current := *input.Order
		for idx, status := range statuses {
			if isNext(current, *status.Order) {
				current = *status.Order
				status.ChangeOrder(current + int8(idx+1))

				err = s.ProjectTaskStatusRepository.UpdateProjectTaskStatus(projectrepo.UpdateProjectTaskStatusParams{ProjectTaskStatus: &status})
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	err = s.ProjectTaskStatusRepository.UpdateProjectTaskStatus(projectrepo.UpdateProjectTaskStatusParams{ProjectTaskStatus: projectTaskStatus})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
