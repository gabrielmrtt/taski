package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type AddSubTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewAddSubTaskService(
	taskRepository taskrepo.TaskRepository,
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *AddSubTaskService {
	return &AddSubTaskService{
		TaskRepository:        taskRepository,
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type AddSubTaskInput struct {
	ProjectIdentity core.Identity
	TaskIdentity    core.Identity
	Name            string
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

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
