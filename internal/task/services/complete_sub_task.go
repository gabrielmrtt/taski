package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type CompleteSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TaskActionRepository  taskrepo.TaskActionRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
	TransactionRepository core.TransactionRepository
}

func NewCompleteSubTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *CompleteSubTaskService {
	return &CompleteSubTaskService{
		TaskRepository:        taskRepository,
		TaskActionRepository:  taskActionRepository,
		ProjectUserRepository: projectUserRepository,
		TransactionRepository: transactionRepository,
	}
}

type CompleteSubTaskInput struct {
	OrganizationIdentity  *core.Identity
	TaskIdentity          core.Identity
	SubTaskIdentity       core.Identity
	UserCompleterIdentity core.Identity
}

func (i CompleteSubTaskInput) Validate() error { return nil }

func (s *CompleteSubTaskService) Execute(input CompleteSubTaskInput) error {
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

	subTask := tsk.GetSubTaskByIdentity(input.SubTaskIdentity)
	if subTask == nil {
		tx.Rollback()
		return core.NewNotFoundError("sub task not found")
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
	if subTask.IsCompleted() {
		actionType = task.TaskActionTypeSubTaskUncomplete
	} else {
		actionType = task.TaskActionTypeSubTaskComplete
	}

	if subTask.IsCompleted() {
		subTask.Uncomplete()
	} else {
		subTask.Complete()
	}

	err = s.TaskRepository.UpdateSubTask(taskrepo.UpdateSubTaskParams{
		Task:    tsk,
		SubTask: subTask,
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
