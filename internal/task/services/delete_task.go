package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type DeleteTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TaskActionRepository  taskrepo.TaskActionRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTaskService {
	return &DeleteTaskService{
		TaskRepository:        taskRepository,
		TaskActionRepository:  taskActionRepository,
		ProjectUserRepository: projectUserRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteTaskInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
	UserDeleterIdentity  core.Identity
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

	userDeleter, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.UserDeleterIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	if userDeleter == nil {
		tx.Rollback()
		return core.NewNotFoundError("project user deleter not found")
	}

	tsk.Delete()

	err = s.TaskRepository.UpdateTask(taskrepo.UpdateTaskParams{
		Task: tsk,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	taskAction := tsk.RegisterAction(task.TaskActionTypeDelete, &userDeleter.User)
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
