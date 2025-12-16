package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type UpdateSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	TaskActionRepository  taskrepo.TaskActionRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateSubTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateSubTaskService {
	return &UpdateSubTaskService{
		TaskRepository:        taskRepository,
		TaskActionRepository:  taskActionRepository,
		ProjectUserRepository: projectUserRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateSubTaskInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
	SubTaskIdentity      core.Identity
	Name                 *string
	UserEditorIdentity   core.Identity
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

	userEditor, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.UserEditorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	if userEditor == nil {
		tx.Rollback()
		return core.NewNotFoundError("project user editor not found")
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

	err = s.TaskRepository.UpdateSubTask(taskrepo.UpdateSubTaskParams{
		Task:    tsk,
		SubTask: subTask,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	taskAction := tsk.RegisterAction(task.TaskActionTypeUpdateSubTask, &userEditor.User)
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
