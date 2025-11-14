package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type ChangeTaskStatusService struct {
	TaskRepository              taskrepo.TaskRepository
	ProjectTaskStatusRepository projectrepo.ProjectTaskStatusRepository
	TransactionRepository       core.TransactionRepository
}

func NewChangeTaskStatusService(
	taskRepository taskrepo.TaskRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	transactionRepository core.TransactionRepository,
) *ChangeTaskStatusService {
	return &ChangeTaskStatusService{
		TaskRepository:              taskRepository,
		ProjectTaskStatusRepository: projectTaskStatusRepository,
		TransactionRepository:       transactionRepository,
	}
}

type ChangeTaskStatusInput struct {
	OrganizationIdentity      *core.Identity
	TaskIdentity              core.Identity
	ProjectTaskStatusIdentity *core.Identity
	AdvanceOrder              bool
	ChangedByUserIdentity     core.Identity
}

func (i ChangeTaskStatusInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.AdvanceOrder {
		if i.ProjectTaskStatusIdentity != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "project_task_status_identity",
				Error: "project task status identity cannot be provided when advancing order is enabled",
			})
		}
	} else {
		if i.ProjectTaskStatusIdentity == nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "project_task_status_identity",
				Error: "project task status identity is required when advancing order is disabled",
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *ChangeTaskStatusService) Execute(input ChangeTaskStatusInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskRepository.SetTransaction(tx)
	s.ProjectTaskStatusRepository.SetTransaction(tx)

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity:         input.TaskIdentity,
		OrganizationIdentity: input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if tsk == nil {
		tx.Rollback()
		return core.NewNotFoundError("task not found")
	}

	projectStatuses, err := s.ProjectTaskStatusRepository.ListProjectTaskStatusesBy(projectrepo.ListProjectTaskStatusesByParams{
		Filters: projectrepo.ProjectTaskStatusFilters{
			ProjectIdentity: &tsk.ProjectIdentity,
		},
		SortInput: &core.SortInput{
			By:        &[]string{"order"}[0],
			Direction: &[]core.SortDirection{"asc"}[0],
		},
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	var nextStatus *project.ProjectTaskStatus = nil

	if input.AdvanceOrder {
		for _, status := range projectStatuses {
			if status.Order != nil {
				if *status.Order == (*tsk.Status.Order + 1) {
					nextStatus = &status
					break
				}
			}
		}

		if nextStatus == nil {
			nextStatus = &projectStatuses[0]
		}
	} else {
		if input.ProjectTaskStatusIdentity == nil {
			tx.Rollback()
			return core.NewConflictError("project task status identity is required")
		}

		projectTaskStatus, err := s.ProjectTaskStatusRepository.GetProjectTaskStatusByIdentity(projectrepo.GetProjectTaskStatusByIdentityParams{
			ProjectTaskStatusIdentity: input.ProjectTaskStatusIdentity,
			ProjectIdentity:           &tsk.ProjectIdentity,
		})
		if err != nil {
			tx.Rollback()
			return err
		}

		if projectTaskStatus == nil {
			tx.Rollback()
			return core.NewNotFoundError("project task status not found")
		}

		nextStatus = projectTaskStatus
	}

	if nextStatus != nil {
		err = tsk.ChangeStatus(nextStatus, &input.ChangedByUserIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = s.TaskRepository.UpdateTask(taskrepo.UpdateTaskParams{Task: tsk})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
