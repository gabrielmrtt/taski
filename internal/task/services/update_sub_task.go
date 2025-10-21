package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type UpdateSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateSubTaskService(
	taskRepository taskrepo.TaskRepository,
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *UpdateSubTaskService {
	return &UpdateSubTaskService{
		TaskRepository:        taskRepository,
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateSubTaskInput struct {
	ProjectIdentity core.Identity
	TaskIdentity    core.Identity
	SubTaskIdentity core.Identity
	Name            *string
	Completed       *bool
}

func (i UpdateSubTaskInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		if _, err := core.NewName(*i.Name); err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateSubTaskService) Execute(input UpdateSubTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if prj == nil {
		tx.Rollback()
		return core.NewNotFoundError("project not found")
	}

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity:    input.TaskIdentity,
		ProjectIdentity: &input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if tsk == nil {
		tx.Rollback()
		return core.NewNotFoundError("task not found")
	}

	subTask := tsk.GetSubTaskByIdentity(input.SubTaskIdentity)
	if subTask == nil {
		tx.Rollback()
		return core.NewNotFoundError("sub task not found")
	}

	if input.Name != nil {
		err = subTask.ChangeName(*input.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Completed != nil {
		if *input.Completed {
			subTask.Complete()
		} else {
			subTask.Uncomplete()
		}
	}

	err = s.TaskRepository.UpdateSubTask(taskrepo.UpdateSubTaskParams{
		Task:    tsk,
		SubTask: subTask,
	})
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
