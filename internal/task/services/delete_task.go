package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type DeleteTaskService struct {
	TaskRepository        taskrepo.TaskRepository
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteTaskService(
	taskRepository taskrepo.TaskRepository,
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTaskService {
	return &DeleteTaskService{
		TaskRepository:        taskRepository,
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteTaskInput struct {
	ProjectIdentity core.Identity
	TaskIdentity    core.Identity
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
		ProjectIdentity: &input.ProjectIdentity,
		TaskIdentity:    input.TaskIdentity,
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
