package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type CompleteTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TaskActionRepository  taskrepo.TaskActionRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
	TransactionRepository core.TransactionRepository
}

func NewCompleteTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *CompleteTaskService {
	return &CompleteTaskService{
		TaskRepository:        taskRepository,
		TaskActionRepository:  taskActionRepository,
		ProjectUserRepository: projectUserRepository,
		TransactionRepository: transactionRepository,
	}
}

type CompleteTaskInput struct {
	OrganizationIdentity  *core.Identity
	TaskIdentity          core.Identity
	UserCompleterIdentity core.Identity
}

func (i CompleteTaskInput) Validate() error { return nil }

func (s *CompleteTaskService) Execute(input CompleteTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskRepository.SetTransaction(tx)
	s.TaskActionRepository.SetTransaction(tx)
	s.ProjectUserRepository.SetTransaction(tx)

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

	userCompleter, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.UserCompleterIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	if userCompleter == nil {
		tx.Rollback()
		return core.NewNotFoundError("project user completer not found")
	}

	var actionType task.TaskActionType
	if tsk.IsCompleted() {
		actionType = task.TaskActionTypeUncomplete
	} else {
		actionType = task.TaskActionTypeComplete
	}

	if tsk.IsCompleted() {
		tsk.Uncomplete()
	} else {
		tsk.Complete()
	}

	err = s.TaskRepository.UpdateTask(taskrepo.UpdateTaskParams{
		Task: tsk,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	taskAction := tsk.RegisterAction(actionType, &userCompleter.User)
	_, err = s.TaskActionRepository.StoreTaskAction(taskrepo.StoreTaskActionParams{
		TaskAction: &taskAction,
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
