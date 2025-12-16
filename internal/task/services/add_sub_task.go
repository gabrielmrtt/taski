package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type AddSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TransactionRepository core.TransactionRepository
	TaskActionRepository  taskrepo.TaskActionRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
}

func NewAddSubTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *AddSubTaskService {
	return &AddSubTaskService{
		TaskRepository:        taskRepository,
		TransactionRepository: transactionRepository,
		TaskActionRepository:  taskActionRepository,
		ProjectUserRepository: projectUserRepository,
	}
}

type AddSubTaskInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
	Name                 string
	UserCreatorIdentity  core.Identity
}

func (i AddSubTaskInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewName(i.Name); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *AddSubTaskService) Execute(input AddSubTaskInput) error {
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

	userCreator, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if userCreator == nil {
		tx.Rollback()
		return core.NewNotFoundError("project user creator not found")
	}

	subTask, err := task.NewSubTask(task.NewSubTaskInput{
		Name: input.Name,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.TaskRepository.AddSubTask(taskrepo.AddSubTaskParams{
		Task:    tsk,
		SubTask: subTask,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	taskAction := tsk.RegisterAction(task.TaskActionTypeAddSubTask, &userCreator.User)
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
