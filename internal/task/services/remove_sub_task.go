package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type RemoveSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TransactionRepository core.TransactionRepository
}

func NewRemoveSubTaskService(
	taskRepository taskrepo.TaskRepository,
	transactionRepository core.TransactionRepository,
) *RemoveSubTaskService {
	return &RemoveSubTaskService{
		TaskRepository:        taskRepository,
		TransactionRepository: transactionRepository,
	}
}

type RemoveSubTaskInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
	SubTaskIdentity      core.Identity
}

func (i RemoveSubTaskInput) Validate() error { return nil }

func (s *RemoveSubTaskService) Execute(input RemoveSubTaskInput) error {
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

	subTask := tsk.GetSubTaskByIdentity(input.SubTaskIdentity)
	if subTask == nil {
		tx.Rollback()
		return core.NewNotFoundError("sub task not found")
	}

	err = s.TaskRepository.RemoveSubTask(taskrepo.RemoveSubTaskParams{
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
