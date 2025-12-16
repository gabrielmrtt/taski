package taskservice

import (
	"strconv"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type CreateTaskService struct {
	TaskRepository                taskrepo.TaskRepository
	TaskActionRepository          taskrepo.TaskActionRepository
	ProjectRepository             projectrepo.ProjectRepository
	ProjectUserRepository         projectrepo.ProjectUserRepository
	ProjectTaskStatusRepository   projectrepo.ProjectTaskStatusRepository
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
	TransactionRepository         core.TransactionRepository
}

func NewCreateTaskService(
	taskRepository taskrepo.TaskRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	projectRepository projectrepo.ProjectRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
	transactionRepository core.TransactionRepository,
) *CreateTaskService {
	return &CreateTaskService{
		TaskRepository:                taskRepository,
		TaskActionRepository:          taskActionRepository,
		ProjectRepository:             projectRepository,
		ProjectUserRepository:         projectUserRepository,
		ProjectTaskStatusRepository:   projectTaskStatusRepository,
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
		TransactionRepository:         transactionRepository,
	}
}

type CreateSubTaskInput struct {
	Name string
}

func (i CreateSubTaskInput) Validate() error {
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

type CreateTaskInput struct {
	OrganizationIdentity *core.Identity
	ProjectIdentity      core.Identity
	StatusIdentity       core.Identity
	CategoryIdentity     *core.Identity
	ParentTaskIdentity   *core.Identity
	Name                 string
	Description          string
	EstimatedMinutes     *int16
	PriorityLevel        task.TaskPriorityLevels
	DueDate              *core.DateTime
	SubTasks             []*CreateSubTaskInput
	Users                []*core.Identity
	ChildrenTasks        []*core.Identity
	UserCreatorIdentity  core.Identity
}

func (i CreateTaskInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewName(i.Name); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if _, err := core.NewDescription(i.Description); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "description",
			Error: err.Error(),
		})
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
		if i.DueDate.IsBefore(core.NewDateTime()) {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "due date",
				Error: "due date cannot be in the past",
			})
		}
	}

	if i.SubTasks != nil {
		for index, subTaskInput := range i.SubTasks {
			if err := subTaskInput.Validate(); err != nil {
				fieldName := "subtasks[" + strconv.Itoa(index) + "].name"
				invalidInputError, ok := err.(*core.InvalidInputError)
				if ok {
					for _, field := range (*invalidInputError).Fields {
						fields = append(fields, core.InvalidInputErrorField{
							Field: fieldName,
							Error: field.Error,
						})
					}
				} else {
					fields = append(fields, core.InvalidInputErrorField{
						Field: fieldName,
						Error: err.Error(),
					})
				}
			}
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateTaskService) Execute(input CreateTaskInput) (*task.TaskDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TaskRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)
	s.ProjectTaskStatusRepository.SetTransaction(tx)
	s.ProjectTaskCategoryRepository.SetTransaction(tx)
	s.TaskActionRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if prj == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project not found")
	}

	userCreator, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
		UserIdentity:    input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if userCreator == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project user creator not found")
	}

	status, err := s.ProjectTaskStatusRepository.GetProjectTaskStatusByIdentity(projectrepo.GetProjectTaskStatusByIdentityParams{
		ProjectIdentity:           &input.ProjectIdentity,
		ProjectTaskStatusIdentity: &input.StatusIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if status == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project task status not found")
	}

	var category *project.ProjectTaskCategory = nil
	if input.CategoryIdentity != nil {
		category, err = s.ProjectTaskCategoryRepository.GetProjectTaskCategoryByIdentity(projectrepo.GetProjectTaskCategoryByIdentityParams{
			ProjectIdentity:             &input.ProjectIdentity,
			ProjectTaskCategoryIdentity: input.CategoryIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if category == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("project task category not found")
		}
	}

	var parentTask *task.Task = nil
	if input.ParentTaskIdentity != nil {
		parentTask, err = s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
			TaskIdentity:    *input.ParentTaskIdentity,
			ProjectIdentity: &input.ProjectIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if parentTask == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("parent task not found")
		}
	}

	subTasks := make([]*task.SubTask, 0)
	for _, subTaskInput := range input.SubTasks {
		subTask, err := task.NewSubTask(task.NewSubTaskInput{
			Name: subTaskInput.Name,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		subTasks = append(subTasks, subTask)
	}

	users := make([]*task.TaskUser, 0)
	for _, userIdentity := range input.Users {
		user, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
			ProjectIdentity: input.ProjectIdentity,
			UserIdentity:    *userIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if user == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("project user not found")
		}

		users = append(users, &task.TaskUser{
			User: &user.User,
		})
	}

	childrenTasks := make([]*task.Task, 0)
	for _, childTaskIdentity := range input.ChildrenTasks {
		childTask, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
			TaskIdentity:    *childTaskIdentity,
			ProjectIdentity: &input.ProjectIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if childTask == nil {
			tx.Rollback()
			return nil, core.NewNotFoundError("child task not found")
		}

		childrenTasks = append(childrenTasks, childTask)
	}

	tsk, err := task.NewTask(task.NewTaskInput{
		ProjectIdentity:     input.ProjectIdentity,
		Status:              status,
		Category:            category,
		ParentTaskIdentity:  &parentTask.Identity,
		Name:                input.Name,
		Description:         input.Description,
		EstimatedMinutes:    input.EstimatedMinutes,
		PriorityLevel:       input.PriorityLevel,
		DueDate:             input.DueDate,
		SubTasks:            subTasks,
		ChildrenTasks:       childrenTasks,
		Users:               users,
		UserCreatorIdentity: &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = s.TaskRepository.StoreTask(taskrepo.StoreTaskParams{
		Task: tsk,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	taskAction := tsk.RegisterAction(task.TaskActionTypeCreate, &userCreator.User)

	_, err = s.TaskActionRepository.StoreTaskAction(taskrepo.StoreTaskActionParams{
		TaskAction: &taskAction,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return task.TaskToDto(tsk), nil
}
