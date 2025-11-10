package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type DeleteTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteTaskService(
	taskRepository taskrepo.TaskRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTaskService {
	return &DeleteTaskService{
		TaskRepository:        taskRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteTaskInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
}

func (i DeleteTaskInput) Validate() error { return nil }

func (s *DeleteTaskService) Execute(input DeleteTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskRepository.SetTransaction(tx)

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

	tsk.Delete()

	err = s.TaskRepository.UpdateTask(taskrepo.UpdateTaskParams{
		Task: tsk,
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
