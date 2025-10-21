package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type UpdateTaskService struct {
	TaskRepository                taskrepo.TaskRepository
	ProjectRepository             projectrepo.ProjectRepository
	ProjectTaskStatusRepository   projectrepo.ProjectTaskStatusRepository
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
	ProjectUserRepository         projectrepo.ProjectUserRepository
	TransactionRepository         core.TransactionRepository
}

func NewUpdateTaskService(
	taskRepository taskrepo.TaskRepository,
	projectRepository projectrepo.ProjectRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateTaskService {
	return &UpdateTaskService{
		TaskRepository:                taskRepository,
		ProjectRepository:             projectRepository,
		ProjectTaskStatusRepository:   projectTaskStatusRepository,
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
		ProjectUserRepository:         projectUserRepository,
		TransactionRepository:         transactionRepository,
	}
}

type UpdateTaskInput struct {
	ProjectIdentity    core.Identity
	TaskIdentity       core.Identity
	StatusIdentity     *core.Identity
	CategoryIdentity   *core.Identity
	ParentTaskIdentity *core.Identity
	Name               *string
	Description        *string
	EstimatedMinutes   *int16
	PriorityLevel      *task.TaskPriorityLevels
	DueDate            *int64
	Users              []*core.Identity
	ChildrenTasks      []*core.Identity
	UserEditorIdentity core.Identity
}

func (i UpdateTaskInput) Validate() error {
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

	if i.Description != nil {
		_, err := core.NewDescription(*i.Description)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "description",
				Error: err.Error(),
			})
		}
	}

	if i.EstimatedMinutes != nil {
		if *i.EstimatedMinutes < 0 {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "estimated minutes",
				Error: "estimated minutes cannot be negative",
			})
		}
	}

	if i.DueDate != nil {
		if *i.DueDate < datetimeutils.EpochNow() {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "due date",
				Error: "due date cannot be in the past",
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateTaskService) Execute(input UpdateTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)
	s.ProjectTaskStatusRepository.SetTransaction(tx)
	s.ProjectTaskCategoryRepository.SetTransaction(tx)

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

	if tsk.ProjectIdentity != input.ProjectIdentity {
		tx.Rollback()
		return core.NewNotFoundError("task not found in project")
	}

	if input.StatusIdentity != nil {
		status, err := s.ProjectTaskStatusRepository.GetProjectTaskStatusByIdentity(projectrepo.GetProjectTaskStatusByIdentityParams{
			ProjectTaskStatusIdentity: input.StatusIdentity,
			ProjectIdentity:           &input.ProjectIdentity,
		})
		if err != nil {
			tx.Rollback()
			return err
		}

		if status == nil {
			tx.Rollback()
			return core.NewNotFoundError("project task status not found")
		}

		err = tsk.ChangeStatus(status, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.CategoryIdentity != nil {
		category, err := s.ProjectTaskCategoryRepository.GetProjectTaskCategoryByIdentity(projectrepo.GetProjectTaskCategoryByIdentityParams{
			ProjectTaskCategoryIdentity: input.CategoryIdentity,
			ProjectIdentity:             &input.ProjectIdentity,
		})
		if err != nil {
			tx.Rollback()
			return err
		}

		if category == nil {
			tx.Rollback()
			return core.NewNotFoundError("project task category not found")
		}

		err = tsk.ChangeCategory(category, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Name != nil {
		err = tsk.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Description != nil {
		err = tsk.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.EstimatedMinutes != nil {
		err = tsk.ChangeEstimatedMinutes(*input.EstimatedMinutes, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.PriorityLevel != nil {
		err = tsk.ChangePriorityLevel(*input.PriorityLevel, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.DueDate != nil {
		err = tsk.ChangeDueDate(*input.DueDate, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Users != nil {
		tsk.ClearUsers()
		for _, userIdentity := range input.Users {
			user, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
				ProjectIdentity: input.ProjectIdentity,
				UserIdentity:    *userIdentity,
			})
			if err != nil {
				tx.Rollback()
				return err
			}

			if user == nil {
				tx.Rollback()
				return core.NewNotFoundError("project user not found")
			}

			tsk.AddUser(&task.TaskUser{
				User: &user.User,
			})
		}
	}

	if input.ChildrenTasks != nil {
		for _, childTask := range tsk.ChildrenTasks {
			tsk.RemoveChildTask(childTask)

			err := s.TaskRepository.UpdateTask(taskrepo.UpdateTaskParams{
				Task: childTask,
			})
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, childTaskIdentity := range input.ChildrenTasks {
			childTask, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
				TaskIdentity:    *childTaskIdentity,
				ProjectIdentity: &input.ProjectIdentity,
			})
			if err != nil {
				tx.Rollback()
				return err
			}

			if childTask == nil {
				tx.Rollback()
				return core.NewNotFoundError("child task not found")
			}

			tsk.AddChildTask(childTask)
		}
	}

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
